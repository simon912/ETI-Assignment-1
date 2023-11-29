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
		
	</script>
</head>
<body>
	
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
