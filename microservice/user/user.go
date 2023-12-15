// user.go - the backend for user management, port 5000
package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	LicenseNo    sql.NullInt64  `json:"License Number,omitempty"`
	PlateNo      sql.NullString `json:"Plate Number,omitempty"`
	CreationDate string         `json:"Account Creation Date"`
}

// REST endpoint for User
func main() {
	router := mux.NewRouter()
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})
	handler := corsOptions.Handler(router)

	// GET function for login
	router.HandleFunc("/api/v1/login/{username}", GetUser).Methods("GET")
	// POST function for register
	router.HandleFunc("/api/v1/register/{username}", CreateUser).Methods("POST")
	// GET function for retrieval of user info
	router.HandleFunc("/api/v1/user/{username}", GetUser).Methods("GET")
	// PUT function for updating user
	router.HandleFunc("/api/v1/updateuser/{username}", UpdateUser).Methods("PUT")
	// PUT function for updating user's user group to Car Owner
	router.HandleFunc("/api/v1/changecarowner/{username}", ChangeToCarOwner).Methods("PUT")
	// DELETE function for Deletion of user, which is called only if the user is over 1 year old/365 days old
	router.HandleFunc("/api/v1/delete/{username}", DeleteUser).Methods("DELETE")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", handler))
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
// Function to get that user's info
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

	userFound := false

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
		userFound = true
	}
	if !userFound {
		http.Error(w, `{"error": "User does not exist"}`, http.StatusNotFound)
		return
	}
}

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
// Function to Create User, mainly used for Registration of User
func CreateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	if userExists(username) {
		http.Error(w, `{"error": "Username already in use`, http.StatusConflict)
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

// Function to update the user's information (Email Address, Mobile Number)
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

// Function to change the Passenger to Car Owner if they provide License Number and Plate Number of their car
func ChangeToCarOwner(w http.ResponseWriter, r *http.Request) {
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

	// Check if user is enrolled before
	if isUserReferencedInEnrollment(username) {
		http.Error(w, `{"error": "User is enrolled in a Trip"}`, http.StatusBadRequest)
		return
	}

	query := "UPDATE Users SET UserGroup=?, LicenseNo=?, PlateNo=? WHERE Username=?"
	_, err = db.Exec(query, "Car Owner", updateUser.LicenseNo, updateUser.PlateNo, username)
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

// Deletion of User only if the user is over 1 year old
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening database: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Retrieve the creation date of the user from the database as a string
	var creationDateStr string
	query := "SELECT CreationDateTime FROM Users WHERE Username = ?"
	err = db.QueryRow(query, username).Scan(&creationDateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user creation date: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse the creation date string into a time.Time value
	creationDate, err := time.Parse("2006-01-02", creationDateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing creation date: %v", err), http.StatusInternalServerError)
		return
	}

	// Calculate the age of the user
	age := time.Since(creationDate).Hours() / 24 // Convert to days
	// Check if the user is over 1 year old (365 days)
	if age < 365 {
		http.Error(w, `{"error": "User is not over 1 year old yet"}`, http.StatusBadRequest)
		return
	}
	if isUserReferencedInTrip(username) {
		http.Error(w, `{"error": "User has published a Trip."}`, http.StatusBadRequest)
		return
	}
	if isUserReferencedInEnrollment(username) {
		http.Error(w, `{"error": "User is enrolled in a Trip"}`, http.StatusBadRequest)
		return
	}

	deleteUser(w, db, username)
	// Respond with success message or redirect as needed
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s deleted successfully", username)

}

// Helper function to delete the user
func deleteUser(w http.ResponseWriter, db *sql.DB, username string) error {
	// Delete the user
	query := "DELETE FROM Users WHERE Username = ?"
	result, err := db.Exec(query, username)
	if err != nil {
		// Check for foreign key constraint violation
		if strings.Contains(err.Error(), "foreign key constraint") {
			return errors.New("user cannot be deleted: The user is referenced in other records")
		}
		return fmt.Errorf("error deleting user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Helper function to check if the user has published a trip (Car Owner)
func isUserReferencedInTrip(username string) bool {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		// Handle error
		fmt.Printf("Error opening database: %v\n", err)
	}
	defer db.Close()

	query := "SELECT COUNT(*) FROM Trips WHERE Publisher = ?"
	var count int
	err = db.QueryRow(query, username).Scan(&count)
	if err != nil {
		// Handle error
		fmt.Printf("Error querying database: %v\n", err)
		return false
	}

	return count > 0
}

// Helper function to check if the user is enrolled in a trip (Passenger)
func isUserReferencedInEnrollment(username string) bool {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpoolingtrip")
	if err != nil {
		// Handle error
		fmt.Printf("Error opening database: %v\n", err)
	}
	defer db.Close()

	query := "SELECT COUNT(*) FROM Enrollment WHERE Username = ?"
	var count int
	err = db.QueryRow(query, username).Scan(&count)
	if err != nil {
		// Handle error
		fmt.Printf("Error querying database: %v\n", err)
		return false
	}

	return count > 0
}
