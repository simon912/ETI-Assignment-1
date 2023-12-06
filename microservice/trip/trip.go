// main.go - the backend
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
	ID                  int            `json:"ID"`
	PickUpLocation      string         `json:"Pick-Up Location"`
	AltPickUpLocation   sql.NullString `json:"Alternate Pick-Up Location"`
	StartTravelTime     string         `json:"Start Traveling Time"`
	DestinationLocation string         `json:"Destination Location"`
	PassengerNoLeft     int            `json:"Number of Passengers Left"`
	MaxPassengerNo      int            `json:"Maximum Number of Passengers"`
	Passenger1          sql.NullString `json:"Passenger 1"`
	Passenger2          sql.NullString `json:"Passenger 2"`
	Passenger3          sql.NullString `json:"Passenger 3"`
	Status              string         `json:"Status"` // pending or active
	Publisher           string         `json:"Publisher"`
}

// Register REST endpoint
func main() {
	router := mux.NewRouter()
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})
	handler := corsOptions.Handler(router)
	// Endpoint for Car-Pooling Trips
	// To retrieve all trips for Passenger to View
	router.HandleFunc("/api/v1/trips", GetAllTrip).Methods("GET")
	// For Car Owner to publish trip
	router.HandleFunc("/api/v1/publishtrip/{username}", PublishTrip).Methods("POST")
	// For Car Owner to start trip
	// curl http://localhost:5000/api/v1/starttrip/3 -X PUT -d "{\"Number of Passengers Left\": 1, \"Maximum Number of Passengers\": 2}"
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
		err = results.Scan(&t.ID, &t.PickUpLocation, &t.AltPickUpLocation, &t.StartTravelTime, &t.DestinationLocation, &t.PassengerNoLeft, &t.MaxPassengerNo, &t.Passenger1, &t.Passenger2, &t.Passenger3, &t.Status, &t.Publisher)
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

		// Convert the time to UTC if necessary (as needed in your application)
		parsedTime = parsedTime.UTC()

		// Format the time as 12-hour format with AM/PM
		formattedTime := parsedTime.Format("03:04 PM")

		// Assign the formatted time to StartTravelTime in the Trips struct
		t.StartTravelTime = formattedTime

		trips = append(trips, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trips)
}

func EnrollUser(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Trip is already full", http.StatusBadRequest)
		return
	}

	// Check if the user is already enrolled in the trip
	if isUserEnrolled(username, trip) {
		http.Error(w, "User is already enrolled in this trip", http.StatusBadRequest)
		return
	}

	// Update the trip and decrement the available passenger slots
	updateTrip(tripIDInt, username, trip)

	// Return success response
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "User %s has been enrolled in trip ID %s!\n", username, tripID)
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

	query := "SELECT * FROM Trips WHERE ID = ?"
	var t Trips
	var startTravelTimeStr string
	err = db.QueryRow(query, tripID).Scan(
		&t.ID, &t.PickUpLocation, &t.AltPickUpLocation, &startTravelTimeStr, &t.DestinationLocation,
		&t.PassengerNoLeft, &t.MaxPassengerNo, &t.Passenger1, &t.Passenger2, &t.Passenger3, &t.Status, &t.Publisher,
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

// Helper function to check if a user is already enrolled in the trip
func isUserEnrolled(username string, trip Trips) bool {
	if trip.Passenger1.Valid && trip.Passenger1.String == username {
		return true
	}
	if trip.Passenger2.Valid && trip.Passenger2.String == username {
		return true
	}
	if trip.Passenger3.Valid && trip.Passenger3.String == username {
		return true
	}
	return false
}

// Helper function to update the trip and decrement the available passenger slots
func updateTrip(tripID int, username string, trip Trips) {
	// Update the trip with the new passenger
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		// Handle error
		return
	}
	defer db.Close()

	// Identify the empty slot in the passenger list
	var passengerField string
	if trip.Passenger1.Valid && trip.Passenger2.Valid && trip.Passenger3.Valid {
		// All passenger slots are full
		return
	} else if !trip.Passenger1.Valid {
		passengerField = "Passenger1"
	} else if !trip.Passenger2.Valid {
		passengerField = "Passenger2"
	} else if !trip.Passenger3.Valid {
		passengerField = "Passenger3"
	}

	// Update the trip with the new passenger
	query := fmt.Sprintf("UPDATE Trips SET %s=?, PassengerNoLeft=? WHERE ID=?", passengerField)
	_, err = db.Exec(query, username, trip.PassengerNoLeft-1, tripID)
	if err != nil {
		// Handle error
		fmt.Printf("Error updating trip: %v\n", err)
		return
	}
}

// for Car Owner only

// Custom unmarshal function for LicenseNo
func (t *Trips) UnmarshalJSON(data []byte) error {
	type Alias Trips
	aux := &struct {
		AltPickUpLocation string `json:"Alternate Pick-Up Location"`
		Passenger1        string `json:"Passenger 1"`
		Passenger2        string `json:"Passenger 2"`
		Passenger3        string `json:"Passenger 3"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.AltPickUpLocation = sql.NullString{String: aux.AltPickUpLocation, Valid: aux.AltPickUpLocation != ""}
	t.Passenger1 = sql.NullString{String: aux.Passenger1, Valid: aux.Passenger1 != ""}
	t.Passenger2 = sql.NullString{String: aux.Passenger2, Valid: aux.Passenger2 != ""}
	t.Passenger3 = sql.NullString{String: aux.Passenger3, Valid: aux.Passenger3 != ""}
	return nil
}

// curl -X POST http://localhost:5000/api/v1/publishtrip/naruto55 -d "{\"Pick-Up Location\":\"Boon Lay Road\",\"Alternate Pick-Up Location\":\"\",\"Start Traveling Time\":\"2023-11-13T10:30:00Z\",\"Destination Location\":\"Bukit Timah Road\",\"Number of Passengers Allowed\":3}"
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
	// Set the Publisher field
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
	newTrip.ID = int(newID)

	// Return success response with the new trip
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTrip)
}

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
	query := "SELECT MaxPassengerNo, PassengerNoLeft FROM Trips WHERE ID = ?"
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

	query = "UPDATE Trips SET Status=? WHERE ID = ?"
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
	err = db.QueryRow("SELECT StartTravelTime FROM Trips WHERE ID = ?", tripID).Scan(&dbStartTime)
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

	// Print the adjusted time difference in minutes
	fmt.Println("Adjusted Time Difference:", timeDifference)

	// Check if the trip can be canceled (more than or equal to 30 minutes before the start time)
	if timeDifference.Minutes() >= -30 && timeDifference.Minutes() <= 30 {
		fmt.Println("Cancellation window is within 30 minutes. Cannot cancel the trip.")
		http.Error(w, "Trip cannot be canceled. Cancellation window is within 30 minutes", http.StatusBadRequest)
		return
	}
	// Delete the trip from the Trips table
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
	query := "DELETE FROM Trips WHERE ID = ?"
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
