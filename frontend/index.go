package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const userURL = "http://localhost:5000/api/v1/user"
const tripURL = "http://localhost:5000/api/v1/carpoolingtrip"

func main() {
	var currentUsername string
	for {
		printMenu()
		option := readInput("Enter an option: ")
		switch option {
		case "1":
			success, username, userGroup := loginUser()
			if success {
				currentUsername = username
				if userGroup == "Passenger" {
					printPassengerMenu(currentUsername)
				} else if userGroup == "Car Owner" {
					printCarOwnerMenu(currentUsername)
				}
			}
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
	fmt.Println("Welcome to the Community Car-Pooling Platform")
	fmt.Println("1. Login")
	fmt.Println("2. Register")
	fmt.Println("9. Get All User (For Testing Purpose)")
	fmt.Println("0. Quit")
}

func printPassengerMenu(username string) {
	fmt.Println("=================")
	fmt.Println("Welcome Passenger")
	fmt.Println("1. Change to Car Owner")
	fmt.Println("2. Update Profile")
	fmt.Println("3. Browse Car-Pooling Trips")
	fmt.Println("4. Delete Profile")
	fmt.Println("0. Logout")
	passengerOption := readInput("Enter an option: ")
	switch passengerOption {
	case "1": // Change to Car Owner
		changeToCarOwner(username)
		return
	case "2": // Update Profile
		fmt.Println("Option 2 selected")
	case "3": // Browse Car Pooling Trip
		listAllTrip()
		printPassengerMenu(username)
	case "4": // Delete Profile if user is over a year old
		fmt.Println("Option 4 selected")
	case "0":
		fmt.Println("Logging out...")
		return
	default:
		fmt.Println("Invalid option, please try again")
	}
}

func printCarOwnerMenu(username string) {
	fmt.Println("=================")
	fmt.Println("Welcome Car Owner")
	fmt.Println("1. Publish Car-Pooling Trips")
	fmt.Println("2. Manage Car-Pooling Trips")
	fmt.Println("0. Logout")
	passengerOption := readInput("Enter an option: ")
	switch passengerOption {
	case "1":
		fmt.Println("Option 1 selected")
	case "2":
		fmt.Println("Option 2 selected")
	case "0":
		fmt.Println("Logging out...")
		return
	default:
		fmt.Println("Invalid option, please try again")
	}
}
func readInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// ------------------------------- User ---------------------------------
func listAllUser() {
	response, err := http.Get(userURL)
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

func loginUser() (bool, string, string) {
	username := readInput("Enter your username: ")
	userExists, userGroup, err := checkUserExists(username)
	if err != nil {
		fmt.Println("Error:", err)
		return false, "", ""
	}
	if userExists {
		fmt.Printf("Welcome, user %s!\n", username)
		return true, username, userGroup
	} else {
		fmt.Printf("User %s does not exist. Please register first.\n", username)
		return false, "", ""
	}
}

// check if user exists
func checkUserExists(username string) (bool, string, error) {
	response, err := http.Get(userURL + "/" + username)
	if err != nil {
		return false, "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, "", err
	}
	if response.StatusCode == http.StatusOK {
		userGroupIndex := strings.Index(string(body), "User Group:")
		if userGroupIndex != -1 {
			userGroupStart := userGroupIndex + len("User Group:")
			userGroupEnd := strings.Index(string(body)[userGroupStart:], "\n")
			if userGroupEnd == -1 {
				return false, "", fmt.Errorf("user group not found in response")
			}

			userGroup := strings.TrimSpace(string(body)[userGroupStart : userGroupStart+userGroupEnd])
			return true, userGroup, nil
		} else {
			return false, "", fmt.Errorf("user group not found in response")
		}
	} else if response.StatusCode == http.StatusNotFound {
		return false, "", nil
	} else {
		return false, "", fmt.Errorf("error: %s", response.Status)
	}
}

// Register new user
func createNewUser() {
	username := readInput("Enter your username: ")
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

	response, err := http.Post(userURL+"/"+username, "application/json", bytes.NewBuffer(postData))
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

func changeToCarOwner(username string) {

	licenseNo := readInput("Enter License Number: ")
	plateNo := readInput("Enter Plate Number: ")

	IntlicenseNo, _ := strconv.Atoi(licenseNo)

	updatedUser := map[string]interface{}{
		"User Group":     "Car Owner",
		"License Number": IntlicenseNo,
		"Plate Number":   plateNo,
	}

	putData, err := json.Marshal(updatedUser)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	request, err := http.NewRequest(http.MethodPut, userURL+"/"+username+"/changecarowner", bytes.NewBuffer(putData))
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
		fmt.Println("User's User Group changed to Car Owner.")
	} else if response.StatusCode == http.StatusNotFound {
		fmt.Printf("Error - User %s does not exist\n", username)
	} else {
		fmt.Println("Error:", response.Status)
	}
}

// --------------------------- Car Pooling Trip ---------------------------
func listAllTrip() {
	response, err := http.Get(tripURL)
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
