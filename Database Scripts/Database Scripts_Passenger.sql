-- DATABASE SCRIPT FOR PASSENGER DATABASE
 
 CREATE DATABASE passenger_db;
 USE passenger_db;
 
 
CREATE TABLE Passengers (
 passenger_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY UNIQUE,
 passenger_first_name VARCHAR(20) NOT NULL,
 passenger_last_name VARCHAR(20) NOT NULL,
 passenger_email VARCHAR(320) NOT NULL unique,
 passenger_mobileNo varchar(8) NOT NULL unique);
 
 
-- Dummy data for Passenger Table
 INSERT INTO Passengers (passenger_first_name,passenger_last_name,passenger_email,passenger_mobileNo) VALUES ("John","Doe","jd@gmail.com","98765421");
INSERT INTO Passengers (passenger_first_name,passenger_last_name,passenger_email,passenger_mobileNo) VALUES ("Dim","Sum","ds@hotmail.com","91234567");
 INSERT INTO Passengers (passenger_first_name,passenger_last_name,passenger_email,passenger_mobileNo) VALUES ("Le","Bezos","lbezo@gmail.com","94324563");
 
USE passenger_db;
SELECT * FROM Passengers
 