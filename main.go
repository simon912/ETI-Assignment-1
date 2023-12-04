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

type Users struct {
	Username     string `json:"Username"`
	Password     string `json:"Password"`
	Usergroup    string `json:"User Group"`
	Firstname    string `json:"First Name"`
	Lastname     string `json:"Last Name"`
	MobileNumber int    `json:"Mobile Number"`
	EmailAddr    string `json:"Email Address"`
	//This attribute will only be used if the User's User Group is Car Owner
	LicenseNo    sql.NullInt64  `json:"License Number,omitempty"`
	PlateNo      sql.NullString `json:"Plate Number,omitempty"`
	CreationDate string         `json:"Account Creation Date"`
}

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
		AllowedOrigins: []string{"http://localhost:5001"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})
	handler := corsOptions.Handler(router)

	// Endpoint for User
	//This GET method retrieves the relevant course information.
	// For Login
	router.HandleFunc("/api/v1/login/{username}", GetUser).Methods("GET")
	//test case: curl http://localhost:5000/api/v1/user/naruto55 -X POST -d "{\"User Group\":\"Car Owner\", \"First Name\":\"Naruto\", \"Last Name\":\"Uzumaki\", \"Mobile Number\":99987634, \"Email Address\":\"naruto@gmail.com\"}"
	// For Register
	router.HandleFunc("/api/v1/register/{username}", CreateUser).Methods("POST")
	router.HandleFunc("/api/v1/user/{username}", GetUser).Methods("GET")
	router.HandleFunc("/api/v1/updateuser/{username}", UpdateUser).Methods("PUT")
	// test case: curl http://localhost:5000/api/v1/changecarowner/naruto55 -X PUT -d "{\"License Number\": 111123335, \"Plate Number\": \"ABC123\"}"
	router.HandleFunc("/api/v1/changecarowner/{username}", ChangeToCarOwner).Methods("PUT")
	// For Deletion of User, which is called properly if its over 1 year old/365 days old
	router.HandleFunc("/api/v1/delete/{username}", DeleteUser).Methods("DELETE")
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
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", handler))
}

// Helper function to handle NULL values in SQL parameters
func nullOrValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	return value
}

// Custom unmarshal function for LicenseNumber & PlateNumber
func (u *Users) UnmarshalJSON(data []byte) error {
	type Alias Users
	aux := &struct {
		LicenseNumber int    `json:"License Number"`
		PlateNumber   string `json:"Plate Number"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	u.LicenseNo = sql.NullInt64{Int64: int64(aux.LicenseNumber), Valid: true}
	u.PlateNo = sql.NullString{String: aux.PlateNumber, Valid: aux.PlateNumber != ""}
	return nil
}

// ----------------------------- Endpoint for User ----------------------------------------

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query := "SELECT * FROM Users WHERE Username = ?"
	result, err := db.Query(query, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer result.Close()

	found := false
	for result.Next() {
		var u Users
		err = result.Scan(&u.Username, &u.Password, &u.Usergroup, &u.Firstname, &u.Lastname, &u.MobileNumber, &u.EmailAddr, &u.LicenseNo, &u.PlateNo, &u.CreationDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := struct {
			Username     string `json:"Username"`
			Password     string `json:"Password"`
			Usergroup    string `json:"User Group"`
			Firstname    string `json:"First Name"`
			Lastname     string `json:"Last Name"`
			MobileNumber int    `json:"Mobile Number"`
			EmailAddr    string `json:"Email Address"`
		}{
			Username:     u.Username,
			Password:     u.Password,
			Usergroup:    u.Usergroup,
			Firstname:    u.Firstname,
			Lastname:     u.Lastname,
			MobileNumber: u.MobileNumber,
			EmailAddr:    u.EmailAddr,
		}

		// Convert the response to JSON and send it
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		found = true
	}
	if !found {
		fmt.Println("User doesn't exist")
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	if userExists(username) {
		http.Error(w, "Username already in use", http.StatusConflict)
		return
	}

	var newUser Users
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	newUser.Username = username
	query := "INSERT INTO Users (Username, Password, UserGroup, FirstName, LastName, MobileNumber, EmailAddress, CreationDateTime) VALUES (?, ?, ?, ?, ?, ?, ?, NOW())"
	_, err = db.Exec(query, newUser.Username, newUser.Password, "Passenger", newUser.Firstname, newUser.Lastname, newUser.MobileNumber, newUser.EmailAddr)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "User %s has been registered!\n", username)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	var updateUser Users
	err := json.NewDecoder(r.Body).Decode(&updateUser)
	if err != nil {
		panic(err.Error())
	}
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	query := "UPDATE Users SET MobileNumber=?, EmailAddress=? WHERE Username=?"
	_, err = db.Exec(query, updateUser.MobileNumber, updateUser.EmailAddr, username)
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
	fmt.Fprintf(w, "User %s has been updated!\n", username)
}

func ChangeToCarOwner(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	var updateUser Users
	fmt.Println("Received JSON:", r.Body)
	err := json.NewDecoder(r.Body).Decode(&updateUser)
	if err != nil {
		panic(err.Error())
	}
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	query := "UPDATE Users SET UserGroup=?, LicenseNo=?, PlateNo=? WHERE Username=?"
	_, err = db.Exec(query, "Car Owner", nullOrValue(updateUser.LicenseNo), nullOrValue(updateUser.PlateNo), username)
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
	fmt.Fprintf(w, "User %s has been changed to Car Owner!\n", username)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Retrieve the creation date of the user from the database as a string
	var creationDateStr string
	query := "SELECT CreationDateTime FROM Users WHERE Username = ?"
	err = db.QueryRow(query, username).Scan(&creationDateStr)
	if err != nil {
		fmt.Printf("Error retrieving user creation date: %v\n", err)
		http.Error(w, "Error retrieving user creation date", http.StatusInternalServerError)
		return
	}

	// Parse the creation date string into a time.Time value
	creationDate, err := time.Parse("2006-01-02", creationDateStr)
	if err != nil {
		fmt.Printf("Error parsing creation date: %v\n", err)
		http.Error(w, "Error parsing creation date", http.StatusInternalServerError)
		return
	}

	// Calculate the age of the user
	age := time.Since(creationDate).Hours() / 24 // Convert to days

	// Check if the user is over 1 year old (365 days)
	if age >= 365 {
		// Check if the user is referenced in the trips table
		tripsReferenced, err := isUserReferencedInTrips(db, username)
		if err != nil {
			fmt.Printf("Error checking references in Trips table: %v\n", err)
			http.Error(w, "Error checking references in Trips table", http.StatusInternalServerError)
			return
		}

		if tripsReferenced {
			http.Error(w, "User cannot be deleted. User has already published Trip(s).", http.StatusBadRequest)
			return
		}

		// If the user is over 1 year old and not referenced in trips, delete the user
		deleteUser(w, username)
	} else {
		// If the user is not over 1 year old, return error
		http.Error(w, "User must be over 1 year old for deletion", http.StatusBadRequest)
	}
}

// Helper function to delete the user
func deleteUser(w http.ResponseWriter, username string) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Check if the user is referenced in the trips table
	tripsReferenced, err := isUserReferencedInTrips(db, username)
	if err != nil {
		fmt.Printf("Error checking references in Trips table: %v\n", err)
		http.Error(w, "Error checking references in Trips table", http.StatusInternalServerError)
		return
	}

	if tripsReferenced {
		http.Error(w, "User cannot be deleted. User has already published Trip(s).", http.StatusBadRequest)
		return
	}

	// Delete the user
	query := "DELETE FROM Users WHERE Username = ?"
	result, err := db.Exec(query, username)
	if err != nil {
		panic(err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}

	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s deleted successfully", username)
}

// Helper function to check if the user is referenced in the Trips table
func isUserReferencedInTrips(db *sql.DB, username string) (bool, error) {
	query := "SELECT COUNT(*) FROM Trips WHERE Publisher = ?"
	var count int
	err := db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ----------------------------- Endpoint for Car Pooling Trips ----------------------------------------
func GetAllTrip(w http.ResponseWriter, r *http.Request) {
	//test case for retrieve all: curl -X GET http://localhost:5000/api/v1/trips
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	// handle error
	if err != nil {
		panic(err.Error())
	}
	// database operation
	defer db.Close()
	results, err := db.Query("SELECT * FROM Trips")
	if err != nil {
		panic(err.Error())
	}
	var trips []Trips
	for results.Next() {
		var t Trips
		var startTravelTimeStr string
		err = results.Scan(&t.ID, &t.PickUpLocation, &t.AltPickUpLocation, &startTravelTimeStr, &t.DestinationLocation, &t.PassengerNoLeft, &t.MaxPassengerNo, &t.Passenger1, &t.Passenger2, &t.Passenger3, &t.Status, &t.Publisher)
		if err != nil {
			panic(err.Error())
		}
		// Parse the time string and set it in the Trips struct
		parsedTime, err := time.Parse("15:04:05", startTravelTimeStr)
		if err != nil {
			// Handle the error, e.g., log it and proceed with a default value
			fmt.Printf("Error parsing time for Trip ID %d: %v\n", t.ID, err)
			parsedTime = time.Time{} // Default to zero time
		}
		t.StartTravelTime = parsedTime.Format(time.RFC3339)
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
	parsedTime, err := time.Parse(time.RFC3339, newTrip.StartTravelTime)
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
	query := "INSERT INTO Trips (PickUpLocation, AltPickUpLocation, StartTravelTime, DestinationLocation, PassengerNoLeft, MaxPassengerNo, Status, Publisher) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
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

	var startTrip Trips
	fmt.Println("Received JSON:", r.Body)

	err := json.NewDecoder(r.Body).Decode(&startTrip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if PassengerNoLeft is equal to MaxPassengerNo
	if startTrip.MaxPassengerNo == startTrip.PassengerNoLeft {
		http.Error(w, "Trip cannot start if no one is enrolled into your trip", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	query := "UPDATE Trips SET Status=? WHERE ID = ?"
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
