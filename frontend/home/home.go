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
</style>
 <script>
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
		function redirectedToProfile() {
			const urlParams = new URLSearchParams(window.location.search);
        	const username = urlParams.get('username');
        
        // Make a GET request to the server to handle the redirection
        fetch('/redirect-profile?username=' + encodeURIComponent(username))
            .then(response => {
                if (response.ok) {
                    // Redirect to the profile page
                    window.location.href = '/profile?username=' + encodeURIComponent(username);
                } else {
                    console.error('Error handling redirection:', response.status, response.statusText);
                }
            })
            .catch(error => {
                console.error('Error handling redirection:', error);
            });
		}
		function logOutUser() {

		}
		function retrieveUserData() {
			const urlParams = new URLSearchParams(window.location.search);
			const username = urlParams.get('username');
		
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
	</script>
</head>
<body>
	<div id="userButtonList">
	<button type="button" onclick="redirectedToProfile()">View Profile</button>
	<button type="button" onclick="redirectedToProfile()">View Trips</button>
	<button type="button" onclick="logOutUser()">Log Out</button>
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
	http.HandleFunc("/profile", profile.ProfileHandler)
	http.HandleFunc("/display-profile", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		redirectToProfileSuccess(w, r, username)
	})

}
func redirectToProfileSuccess(w http.ResponseWriter, r *http.Request, username string) {
	http.Redirect(w, r, "http://localhost:5001/profile?username="+username, http.StatusSeeOther)
}
