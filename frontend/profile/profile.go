package profile

import (
	"net/http"
	"text/template"
)

const profileTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Car Pooling Trip - Profile</title>
	<h1>Welcome to your profile</h1>
	<h2>Welcome, User {{.Username}}</h2>
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
</style>
 <script>
 		const urlParams = new URLSearchParams(window.location.search);
		const username = urlParams.get('username');
        function showMessage(message) {
            const messageContainer = document.getElementById('message');
            messageContainer.innerHTML = message;
            messageContainer.style.display = 'block';

            setTimeout(() => {
                messageContainer.style.display = 'none';
            }, 3000);
        }
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

			document.getElementById('licenseNo').value = '';
			document.getElementById('carPlateNo').value = '';
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
	</script>
</head>
<body>
	<button type="button" onclick="showChangeCarOwner()">Change to Car Owner</button>
	<button type="button" onclick="showUserInfo()">Update Profile</button>
	<button type="button" onclick="logOutUser()">Log Out</button>
	<div id="changeCarOwnerContainer">
		<form id="changeCarOwnerForm">
			<label for="licenseNo">Your Driver's License Number:</label>
			<input type="text" id="licenseNo" required><br>
			<label for="carPlateNo">Your Car Plate Number:</label>
			<input type="text" id="carPlateNo" required><br>
			<button type="button" onclick="showUserInfo()">Change to Car Owner</button>
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

// HomeHandler handles requests to the home page
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract username from the URL
	username := r.URL.Query().Get("username")

	// Pass the username to the template
	data := struct {
		Username string
	}{
		Username: username,
	}

	tmpl, err := template.New("profile").Parse(profileTemplate)
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

/*
function showUserInfo() {
			const updateUserContainer = document.getElementById('updateUserContainer');
			const changeCarOwnerContainer = document.getElementById('changeCarOwnerContainer');
			updateUserContainer.style.display = 'block';
			changeCarOwnerContainer.style.display = 'none';

			document.getElementById('updateUsername').value = '';
			document.getElementById('updatePassword').value = '';
			document.getElementById('updateMobileNo').value = '';
			document.getElementById('updateEmailAddr').value = '';
            document.getElementById('updateUserForm').style.display = 'block';
		}
		function showChangeCarOwner() {
			const updateUserContainer = document.getElementById('updateUserContainer');
			const changeCarOwnerContainer = document.getElementById('changeCarOwnerContainer');
			updateUserContainer.style.display = 'none';
			changeCarOwnerContainer.style.display = 'block';

			document.getElementById('licenseNo').value = '';
			document.getElementById('carPlateNo').value = '';
            document.getElementById('changeCarOwnerForm').style.display = 'block';
		}
*/
