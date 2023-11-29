// main.go - the backend
package main

import (
	"database/sql"
	"encoding/json"
	"eti-assignment-1/frontend/profile"
	"fmt"
	"log"
	"net/http"
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
	PickUpLocation      string         `json:"Pick-Up Location"`
	AltPickUpLocation   sql.NullString `json:"Alternate Pick-Up Location"` // can be nullable
	StartTravelTime     time.Time      `json:"Start Traveling Time"`
	DestinationLocation string         `json:"Destination Location"`
	NoPassengers        int            `json:"Number of Passengers Allowed"`
	Author              string         `json:"Published By"`
}

// Register REST endpoint
func main() {
	router := mux.NewRouter()
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5001"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})
	handler := corsOptions.Handler(router)

	// Endpoint for User
	//This GET method retrieves the relevant course information.
	// For Login
	router.HandleFunc("/api/v1/login/{username}", GetUser).Methods("GET")
	router.HandleFunc("/api/v1/logout", LogoutHandler).Methods("GET", "PUT", "POST")

	//test case: curl http://localhost:5000/api/v1/user/naruto55 -X POST -d "{\"User Group\":\"Car Owner\", \"First Name\":\"Naruto\", \"Last Name\":\"Uzumaki\", \"Mobile Number\":99987634, \"Email Address\":\"naruto@gmail.com\"}"
	// For Register
	router.HandleFunc("/api/v1/register/{username}", CreateUser).Methods("POST")
	router.HandleFunc("/api/v1/user/{username}", GetUser).Methods("GET")
	router.HandleFunc("/api/v1/updateuser/{username}", UpdateUser).Methods("PUT")
	// test case: curl http://localhost:5000/api/v1/changecarowner/naruto55 -X PUT -d "{\"License Number\": 111123335, \"Plate Number\": \"ABC123\"}"
	router.HandleFunc("/api/v1/changecarowner/{username}", ChangeToCarOwner).Methods("PUT")

	http.HandleFunc("/profile", profile.ProfileHandler)
	http.HandleFunc("/display-profile", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		redirectToProfileSuccess(w, r, username)
	})
	router.HandleFunc("/redirect-profile/", RedirectProfileHandler).Methods("GET")
	// handlefunc here which deletes account after its 1 year old/365 days old

	// Endpoint for Car-Pooling Trips
	//router.HandleFunc("/api/v1/carpoolingtrip", GetAllTrip).Methods("GET")
	//router.HandleFunc("/api/v1/carpoolingtrip/{tripid}", PublishTrip).Methods("POST")
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5001")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNoContent)
	})
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", handler))
}
func redirectToProfileSuccess(w http.ResponseWriter, r *http.Request, username string) {
	http.Redirect(w, r, "http://localhost:5001/profile?username="+username, http.StatusSeeOther)
}
func RedirectProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the username from the query parameters
	username := r.URL.Query().Get("username")

	// Redirect to the profile page
	http.Redirect(w, r, "/profile?username="+username, http.StatusSeeOther)
}

// Helper function to handle null values in SQL parameters
func nullOrValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	return value
}

// Custom unmarshal function for LicenseNo
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

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received logout request")
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5001")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		// Add other CORS headers as needed
		w.WriteHeader(http.StatusOK)
		return
	}

	// Respond to the actual request (GET)
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5001")
	// Add other CORS headers as needed

	// Perform the redirect after the actual request is handled
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

/*
// ----------------------------- Endpoint for Car Pooling Trips ----------------------------------------
func GetAllTrip(w http.ResponseWriter, r *http.Request) {
	//test case for retrieve all: curl -X GET http://localhost:5000/api/v1/carpoolingtrip
	for tripid, carPoolingTrip := range carPoolingTrip {
		startTime := carPoolingTrip.StartTravelTime.Format("2006-01-02 15:04:05 MST")
		altPickUpLocation := ""
		if carPoolingTrip.AltPickUpLocation != nil {
			altPickUpLocation = *carPoolingTrip.AltPickUpLocation
		}
		fmt.Fprintf(w, "Trip ID: %d\nPick-Up Location: %s\nAlternate Pick-Up Location: %s\nStarting Traveling Time: %s\nDestination Location: %s\nNumber of Passengers: %d\nPublished By: %s\n\n", tripid, carPoolingTrip.PickUpLocation, altPickUpLocation, startTime, carPoolingTrip.DestinationLocation, carPoolingTrip.NoPassengers, carPoolingTrip.Author)
	}
}

// for Car Owner only
// curl http://localhost:5000/api/v1/carpoolingtrip/3 -X POST -d "{\"Pick-Up Location\":\"Choa Chu Kang Road\", \"Alternate Pick-Up Location\":\"\", \"Start Traveling Time\":\"2023-11-13T10:30:00Z\", \"Destination Location\":\"Bukit Timah Road\", \"Number of Passengers Allowed\":4, \"Published By\":\"jane456\"}"
func PublishTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tripid := params["tripid"]

	//convert tripid to int
	tripidInt, err := strconv.Atoi(tripid)
	if err != nil {
		http.Error(w, "Invalid trip ID", http.StatusBadRequest)
		return
	}
	_, found := carPoolingTrip[tripidInt]
	if found {
		http.Error(w, "This car pooling trip already exist", http.StatusConflict)
		return
	}

	var newTrip carPoolingTripAttr
	err = json.NewDecoder(r.Body).Decode(&newTrip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	carPoolingTrip[tripidInt] = carPoolingTripAttr(newTrip)

	// status code 201 - Created
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Your car publishing trip ID %s has been registered!\n", tripid)
}
*/
