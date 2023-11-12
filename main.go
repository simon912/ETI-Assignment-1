package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type userAttribute struct {
	Usergroup    string `json:"User Group"`
	Firstname    string `json:"First Name"`
	Lastname     string `json:"Last Name"`
	MobileNumber int    `json:"Mobile Number"`
	EmailAddr    string `json:"Email Address"`
	//This attribute will only be used if the User's User Group is Car Owner
	//LicenseNo int `json:"License Number"`
	//PlateNo string `json:"Car Plate"`
}

var user = map[string]userAttribute{
	"john123": {"Passenger", "John", "Doe", 98765432, "john123@gmail.com"},
	"jane456": {"Car Owner", "Jane", "Doe", 98534243, "janedoe@gmail.com"},
	"lee44":   {"Passenger", "Bryan", "Lee", 95732952, "bryan@gmail.com"},
	"tjm95":   {"Car Owner", "Jun Ming", "Tan", 98643435, "tjm@gmail.com"},
}

// Register REST endpoint
func main() {
	router := mux.NewRouter()
	//This GET method display all car-pooling trip
	//router.HandleFunc("/api/v1/carpooling", GetAllCourses).Methods("GET")
	//This GET & DELETE method retrieves the relevant course information.
	//router.HandleFunc("/api/v1/courses/{course id}", GetOrDelete).Methods("GET", "DELETE")
	//This POST method creates or updates a user
	router.HandleFunc("/api/v1/user", GetAllUser).Methods("GET")
	//curl http://localhost:5000/api/v1/user/naruto55 -X POST -d "{\"User Group\":\"Car Owner\", \"First Name\":\"Naruto\", \"Last Name\":\"Uzumaki\", \"Mobile Number\":99987634, \"Email Address\":\"naruto@gmail.com\"}"
	router.HandleFunc("/api/v1/user/{username}", CreateUser).Methods("POST")
	router.HandleFunc("/api/v1/user/{username}", UpdateUser).Methods("PUT")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
func PrintMenu() {
	fmt.Println("======")
	fmt.Println("Course Management Console")
	fmt.Println("1. List all courses")
	fmt.Println("2. Create new course")
	fmt.Println("3. Update course")
	fmt.Println("4. Delete course")
	fmt.Println("9. Quit")
}
func UserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

/*
	func GetOrDelete(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		courseID := params["course id"]
		course, found := user[courseID]
		if r.Method == http.MethodGet { ///test case: curl -X GET http://localhost:5000/api/v1/courses/IT
			if !found {
				http.Error(w, "Invalid Course ID", http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "Course Name: %s\nPlanned Intake: %d\nMin GPA: %d\nMax GPA: %d", course.Name, course.PlannedIntake, course.MinGPA, course.MaxGPA)
		} else if r.Method == http.MethodDelete { //test case: curl -X DELETE http://localhost:5000/api/v1/courses/IT
			if found {
				delete(user, courseID)
				fmt.Fprintf(w, "%s Deleted", courseID)
			} else {
				http.Error(w, "Invalid Course ID", http.StatusNotFound)
			}
		}
	}
*/

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	/*query := r.URL.Query().Get("q")
	gpaStr := r.URL.Query().Get("gpa")
	var gpa int
	if gpaStr != "" {
		var err error
		gpa, err = strconv.Atoi(gpaStr)
		if err != nil {
			http.Error(w, "Invalid GPA format", http.StatusBadRequest)
			return
		}
	}*/
	//Retrieve all courses when there is no query string
	//test case for retrieve all: curl -X GET http://localhost:5000/api/v1/user
	for username, user := range user {
		//Return the search results when there is a partial match in the Name
		fmt.Fprintf(w, "Username: %s\nUser Group: %s\nFirst Name: %s\nLast Name: %s\nMobile Number: %d\nEmail Address: %s\n\n", username, user.Usergroup, user.Firstname, user.Lastname, user.MobileNumber, user.EmailAddr)
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
	courseID := params["username"]

	// Check if the course ID exists
	_, found := user[courseID]
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

	user[courseID] = updatedUser

	// Status Code 202 - Accepted
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "User info has been updated\n")
}
