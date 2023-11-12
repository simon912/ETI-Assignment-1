package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const baseURL = "http://localhost:5000/api/v1/user"

func main() {
	for {
		printMenu()
		option := readInput("Enter an option: ")
		switch option {
		case "1":
			loginUser()
		case "2":
			createNewUser()
		case "0":
			fmt.Println("Exiting the program...")
			return
		case "9":
			listAllUser()
		default:
			fmt.Println("Invalid option, please try again")
		}
	}
}

func printMenu() {
	fmt.Println("=================")
	fmt.Println("Welcome to the Commnity Car-Pooling Platform")
	fmt.Println("1. Login")
	fmt.Println("2. Register")
	fmt.Println("9. Get All User (For Testing Purpose)")
	fmt.Println("0. Quit")
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func listAllUser() {
	response, err := http.Get(baseURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Error:", response.Status)
		return
	}

	var coursesBuffer bytes.Buffer
	_, _ = coursesBuffer.ReadFrom(response.Body)
	fmt.Println(coursesBuffer.String())
}
func loginUser() {
	username := readInput("Enter your username: ")
	userExists, err := checkUserExists(username)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if userExists {
		fmt.Printf("Welcome, user %s!.\n", username)
	} else {
		fmt.Printf("User %s does not exist. Please register first.\n", username)
	}
}

// check if user exists
func checkUserExists(username string) (bool, error) {
	response, err := http.Get(baseURL + "/" + username)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return true, nil
	} else if response.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		return false, fmt.Errorf(response.Status)
	}
}

// Register new user
func createNewUser() {
	username := readInput("Enter a username of your choice: ")
	firstname := readInput("Enter your first name: ")
	lastname := readInput("Enter your last name: ")
	mobileno := readInput("Enter your mobile number: ")
	emailaddr := readInput("Enter your email address: ")
	Intmobileno, _ := strconv.Atoi(mobileno)
	newUser := map[string]interface{}{
		"User Group":    "Passenger",
		"First Name":    firstname,
		"Last Name":     lastname,
		"Mobile Number": Intmobileno,
		"Email Address": emailaddr,
		"License No":    nil,
		"PlateNo":       nil,
	}

	postData, err := json.Marshal(newUser)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	response, err := http.Post(baseURL+"/"+username, "application/json", bytes.NewBuffer(postData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusAccepted {
		fmt.Printf("User %s created successfuly\n", username)
	} else if response.StatusCode == http.StatusConflict {
		fmt.Println("This username is already in use!")
	} else {
		fmt.Println("error:", response.Status)
	}

}

/*
func updateCourse() {
	courseID := readInput("Enter Course ID: ")
	name := readInput("Enter Course Name: ")
	plannedIntake := readInput("Enter Planned Intake: ")
	minGPA := readInput("Enter Min GPA: ")
	maxGPA := readInput("Enter Max GPA: ")

	IntplannedIntake, _ := strconv.Atoi(plannedIntake)
	IntminGPA, _ := strconv.Atoi(minGPA)
	IntmaxGPA, _ := strconv.Atoi(maxGPA)

	updatedCourse := map[string]interface{}{
		"Name":           name,
		"Planned Intake": IntplannedIntake,
		"Min GPA":        IntminGPA,
		"Max GPA":        IntmaxGPA,
	}

	putData, err := json.Marshal(updatedCourse)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	request, err := http.NewRequest(http.MethodPut, baseURL+"/"+courseID, bytes.NewBuffer(putData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusAccepted {
		fmt.Println("Course updated successfully.")
	} else if response.StatusCode == http.StatusNotFound {
		fmt.Printf("Error - Course %s does not exist\n", courseID)
	} else {
		fmt.Println("Error:", response.Status)
	}
}
*/
/*
func deleteCourse() {
	courseID := readInput("Enter Course ID: ")

	request, err := http.NewRequest(http.MethodDelete, baseURL+"/"+courseID, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Course deleted successfully.")
	} else {
		fmt.Println("Error:", response.Status)
	}
}
*/
