# ETI Assignment 1

### Instruction for setting up and running the microservices
1. Run the MySQL Script for the Creation of Database and Table as well as Insertion of Records
2. Run the two microservices (user.go & trip.go)
3. Open index.html

### Design Considerations/Microservices
The backend is split into two Microservices that do separate functions. Each Microservice has their own role and responsbility to fulfill while being linked to each other.
As for the frontend, it is done through the use of HTML CSS Javascript to display the relevant information as well as make use of several HTTP method such as GET, POST, PUT & DELETE.
The frontend and backend are integrated through the use of REST APIs.

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
![architecture_diagram-v2](https://github.com/simon912/ETI-Assignment-1/assets/93958709/4e89fd8e-c6f1-475a-afe0-469da707e88a)

