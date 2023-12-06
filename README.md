# ETI Assignment 1

### How to run the program
1. Run main.go
2. Launch index.html 

### Design Considerations/Microservices
The backend is split into three Microservices that do separate functions. Each Microservice has their own role and responsbility to fulfill while being linked to each other.
#### Login Microservice
The Login Microservice is designed to handle operations such as:
- Login of user
- Registration of user
#### Passenger Microservice
The Passenger Microservice is designed to handle operations such as:
- Updating user information
- Changing the Passenger's User Group to Car Owner if they provide the License Number and Plate Number of their car
- Deletion of user if their account is over a year old
- Viewing trips that has been published by Car Owners
- Enrolling into trips if there is vacancy in the trip
#### Car Owner Microservice
While a little similar to Passenger, The Car Owner Microservice is designed to handle operations such as:
- Updating user information
- Deletion of user if their account is over a year old
- Viewing trips that has been published by Car Owners
- Publishing trips with their own information such as Start Traveling Time, Pick-Up Location and Destination Location



### Architecture Diagram
![architecture_diagram-v1](https://github.com/simon912/ETI-Assignment-1/assets/93958709/00d8026e-a4f0-4dce-ad44-0e1aab2c1365)
