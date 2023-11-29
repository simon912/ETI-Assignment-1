// home.go - the page that the user is redirected to from index.go
package home

import (
	"eti-assignment-1/frontend/profile"
	"net/http"
	"text/template"
)

const homeTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Car Pooling Trip - Home</title>
	<h1>Welcome to the Car Pooling Trip Platform!</h1>
	<h2>Welcome <span id="usergroupspan"></span> <span id="firstnamespan"></span> <span id="lastnamespan"></span></h2>
	<style>
    #message {
        display: none;
    }
	#updateUserContainer {
		display: none;
	}
	#changeCarOwnerContainer {
		display: none;
	}
	#viewProfileContainer {
		display: none;
	}
</style>
 <script>
 		const urlParams = new URLSearchParams(window.location.search);
		const username = urlParams.get('username');
		window.onload = function () {
			retrieveUserData();
		}
        function showMessage(message) {
            const messageContainer = document.getElementById('message');
            messageContainer.innerHTML = message;
            messageContainer.style.display = 'block';

            setTimeout(() => {
                messageContainer.style.display = 'none';
            }, 3000);
        }
		function showProfileContainer() {
			document.getElementById('viewProfileContainer').style.display = 'block';
		}
		function logOutUser() {
			fetch('http://localhost:5000/api/v1/logout', {
				method: 'GET',
				headers: {
					'Content-Type': 'application/json',
				},
			})
			.then(response => {
				if (response.ok) {
					// Redirect to the login page
					window.location.href = 'http://localhost:5001/index';
				} else {
					throw new Error('Logout failed');
				}
			})
			.catch(error => {
				// Display an error message or handle the error appropriately
				showMessage('Logout failed.');
			});
		}
		function retrieveUserData() {
			
		
			// Send a GET request to /api/v1/user/{username} endpoint for user data
			fetch('http://localhost:5000/api/v1/user/' + username, {
				method: 'GET',
				headers: {
					'Content-Type': 'application/json',
				},
			})
			.then(response => {
				if (response.ok) {
					return response.json();
				} else {
					throw new Error('Failed to retrieve user data');
				}
			})
			.then(data => {
				// Assuming the backend sends JSON data with user details
				console.log('Response Data:', data);
				// Extract usergroup, firstname, and lastname from the data
				const usergroup = data['User Group'];
				const firstname = data['First Name'];
				const lastname = data['Last Name'];
				
				const userGroupSpan = document.getElementById('usergroupspan');
            	userGroupSpan.innerHTML = usergroup;
				const firstNameSpan = document.getElementById('firstnamespan');
            	firstNameSpan.innerHTML = firstname;
				const lastNameSpan = document.getElementById('lastnamespan');
            	lastNameSpan.innerHTML = lastname;
			})
			.catch(error => {
				// Display an error message or handle the error appropriately
				showMessage('Failed to retrieve user data.');
			});
		}
		// For Profile
		function showUserInfo() {
			// Send a GET request to /api/v1/user/{username} endpoint for user data
			fetch('http://localhost:5000/api/v1/user/' + username, {
				method: 'GET',
				headers: {
					'Content-Type': 'application/json',
				},
			})
			.then(response => {
				if (response.ok) {
					return response.json();
				} else {
					throw new Error('Failed to retrieve user data');
				}
			})
			.then(data => {
				// Assuming the backend sends JSON data with user details
				console.log('Response Data:', data);
				// Extract usergroup, firstname, and lastname from the data
				const usergroup = data['User Group'];
				const firstname = data['First Name'];
				const lastname = data['Last Name'];
				const mobileno = data['Mobile Number'];
				const emailaddr = data['Email Address'];
				document.getElementById('updateUsername').value = username;
				document.getElementById('updateMobileNo').value = mobileno;
				document.getElementById('updateEmailAddr').value = emailaddr;
			})
			.catch(error => {
				// Display an error message or handle the error appropriately
				showMessage('Failed to retrieve user data.');
			});
			const updateUserContainer = document.getElementById('updateUserContainer');
			const changeCarOwnerContainer = document.getElementById('changeCarOwnerContainer');
			updateUserContainer.style.display = 'block';
			changeCarOwnerContainer.style.display = 'none';
            document.getElementById('updateUserForm').style.display = 'block';
		}
		function showChangeCarOwner() {
			const updateUserContainer = document.getElementById('updateUserContainer');
			const changeCarOwnerContainer = document.getElementById('changeCarOwnerContainer');
			updateUserContainer.style.display = 'none';
			changeCarOwnerContainer.style.display = 'block';

			document.getElementById('carownerLicenseNo').value = '';
			document.getElementById('carownerCarPlateNo').value = '';
            document.getElementById('changeCarOwnerForm').style.display = 'block';
		}
		function updateUserInfo(username) {
            const mobilenumber = parseInt(document.getElementById('updateMobileNo').value);
            const emailaddr = document.getElementById('updateEmailAddr').value;
        
            fetch('http://localhost:5000/api/v1/updateuser/' + username, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    "Mobile Number": mobilenumber,
                    "Email Address": emailaddr,
                }),
            })
            .then(response => {
                if (response.ok) {
                    document.getElementById('updateUserForm').style.display = 'none';
                    document.getElementById('message').style.display = 'block';
                    document.getElementById('message').textContent = "User " + username + "'s info updated successfully";
                } else {
                    throw new Error('User update failed');
                }
            })
            .catch(error => {
                console.error('Error updating user:', error.message);
            });
        }
		function changeCarOwner(username) {
            const licenseno = parseInt(document.getElementById('carownerLicenseNo').value);
            const plateno = document.getElementById('carownerCarPlateNo').value;
            fetch('http://localhost:5000/api/v1/changecarowner/' + username, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    "License Number": licenseno,
                    "Plate Number": plateno,
                }),
            })
            .then(response => {
                if (response.ok) {
                    document.getElementById('changeCarOwnerForm').style.display = 'none';
                    document.getElementById('message').style.display = 'block';
                    document.getElementById('message').textContent = "User " + username + " changed to Car Owner";
                } else {
                    throw new Error('User update failed');
                }
            })
            .catch(error => {
                console.error('Error updating user:', error.message);
            });
        }
	</script>
</head>
<body>
	<div id="userButtonList">
	<button type="button" onclick="showProfileContainer()">View Profile</button>
	<button type="button" onclick="redirectedToProfile()">View Trips</button>
	<button type="button" onclick="logOutUser()">Log Out</button>
	</div>
	<div id="viewProfileContainer">
		<button type="button" onclick="showChangeCarOwner()">Change to Car Owner</button>
		<button type="button" onclick="showUserInfo()">Update Profile</button>
		<div id="changeCarOwnerContainer">
			<form id="changeCarOwnerForm">
				<label for="carownerLicenseNo">Your Driver's License Number:</label>
				<input type="text" id="carownerLicenseNo" required><br>
				<label for="carownerCarPlateNo">Your Car Plate Number:</label>
				<input type="text" id="carownerCarPlateNo" required><br>
				<button type="button" onclick="changeCarOwner(username)">Change to Car Owner</button>
			</form>
		</div>
		<div id="updateUserContainer">
			<form id="updateUserForm">
				<label for="updateUsername">Username:</label>
				<input type="text" id="updateUsername" readonly><br>
				<label for="updateMobileNo">Mobile Number:</label>
				<input type="text" id="updateMobileNo" required><br>
				<label for="updateEmailAddr">Email Address:</label>
				<input type="text" id="updateEmailAddr" required><br>
				<button type="button" onclick="updateUserInfo(username)">Update</button>
			</form>
	</div>
    <div id="message"></div>
</body>
</html>
`

// ProfileHandler handles requests to the profile page
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	// Extract username from the URL
	username := r.URL.Query().Get("username")

	// Pass the username and additional user details to the template
	data := struct {
		Username string
	}{
		Username: username,
	}

	tmpl, err := template.New("home").Parse(homeTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func init() {
	http.HandleFunc("/redirect-profile", func(w http.ResponseWriter, r *http.Request) {
		// You may need to import the "profile" package if it's not already imported
		profile.ProfileHandler(w, r)
	})
}
