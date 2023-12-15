# ETI Assignment 1

### Instruction for setting up and running the microservices
1. Run the MySQL Script for the Creation of Database and Table as well as Insertion of Records
2. Run the two microservices (user.go & trip.go)
3. Open index.html

### Design Considerations
The backend is split into two Microservices that do separate functions. Each Microservice has their own role and responsbility to fulfill while being linked to each other.
As for the frontend, it is done through the use of HTML CSS Javascript to display the relevant information as well as make use of several HTTP method such as GET, POST, PUT & DELETE.
The frontend and backend are integrated through the use of REST APIs.

### User Microservice 🧍
The **User Microservice** is designed to handle operations that mainly involves the User (Passenger & Car Owner).
The following operations as well as the condition for them to work are:
* **Login of User**
  * If the user provides the correct username and password 
* **Registration of User**
  * If the username of their choice is not taken
* **Updating User Information** (Only Mobile Number & Email Address can be updated)
* **Deletion of User**
  * If the user is over 1 year old or 365 days old
  * If the user has not enrolled into a trip
  * If the user has not published any trip
* **Changing the Passenger's User Group to Car Owner** if they provide the License Number and Plate Number of their car
  * If the user has not enrolled into a trip 
  
### Trip Microservice 🚗
The **Trip Microservice** is designed to handle operations that involves managing the Car Pooling trips published by Car Owners.
The following operations as well as the condition for them to work are:
* **Viewing Trips** that has been published by Car Owners 
* **Enrollment of Trips**
  * Only Passengers can enroll into trips
  * If there is vacancy in the trip
  * If the trip status is Pending
  * If the user has not enrolled into that trip
* **Publishing Trips** with their own information such as Start Traveling Time, Pick-Up Location and Destination Location
  * Only Car Owner can publish trips
* **Starting The Trip**
  * If Passenger has enrolled into the trip 
* **Canceling the Trip** 
  * If the Car Owner is not within 30 minutes cancellation window[^1]
[^1]: For example, if the Trip's Start Traveling Time is 8:30 pm and the current time is 8:05 pm, the Car Owner will not be able to cancel the trip
### Other external tools or library
* **moment.js**[^2] for displaying the Date and Time on frontend
* **Bootstrap**[^3] to enhance the frontend design
[^2]:https://momentjs.com/
[^3]:https://getbootstrap.com/
### Other images used
* **[index-background.jpg](frontend/images/index-background.jpg)**[^4]
[^4]:https://www.pexels.com/photo/man-inside-vehicle-13861/
### Architecture Diagram
![architecture_diagram-v2](https://github.com/simon912/ETI-Assignment-1/assets/93958709/4e89fd8e-c6f1-475a-afe0-469da707e88a)

