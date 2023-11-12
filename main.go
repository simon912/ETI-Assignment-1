package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	//This GET method display all car-pooling trip
	//router.HandleFunc("/api/v1/carpooling", GetAllCourses).Methods("GET")
	//This GET method retrieves the relevant course information.
	router.HandleFunc("/api/v1/user/{username}", GetUser).Methods("GET")
	//This POST method creates or updates a user
	router.HandleFunc("/api/v1/user", GetAllUser).Methods("GET")
	//curl http://localhost:5000/api/v1/user/naruto55 -X POST -d "{\"User Group\":\"Car Owner\", \"First Name\":\"Naruto\", \"Last Name\":\"Uzumaki\", \"Mobile Number\":99987634, \"Email Address\":\"naruto@gmail.com\"}"
	router.HandleFunc("/api/v1/user/{username}", CreateUser).Methods("POST")
	router.HandleFunc("/api/v1/user/{username}", UpdateUser).Methods("PUT")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	user, found := user[username]
	if r.Method == http.MethodGet { ///test case: curl -X GET http://localhost:5000/api/v1/courses/IT
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
