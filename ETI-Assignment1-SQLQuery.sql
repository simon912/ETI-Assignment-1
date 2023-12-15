CREATE database CarPoolingTrip;
USE CarPoolingTrip;

-- Creation of Users Table
CREATE Table Users (
	Username VARCHAR(20) NOT NULL PRIMARY KEY,
    Password VARCHAR(20) NOT NULL,
    UserGroup VARCHAR(20),
    FirstName VARCHAR(10),
    LastName VARCHAR(10),
    MobileNumber INT,
    EmailAddress VARCHAR(50),
    LicenseNo INT NULL,
    PlateNo VARCHAR(10) NULL,
    CreationDateTime Date
);
-- Insertion of Record for Users Table
INSERT INTO Users (Username, Password, UserGroup, FirstName, LastName, MobileNumber, EmailAddress, CreationDateTime) VALUES ("john123", "pass1234", "Passenger", "John", "Doe", 96232423, "johndoe@email.com", STR_TO_DATE("29-02-2020", "%d-%m-%Y"));
INSERT INTO Users (Username, Password, UserGroup, FirstName, LastName, MobileNumber, EmailAddress, LicenseNo, PlateNo, CreationDateTime) VALUES ("jane456", "pass1234", "Car Owner", "Jane", "Doe", 98534243, "jane123@email.com", 103436331, "SKW22G", STR_TO_DATE("24-05-2022", "%d-%m-%Y"));
INSERT INTO Users (Username, Password, UserGroup, FirstName, LastName, MobileNumber, EmailAddress, CreationDateTime) VALUES ("lee44", "amongus", "Passenger", "Bryan", "Lee", 95732952, "bryan@email.com", STR_TO_DATE("15-04-2023", "%d-%m-%Y"));
INSERT INTO Users (Username, Password, UserGroup, FirstName, LastName, MobileNumber, EmailAddress, LicenseNo, PlateNo, CreationDateTime) VALUES ("jojop2", "amogus", "Car Owner", "Joseph", "Joestar", 98643435, "joseph@email.com", 104953432, "SLT45G", STR_TO_DATE("12-12-2021", "%d-%m-%Y"));
INSERT INTO Users (Username, Password, UserGroup, FirstName, LastName, MobileNumber, EmailAddress, CreationDateTime) VALUES ("naruto55", "pass1234", "Passenger", "Naruto", "Uzumaki", 95732151, "hokage@email.com", STR_TO_DATE("24-08-2022", "%d-%m-%Y"));
INSERT INTO Users (Username, Password, UserGroup, FirstName, LastName, MobileNumber, EmailAddress, CreationDateTime) VALUES ("songoku", "dokkanbattle", "Passenger", "Son", "Goku", 91595959, "kamehameha@email.com", STR_TO_DATE("09-05-2022", "%d-%m-%Y"));
INSERT INTO Users (Username, Password, UserGroup, FirstName, LastName, MobileNumber, EmailAddress, LicenseNo, PlateNo, CreationDateTime) VALUES ("pepethefrog", "pass1234", "Car Owner", "Pepe", "Sadge", 91111111, "pepega@email.com", 101111111, "SHL22R", STR_TO_DATE("06-06-2023", "%d-%m-%Y"));
INSERT INTO Users (Username, Password, UserGroup, FirstName, LastName, MobileNumber, EmailAddress, CreationDateTime) VALUES ("petergriffin", "pass1234", "Passenger", "Peter", "Griffin", 91242424, "peter@email.com", STR_TO_DATE("03-10-2021", "%d-%m-%Y"));
INSERT INTO Users (Username, Password, UserGroup, FirstName, LastName, MobileNumber, EmailAddress, CreationDateTime) VALUES ("samuel2", "pass1234", "Passenger", "Samuel", "Tan", 91233322, "samuel2@email.com", STR_TO_DATE("10-05-2020", "%d-%m-%Y"));
INSERT INTO Users (Username, Password, UserGroup, FirstName, LastName, MobileNumber, EmailAddress, CreationDateTime) VALUES ("tomchan", "pass1234", "Passenger", "Tom", "Chan", 91511584, "tom@email.com", STR_TO_DATE("10-05-2022", "%d-%m-%Y"));

-- Creation of Trip Table 
CREATE Table Trips (
	TripID INT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
    PickUpLocation VARCHAR(50), 
    AltPickUpLocation VARCHAR(50) NULL, 
    StartTravelTime Time, 
    DestinationLocation VARCHAR(50), 
    PassengerNoLeft INT,
    MaxPassengerNo INT,
    Status VARCHAR(10) CHECK (Status IN('Pending', 'Active')), 
    Publisher VARCHAR(20), 
    FOREIGN KEY (Publisher) REFERENCES Users(Username)
);

-- Insertion of Record for Trips Table
INSERT INTO Trips (TripID, PickUpLocation, StartTravelTime, DestinationLocation, PassengerNoLeft, MaxPassengerNo, Status, Publisher) VALUES (1, "Ang Mo Kio Road", STR_TO_DATE('12:35', '%H:%i'), "Ang Mio Kio Hub", 2, 3, "Pending", "jane456");
INSERT INTO Trips (TripID, PickUpLocation, AltPickUpLocation, StartTravelTime, DestinationLocation, PassengerNoLeft, MaxPassengerNo, Status, Publisher) VALUES (2, "Bukit Panjang Ring Road", "Jelapang Road", STR_TO_DATE('14:20', '%H:%i'), "Lot One", 3, 3, "Pending", "jojop2");
INSERT INTO Trips (TripID, PickUpLocation, AltPickUpLocation, StartTravelTime, DestinationLocation, PassengerNoLeft, MaxPassengerNo, Status, Publisher) VALUES (3, "Choa Chu Kang Road", "Teck Whye Road", STR_TO_DATE('18:20', '%H:%i'), "Orchard Road", 0, 3, "Pending", "jane456");
INSERT INTO Trips (TripID, PickUpLocation, StartTravelTime, DestinationLocation, PassengerNoLeft, MaxPassengerNo, Status, Publisher) VALUES (4, "Jurong West Road", STR_TO_DATE('20:20', '%H:%i'), "Bugis Road", 1, 2, "Pending", "pepethefrog");
INSERT INTO Trips (TripID, PickUpLocation, StartTravelTime, DestinationLocation, PassengerNoLeft, MaxPassengerNo, Status, Publisher) VALUES (5, "Kembengan Road", STR_TO_DATE('8:40', '%H:%i'), "Bedok Mall", 0, 2, "Active", "jane456");
INSERT INTO Trips (TripID, PickUpLocation, StartTravelTime, DestinationLocation, PassengerNoLeft, MaxPassengerNo, Status, Publisher) VALUES (6, "Orchard Road", STR_TO_DATE('10:35', '%H:%i'), "Bugis Road", 2, 3, "Active", "pepethefrog");

-- Creation of Enrollment Table
CREATE Table Enrollment (
	EnrollID INT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
    Username VARCHAR(20),
    TripID INT, 
    FOREIGN KEY (Username) REFERENCES Users(Username),
    FOREIGN KEY (TripID) REFERENCES Trips(TripID)
);

-- Insertion of Record for Enrollment Table
INSERT INTO Enrollment (EnrollID, Username, TripID) VALUES (1, "john123", 1);
INSERT INTO Enrollment (EnrollID, Username, TripID) VALUES (2, "john123", 3);
INSERT INTO Enrollment (EnrollID, Username, TripID) VALUES (3, "naruto55", 3);
INSERT INTO Enrollment (EnrollID, Username, TripID) VALUES (4, "lee44", 3);
INSERT INTO Enrollment (EnrollID, Username, TripID) VALUES (5, "songoku", 4);
INSERT INTO Enrollment (EnrollID, Username, TripID) VALUES (6, "john123", 5);
INSERT INTO Enrollment (EnrollID, Username, TripID) VALUES (7, "petergriffin", 5);
INSERT INTO Enrollment (EnrollID, Username, TripID) VALUES (8, "petergriffin", 6);