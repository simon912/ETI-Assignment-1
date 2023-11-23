package main

import (
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Car Pooling Trip - Login</title>
	<h1>Welcome to the Car Pooling Trip Platform!</h1>
	<h2>You need to login or register to participate in the car pooling service</h2>
	<style>
    #loginFormContainer {
        display: none;
    }
    #registerFormContainer {
        display: none;
    }
    #message {
        display: none;
    }
</style>
 <script>
        function showMessage(message) {
            const messageContainer = document.getElementById('message');
            messageContainer.innerHTML = message;
            messageContainer.style.display = 'block';

            setTimeout(() => {
                messageContainer.style.display = 'none';
            }, 3000);
        }
        function showLoginForm() {
            document.getElementById('loginFormContainer').style.display = 'block';
            document.getElementById('registerFormContainer').style.display = 'none';
            document.getElementById('message').style.display = 'none';
            //reset the form field
            document.getElementById('loginUsername').value = '';
            document.getElementById('loginPassword').value = '';
            document.getElementById('loginForm').style.display = 'block';
        }

        function showRegisterForm() {
            document.getElementById('loginFormContainer').style.display = 'none';
            document.getElementById('registerFormContainer').style.display = 'block';
            document.getElementById('message').style.display = 'none';
            //reset the form field
            document.getElementById('registerUsername').value = '';
            document.getElementById('registerPassword').value = '';
            document.getElementById('registerFirstname').value = '';
            document.getElementById('registerLastname').value = '';
			document.getElementById('registerMobilenumber').value = '';
			document.getElementById('registerEmail').value = '';
            document.getElementById('registerForm').style.display = 'block';
        }
		function loginUser() {
			const username = document.getElementById('loginUsername').value;
    		const password = document.getElementById('loginPassword').value;

    		// Send a GET request to your /api/v1/login/{username} endpoint for authentication
    		fetch('http://localhost:5000/api/v1/login/' + username, {
        		method: 'GET',
        		headers: {
            		'Content-Type': 'application/json',
        	},
    	})
        	.then(response => {
            	if (response.ok) {
                	return response.json(); // Assuming your backend sends JSON data on successful authentication
            	} else {
                	throw new Error('Login failed');
            	}
        	})
        	.then(data => {
            	// Check if the password from the response matches the entered password
            	if (data && data.Password === password) {
                	// Display a success message or redirect to another page if needed
                	showMessage('Login successful');
            	}
        	})
        	.catch(error => {
            	// Display an error message or handle the error appropriately
            	showMessage('Login failed. Invalid username or password.');
        	});
		}
		function registerUser() {
			// Get values from the registration form
			const username = document.getElementById('registerUsername').value;
			const password = document.getElementById('registerPassword').value;
			const firstname = document.getElementById('registerFirstname').value;
			const lastname = document.getElementById('registerLastname').value;
			const mobileNumber = document.getElementById('registerMobilenumber').value;
			const email = document.getElementById('registerEmail').value;
	
			// You can perform any additional client-side validation here if needed
	
			// Send a POST request to your /api/v1/user/{username} endpoint with the registration data
			fetch("'/api/v1/register/' + username", {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					Username: username,
					Password: password,
					Firstname: firstname,
					Lastname: lastname,
					MobileNumber: parseInt(mobileNumber), // assuming MobileNumber is an integer
					EmailAddr: email,
				}),
			})
				.then(response => {
					if (response.ok) {
						return response.text();
					} else {
						throw new Error('Registration failed');
					}
				})
				.then(data => {
					// Display a success message or redirect to another page if needed
					showMessage('Register successful');
				})
				.catch(error => {
					// Display an error message or handle the error appropriately
					showMessage('Registration failed. Username may already be in use.');
				});
		}
	</script>
	
	</script>
</head>
<body>
	<button type="button" onclick="showLoginForm()">Login</button>
	<button type="button" onclick="showRegisterForm()">Register</button>
    <div id="courseListContainer">
        <ul id="courseList"></ul>
    </div>

    <div id="loginFormContainer">
        <form id="loginForm">
            <label for="loginUsername">Username:</label>
            <input type="text" id="loginUsername" required><br>
            <label for="loginPassword">Password:</label>
            <input type="text" id="loginPassword" required><br>
			<button type="button" onclick="loginUser()">Login</button>
        </form>
    </div>
    <div id="registerFormContainer">
        <form id="registerForm">
            <label for="registerUsername">Username:</label>
            <input type="text" id="registerUsername" required><br>
            <label for="registerPassword">Password:</label>
            <input type="text" id="registerPassword" required><br>
            <label for="registerFirstname">First Name:</label>
            <input type="text" id="registerFirstname" required><br>
            <label for="registerLastname">Last Name:</label>
            <input type="text" id="registerLastname" required><br>
            <label for="registerMobilenumber">Mobile Number:</label>
            <input type="number" id="registerMobilenumber" required><br>
			<label for="registerEmail">Email Address:</label>
            <input type="text" id="registerEmail" required><br>
            <button type="button" onclick="registerUser()">Register</button>
        </form>
    </div>
    <div id="message"></div>
</body>
</html>
`

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("login").Parse(htmlTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Listening at http://localhost:5001")
	http.ListenAndServe(":5001", nil)
}
