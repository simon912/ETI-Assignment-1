// trip.go - the backend for trip management
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Trips struct {
	TripID              int            `json:"Trip ID"`
	PickUpLocation      string         `json:"Pick-Up Location"`
	AltPickUpLocation   sql.NullString `json:"Alternate Pick-Up Location"`
	StartTravelTime     string         `json:"Start Traveling Time"`
	DestinationLocation string         `json:"Destination Location"`
	PassengerNoLeft     int            `json:"Number of Passengers Left"`
	MaxPassengerNo      int            `json:"Maximum Number of Passengers"`
	Status              string         `json:"Status"` // pending or active
	Publisher           string         `json:"Publisher"`
}

type Enrollment struct {
	EnrollmentID int    `json:"Enrollment ID"`
	Username     string `json:"Username"`
	TripID       int    `json:"Trip ID"`
}

// REST endpoint for Trips
func main() {
	router := mux.NewRouter()
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})
	handler := corsOptions.Handler(router)
	// To retrieve all trips for Passenger to View
	router.HandleFunc("/api/v1/trips", GetAllTrip).Methods("GET")
	// For Car Owner to publish trip
	router.HandleFunc("/api/v1/publishtrip/{username}", PublishTrip).Methods("POST")
	// For Car Owner to start trip
	router.HandleFunc("/api/v1/starttrip/{tripid}", StartTrip).Methods("PUT")
	// For Passenger to enroll themselves into any trip
	router.HandleFunc("/api/v1/enroll/{tripid}/{username}", EnrollUser).Methods("PUT")
	// For Car Owner to cancel trip, works if they cancel it 30 mins before start travel time
	router.HandleFunc("/api/v1/canceltrip/{tripid}", CancelTrip).Methods("DELETE")
	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", handler))
}

// Helper function to handle NULL values in SQL parameters
func nullOrValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	return value
}

// ----------------------------- Endpoint from User ----------------------------------------
// Helper function to check if user exists in the table
func userExists(username string) bool {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query := "SELECT COUNT(*) FROM Users WHERE Username = ?"
	var count int
	err = db.QueryRow(query, username).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	return count > 0
}

// ----------------------------- Endpoint for Car Pooling Trips ----------------------------------------
// Function to retrieve all trip with the detail
func GetAllTrip(w http.ResponseWriter, r *http.Request) {
	//test case for retrieve all: curl -X GET http://localhost:5000/api/v1/trips
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	// handle error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// database operation
	defer db.Close()
	results, err := db.Query("SELECT * FROM Trips")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var trips []Trips
	for results.Next() {
		var t Trips
		err = results.Scan(&t.TripID, &t.PickUpLocation, &t.AltPickUpLocation, &t.StartTravelTime, &t.DestinationLocation, &t.PassengerNoLeft, &t.MaxPassengerNo, &t.Status, &t.Publisher)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Parse the time string to time.Time
		parsedTime, err := time.Parse("15:04:05", t.StartTravelTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert the time to UTC
		parsedTime = parsedTime.UTC()

		// Format the time as 12-hour format with AM/PM
		formattedTime := parsedTime.Format("03:04 PM")
		t.StartTravelTime = formattedTime
		trips = append(trips, t)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trips)
}

// Function for Passenger to enroll into a trip of their coice
func EnrollUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	// handle error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	tripID := params["tripid"]
	username := params["username"]

	// Convert tripID to int
	tripIDInt, err := strconv.Atoi(tripID)
	if err != nil {
		http.Error(w, "Invalid trip ID", http.StatusBadRequest)
		return
	}

	// Check if the trip exists
	trip, found := getTripByID(tripIDInt)
	if !found {
		http.Error(w, "Trip not found", http.StatusNotFound)
		return
	}

	// Check if the trip is full
	if trip.PassengerNoLeft <= 0 {
		http.Error(w, `{"error": "Trip is already full"}`, http.StatusBadRequest)
		return
	}
	// Check if user is already enrolled
	if isUserEnrolled(tripIDInt, username) {
		http.Error(w, `{"error": "User is already enrolled in this trip"}`, http.StatusBadRequest)
		return
	}
	if isTripActive(tripIDInt) {
		http.Error(w, `{"error": "Trip is Active"}`, http.StatusBadRequest)
		return
	}
	// Insert enrollment record
	enrollmentQuery := "INSERT INTO Enrollment (Username, TripID) VALUES (?, ?)"
	_, err = db.Exec(enrollmentQuery, username, tripIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err.Error())
	}

	// Update the trip and decrement the available passenger slots
	updatePassengerNoLeft(tripIDInt, username, trip)

	// Return success response
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "User %s has been enrolled in trip ID %s!\n", username, tripID)
}

// Helper function to check if Trip is active
func isTripActive(tripID int) bool {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return false
	}
	defer db.Close()

	var status string
	query := "SELECT Status FROM Trips WHERE TripID = ?"
	err = db.QueryRow(query, tripID).Scan(&status)
	if err != nil {
		fmt.Printf("Error retrieving trip status: %v\n", err)
		return false
	}

	// Check if the trip is marked as "Active"
	return status == "Active"
}

// Helper function to get a trip by ID
func getTripByID(tripID int) (Trips, bool) {
	// Database query to get the trip by ID
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		// Handle error
		fmt.Printf("Error opening database: %v\n", err)
		return Trips{}, false
	}
	defer db.Close()

	query := "SELECT * FROM Trips WHERE TripID = ?"
	var t Trips
	var startTravelTimeStr string
	err = db.QueryRow(query, tripID).Scan(
		&t.TripID, &t.PickUpLocation, &t.AltPickUpLocation, &startTravelTimeStr, &t.DestinationLocation,
		&t.PassengerNoLeft, &t.MaxPassengerNo, &t.Status, &t.Publisher,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Trip with ID %d not found\n", tripID)
			return Trips{}, false
		}
		// Handle other errors
		fmt.Printf("Error querying database: %v\n", err)
		return Trips{}, false
	}

	return t, true
}

// Helper function to update the trip and decrement the available passenger slots
func updatePassengerNoLeft(tripID int, username string, trip Trips) {
	// Update the trip with the new passenger
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		// Handle error
		return
	}
	defer db.Close()
	// Update the trip with the new passenger
	query := "UPDATE Trips SET PassengerNoLeft = ? WHERE TripID=?"
	_, err = db.Exec(query, trip.PassengerNoLeft-1, tripID)
	if err != nil {
		// Handle error
		fmt.Printf("Error updating trip: %v\n", err)
		return
	}
}

// Custom unmarshal function for LicenseNo
func (t *Trips) UnmarshalJSON(data []byte) error {
	type Alias Trips
	aux := &struct {
		AltPickUpLocation string `json:"Alternate Pick-Up Location"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.AltPickUpLocation = sql.NullString{String: aux.AltPickUpLocation, Valid: aux.AltPickUpLocation != ""}
	return nil
}

// Function for Car Owner to publish trip
func PublishTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	if !userExists(username) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Decode the JSON request into a new trip
	var newTrip Trips
	err := json.NewDecoder(r.Body).Decode(&newTrip)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		panic(err.Error())
	}
	parsedTime, err := time.Parse("15:04:05", newTrip.StartTravelTime)
	if err != nil {
		http.Error(w, "Invalid time format", http.StatusBadRequest)
		panic(err.Error())
	}
	newTrip.StartTravelTime = parsedTime.Format("15:04:05")
	newTrip.Publisher = username

	// Insert the new trip into the Trips table
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Insert the new trip and auto-increment the ID
	query := "INSERT INTO Trips (PickUpLocation, AltPickUpLocation, StartTravelTime, DestinationLocation, PassengerNoLeft, MaxPassengerNo, Status, Publisher) VALUES (?, ?, STR_TO_DATE(?, '%H:%i:%s'), ?, ?, ?, ?, ?)"
	result, err := db.Exec(query, newTrip.PickUpLocation, nullOrValue(newTrip.AltPickUpLocation), newTrip.StartTravelTime, newTrip.DestinationLocation, newTrip.PassengerNoLeft, newTrip.MaxPassengerNo, "Pending", newTrip.Publisher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err.Error())
	}

	// Get the auto-incremented ID
	newID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err.Error())
	}
	// Update the response with the new ID
	newTrip.TripID = int(newID)

	// Return success response with the new trip
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTrip)
}

// Helper function to check if user enrolls in that specific trip
func isUserEnrolled(tripID int, username string) bool {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		// Handle error
		fmt.Printf("Error opening database: %v\n", err)
		return false
	}
	defer db.Close()

	// Check if the user is enrolled in the specified trip
	query := "SELECT COUNT(*) FROM Enrollment WHERE TripID = ? AND Username = ?"
	var count int
	err = db.QueryRow(query, tripID, username).Scan(&count)
	if err != nil {
		// Handle error
		fmt.Printf("Error querying database: %v\n", err)
		return false
	}

	return count > 0
}

// Function to start a trip by changing the status from Pending to Active
func StartTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tripID := params["tripid"]

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var startTrip Trips
	query := "SELECT MaxPassengerNo, PassengerNoLeft FROM Trips WHERE TripID = ?"
	err = db.QueryRow(query, tripID).Scan(&startTrip.MaxPassengerNo, &startTrip.PassengerNoLeft)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if PassengerNoLeft is equal to MaxPassengerNo
	if startTrip.MaxPassengerNo == startTrip.PassengerNoLeft {
		http.Error(w, "Trip cannot start if no one is enrolled into your trip", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	query = "UPDATE Trips SET Status=? WHERE TripID = ?"
	_, err = db.Exec(query, "Active", tripID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Trip %s Status has been changed to Active\n", tripID)
}

// Function to cancel trip by deleting them from the Trip table only if it is within the 30 minutes cancellation window
func CancelTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tripID := params["tripid"]
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Convert tripID to int
	tripIDInt, err := strconv.Atoi(tripID)
	if err != nil {
		http.Error(w, "Invalid trip ID", http.StatusBadRequest)
		return
	}

	// Check if the trip exists
	_, found := getTripByID(tripIDInt)
	if !found {
		http.Error(w, "Trip not found", http.StatusNotFound)
		return
	}

	var dbStartTime string
	err = db.QueryRow("SELECT StartTravelTime FROM Trips WHERE TripID = ?", tripID).Scan(&dbStartTime)
	if err != nil {
		fmt.Println("Error retrieving StartTravelTime from the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse the StartTravelTime from the database
	startTime, err := time.Parse("15:04:05", dbStartTime)
	if err != nil {
		fmt.Println("Error parsing start time from the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get the current time
	currentTimeStr := time.Now().Format("15:04:05")
	currentTime, err := time.Parse("15:04:05", currentTimeStr)
	if err != nil {
		fmt.Println("Error parsing current time string:", err)
		// Handle the error as needed
		return
	}
	// Calculate the time difference
	timeDifference := currentTime.Sub(startTime)
	
	// Check if the trip can be canceled (more than or equal to 30 minutes before the start time)
	if timeDifference.Minutes() >= -30 && timeDifference.Minutes() <= 30 {
		http.Error(w, `{"error": "Cancellation window is within 30 minutes"}`, http.StatusBadRequest)
		return
	}
	// Delete the trip from the Trips table and other related enrollment from Enrollment table
	deleteEnrollmentsForTrip(tripIDInt)
	deleteTrip(w, tripIDInt)
	
	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trip ID %s has been canceled successfully\n", tripID)
}

// Helper function to delete a trip from the Trips table
func deleteTrip(w http.ResponseWriter, tripID int) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Delete the trip
	query := "DELETE FROM Trips WHERE TripID = ?"
	result, err := db.Exec(query, tripID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the trip was found and deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Trip not found", http.StatusNotFound)
		return
	}
}

// Helper function to delete enrollments for a trip
func deleteEnrollmentsForTrip(tripID int) error {
    db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
    if err != nil {
        return err
    }
    defer db.Close()

    // Delete enrollments for the specified trip
    query := "DELETE FROM Enrollment WHERE TripID = ?"
    _, err = db.Exec(query, tripID)
    return err
}
