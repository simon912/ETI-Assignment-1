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
			listAllCourses()
		case "2":
			createNewCourse()
		case "0":
			fmt.Println("Exiting the program...")
			return
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
	fmt.Println("0. Quit")
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func listAllCourses() {
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

func createNewCourse() {
	courseID := readInput("Enter Course ID: ")

	name := readInput("Enter Course Name: ")
	plannedIntake := readInput("Enter Planned Intake: ")
	minGPA := readInput("Enter Min GPA: ")
	maxGPA := readInput("Enter Max GPA: ")

	IntplannedIntake, _ := strconv.Atoi(plannedIntake)
	IntminGPA, _ := strconv.Atoi(minGPA)
	IntmaxGPA, _ := strconv.Atoi(maxGPA)

	newCourse := map[string]interface{}{
		"Name":           name,
		"Planned Intake": IntplannedIntake,
		"Min GPA":        IntminGPA,
		"Max GPA":        IntmaxGPA,
	}

	postData, err := json.Marshal(newCourse)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	response, err := http.Post(baseURL+"/"+courseID, "application/json", bytes.NewBuffer(postData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusAccepted {
		fmt.Println("Course created successfully.")
	} else if response.StatusCode == http.StatusConflict {
		fmt.Printf("Error - course %s exists\n", courseID)
	} else {
		fmt.Println("Error:", response.Status)
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
