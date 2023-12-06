# ETI Assignment 1

### How to run the program
1. Run the two microservices (user.go & trip.go)
2. Open index.html

### Design Considerations/Microservices
The backend is split into three Microservices that do separate functions. Each Microservice has their own role and responsbility to fulfill while being linked to each other.
#### User Microservice
The User Microservice is designed to handle operations that mainly involves the User (Passenger & Car Owner) such as:
- Login of user
- Registration of user
- Updating user information
- Changing the Passenger's User Group to Car Owner if they provide the License Number and Plate Number of their car
#### Trip Microservice
The Trip Microservice is designed to handle operations that involves managing the Car Pooling trips published by Car Owners such as:
- Viewing trips that has been published by Car Owners
- Enrolling into trips if there is vacancy in the trip
- Viewing trips that has been published by Car Owners
- Publishing trips with their own information such as Start Traveling Time, Pick-Up Location and Destination Location

### Architecture Diagram
![architecture_diagram-v2](https://github.com/simon912/ETI-Assignment-1/assets/93958709/70622e04-3916-4edb-b1b4-c9b35ee9b212)
