package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type userAttribute struct {
	Usergroup    string `json:"User Group"`
	Firstname    string `json:"First Name"`
	Lastname     string `json:"Last Name"`
	MobileNumber int    `json:"Mobile Number"`
	EmailAddr    string `json:"Email Address"`
	//This attribute will only be used if the User's User Group is Car Owner
	LicenseNo *int    `json:"License Number,omitempty"`
	PlateNo   *string `json:"Car Plate,omitempty"`
}

var user = map[string]userAttribute{
	"john123": {"Passenger", "John", "Doe", 98765432, "john123@gmail.com", nil, nil},
	"jane456": {"Car Owner", "Jane", "Doe", 98534243, "janedoe@gmail.com", intPtr(103436331), strPtr("SKW22G")},
	"lee44":   {"Passenger", "Bryan", "Lee", 95732952, "bryan@gmail.com", nil, nil},
	"tjm95":   {"Car Owner", "Jun Ming", "Tan", 98643435, "tjm@gmail.com", intPtr(104953432), strPtr("SLT45G")},
}

type carPoolingTripAttr struct {
	PickUpLocation      string    `json:"Pick-Up Location"`
	AltPickUpLocation   *string   `json:"Alternate Pick-Up Location"` // can be nullable
	StartTravelTime     time.Time `json:"Start Traveling Time"`
	DestinationLocation string    `json:"Destination Location"`
	NoPassengers        int       `json:"Number of Passengers Allowed"`
	Author              string    `json:"Published By"`
}

var carPoolingTrip = map[int]carPoolingTripAttr{
	1: {"Ang Mo Kio Road", nil, time.Date(2023, time.November, 13, 10, 30, 0, 0, time.UTC), "Geylang Road", 3, "jane456"},
	2: {"Bukit Panjang Ring Road", strPtr("Bangkit Road"), time.Date(2023, time.November, 11, 15, 00, 0, 0, time.UTC), "Choa Chu Kang Road", 3, "tjm456"},
}

// helper functions to create pointers for int and string values
func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}

// Register REST endpoint
func main() {
	router := mux.NewRouter()
	// Endpoint for User
	//This GET method retrieves the relevant course information.
	router.HandleFunc("/api/v1/user/{username}", GetUser).Methods("GET")
	//This POST method creates or updates a user
	router.HandleFunc("/api/v1/user", GetAllUser).Methods("GET")
	//curl http://localhost:5000/api/v1/user/naruto55 -X POST -d "{\"User Group\":\"Car Owner\", \"First Name\":\"Naruto\", \"Last Name\":\"Uzumaki\", \"Mobile Number\":99987634, \"Email Address\":\"naruto@gmail.com\"}"
	router.HandleFunc("/api/v1/user/{username}", CreateUser).Methods("POST")
	router.HandleFunc("/api/v1/user/{username}", UpdateUser).Methods("PUT")
	// Endpoint for Car-Pooling Trips
	router.HandleFunc("/api/v1/carpoolingtrip", GetAllTrip).Methods("GET")
	router.HandleFunc("/api/v1/carpoolingtrip/{tripid}", PublishTrip).Methods("POST")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

// ----------------------------- Endpoint for User ----------------------------------------
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	user, found := user[username]
	if r.Method == http.MethodGet {
		if !found {
			http.Error(w, "Username does not exist", http.StatusNotFound)
			return
		}
		if user.Usergroup == "Passenger" {
			fmt.Fprintf(w, "Username: %s\nUser Group: %s\nFirst Name: %s\nLast Name: %s\nMobile Number: %d\nEmail Address: %s\n\n", username, user.Usergroup, user.Firstname, user.Lastname, user.MobileNumber, user.EmailAddr)
		} else if user.Usergroup == "Car Owner" {
			fmt.Fprintf(w, "Username: %s\nUser Group: %s\nFirst Name: %s\nLast Name: %s\nMobile Number: %d\nEmail Address: %s\nLicense Number: %d\nPlate Number: %s\n\n", username, user.Usergroup, user.Firstname, user.Lastname, user.MobileNumber, user.EmailAddr, *user.LicenseNo, *user.PlateNo)
		}
	}
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {

	//test case for retrieve all: curl -X GET http://localhost:5000/api/v1/user
	for username, user := range user {
		if user.Usergroup == "Passenger" {
			fmt.Fprintf(w, "Username: %s\nUser Group: %s\nFirst Name: %s\nLast Name: %s\nMobile Number: %d\nEmail Address: %s\n\n", username, user.Usergroup, user.Firstname, user.Lastname, user.MobileNumber, user.EmailAddr)
		} else if user.Usergroup == "Car Owner" {
			fmt.Fprintf(w, "Username: %s\nUser Group: %s\nFirst Name: %s\nLast Name: %s\nMobile Number: %d\nEmail Address: %s\nLicense Number: %d\nPlate Number: %s\n\n", username, user.Usergroup, user.Firstname, user.Lastname, user.MobileNumber, user.EmailAddr, *user.LicenseNo, *user.PlateNo)
		}
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	_, found := user[username]
	if found {
		http.Error(w, "Username already exists!", http.StatusConflict)
		return
	}

	var newUser userAttribute
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user[username] = newUser

	// status code 201 - Created
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "User %s has been registered!\n", username)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	// Check if the username exists
	_, found := user[username]
	if !found {
		http.Error(w, "Username does not exist", http.StatusNotFound)
		return
	}
	// Decode the incoming JSON data to update the course
	var updatedUser userAttribute
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user[username] = updatedUser
	// Status Code 202 - Accepted
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "User info has been updated\n")
}

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
