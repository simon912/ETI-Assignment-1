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
	</script>
</head>
<body>
	<button type="button" onclick="redirectedToProfile()">View Profile</button>
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

	// Pass the username to the template
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
