// carowner_home.go - the page that the user is redirected to from index.go if their user group is car owner
package carowner_home

import (
	"net/http"
	"text/template"
)

const carownerTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Car Pooling Trip - Home</title>
	<h1>Welcome to the Car Pooling Trip Platform!</h1>
	<h2>Welcome <span id="usergroupspan"></span> <span id="firstnamespan"></span> <span id="lastnamespan"></span></h2>
	<h2>Car Owner Redirection Test </h2>
	<style>
    #message {
        display: none;
    }
	#updateUserContainer, #viewProfileContainer, #viewTripContainer, #publishContainer {
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
			document.getElementById('viewTripContainer').style.display = 'none';
			document.getElementById('confirmDeleteContainer').style.display = 'none';
			document.getElementById('publishContainer').style.display = 'none';
		}
		function logOutUser() {
			window.location.href = 'http://localhost:5001/';
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
			document.getElementById('updateUserForm').style.display = 'block';
			document.getElementById('confirmDeleteContainer').style.display = 'none';
			document.getElementById('updateUserContainer').style.display = 'block';
            document.getElementById('updateUserForm').style.display = 'block';
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
		function showDeleteConfirmation() {
			document.getElementById('updateUserContainer').style.display = 'none';
			document.getElementById('confirmDeleteContainer').style.display = 'block';
		}
		function deleteUser(username) {
			fetch('http://localhost:5000/api/v1/delete/' + username, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                },
            })
            .then(response => {
                if (response.ok) {
                    document.getElementById('message').style.display = 'block';
                    document.getElementById('message').textContent = 'User ' + username + ' deleted, you will be logged out';
					window.location.href = 'http://localhost:5001/';
                } else {
                    throw new Error('User deletion failed. The user needs to be over 1 year old to be deleted.');
                }
            })
            .catch(error => {
                console.error('Error deleting user:', error.message);
				document.getElementById('message').style.display = 'block';
				document.getElementById('message').innerHTML = '<p style="color: red;">Error: ' + error.message + '</p>';
    	});
            
		}
		// For Trip
		function showTripsContainer() {
			getAllTrips();
    		document.getElementById('viewProfileContainer').style.display = 'none';
    		document.getElementById('viewTripContainer').style.display = 'block';
    		document.getElementById('message').style.display = 'none';
			document.getElementById('publishContainer').style.display = 'none';
		}
		function getAllTrips() {
			var xhr = new XMLHttpRequest();
			xhr.open('GET', 'http://localhost:5000/api/v1/trips', true);
		
			xhr.onload = function () {
				if (xhr.status >= 200 && xhr.status < 300) {
					console.log('Raw Response:', xhr.responseText);
		
					if (xhr.responseText.trim() !== "") {
						try {
							var data = JSON.parse(xhr.responseText);
							console.log('Trips:', data);
							updateTripList(data);
                            document.getElementById('viewTripContainer').style.display = 'block';
                            document.getElementById('viewProfileContainer').style.display = 'none';
                            document.getElementById('message').style.display = 'none';
						} catch (error) {
							console.error('Error parsing JSON:', error);
						}
					} else {
						console.error('Error: Empty response body');
					}
				} else {
					console.error('Error:', xhr.status, xhr.statusText);
				}
			};
			xhr.onerror = function () {
				console.error('Network error occurred');
			};
			xhr.send();
		}
		
		function updateTripList(trips) {
			const tripList = document.getElementById('tripList');
			tripList.innerHTML = '';
    		trips.forEach(trip => {
        		const tripDiv = document.createElement('div');
        		const listItem = document.createElement('p');
				
				const alternatePickUpLocation = trip['Alternate Pick-Up Location'] ? trip['Alternate Pick-Up Location']['String'] : '';
				const startTravelingTime = new Date(trip['Start Traveling Time']);
				const formattedStartTime = startTravelingTime.toLocaleTimeString();
        		listItem.innerHTML = "ID: " + trip.ID + "<br>" +
                             "Pick-Up Location: " + trip['Pick-Up Location'] + "<br>" +
                             "Alternate Pick-Up Location: " + alternatePickUpLocation + "<br>" +
                             "Start Traveling Time: " +  formattedStartTime + "<br>" +
                             "Destination Location: " + trip['Destination Location'] + "<br>" +
							 "Vacancies: " + trip['Number of Passengers Allowed'] + "<br>" +
                             "Published By: " + trip['Publisher'];
        		tripDiv.appendChild(listItem);
        		tripList.appendChild(tripDiv);
    		});
		}
		function showPublishContainer() {
			document.getElementById('publishContainer').style.display = 'block';
			document.getElementById('viewProfileContainer').style.display = 'none';
			document.getElementById('viewTripContainer').style.display = 'none';
		}

		function publishTrip() {
			const pickUpLocation = document.getElementById('pickUpLocation').value;
			const altPickUpLocation = document.getElementById('altPickUpLocation').value;
			const timeValue = document.getElementById('startTravelTime').value;

    		// Format the time in ISO 8601 format with a fixed date
    		const startTravelTime = new Date('2000-01-01T' + timeValue + ':00Z');
			const destinationLocation = document.getElementById('destinationLocation').value;
			const passengerNo = parseInt(document.getElementById('passengerNo').value);
			const requestBody = {
				"Pick-Up Location": pickUpLocation,
				"Alternate Pick-Up Location": altPickUpLocation,
				"Start Traveling Time": startTravelTime.toISOString(),
				"Destination Location": destinationLocation,
				"Number of Passengers Allowed": passengerNo
			};
			fetch('http://localhost:5000/api/v1/publishtrip/' + username, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(requestBody),
			})
			.then(response => {
				if (!response.ok) {
					throw new Error('HTTP error! Status: ' + response.status);
				}
				return response.text();
			})
			.then(data => {
				// Handle the successful response
			})
			.catch(error => {
				console.error('Fetch error:', error);
				// Handle the error
				alert('Error publishing trip: ' + error.message);
			});
		}
		
	</script>
</head>
<body>
	<div id="userButtonList">
	<button type="button" onclick="showProfileContainer()">View Profile</button>
	<button type="button" onclick="showTripsContainer()">View Trips</button>
	<button type="button" onclick="showPublishContainer()">Publish Trip</button>
	<a href="/logout"><button type="button">Log Out</button></a>
	</div>
	<div id="viewProfileContainer">
		<button type="button" onclick="showUserInfo()">Update Profile</button>
		<button type="button" onclick="showDeleteConfirmation()">Delete Profile</button>
		<div id="confirmDeleteContainer">
			<span>Are you sure you want to delete your user?</span>
			<div>
				<button type="button" onclick="deleteUser(username)">Delete</button>
			</div>
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
	</div>
	<div id="viewTripContainer">
		<ul id="tripList"></ul>
	</div>
	<div id="publishContainer">
			<form id="publishTripForm">
				<label for="pickUpLocation">Pick-Up Location:</label>
				<input type="text" id="pickUpLocation" required><br>
				<label for="altPickUpLocation">Alternate Pick-Up Location:</label>
				<input type="text" id="altPickUpLocation" required><br>
				<label for="startTravelTime">Start Traveling Time:</label>
				<input type="time" id="startTravelTime" required><br>
				<label for="destinationLocation">Destination Location:</label>
				<input type="text" id="destinationLocation" required><br>
				<label for="passengerNo">Maximum Number of Passengers:</label>
				<input type="text" id="passengerNo" required><br>
				<button type="button" onclick="publishTrip()">Publish</button>
			</form>
		</div>
    <div id="message"></div>
</body>
</html>
`

// ProfileHandler handles requests to the profile page
func CarOwnerHandler(w http.ResponseWriter, r *http.Request) {

	// Extract username from the URL
	username := r.URL.Query().Get("username")

	// Pass the username and additional user details to the template
	data := struct {
		Username string
	}{
		Username: username,
	}

	tmpl, err := template.New("home").Parse(carownerTemplate)
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
