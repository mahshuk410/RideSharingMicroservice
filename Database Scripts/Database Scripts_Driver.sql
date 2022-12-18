-- DATABASE SCRIPT FOR DRIVER DATABASE
 
 CREATE DATABASE driver_db;
 USE driver_db;
 
 
CREATE TABLE Drivers (
 driver_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY UNIQUE,
 driver_first_name VARCHAR(20) NOT NULL,
 driver_last_name VARCHAR(20) NOT NULL,
 driver_email VARCHAR(320) NOT NULL unique,
 driver_mobileNo varchar(8) NOT NULL unique,
 driver_idNo VARCHAR(9) NOT NULL unique,
 driver_licenseNo VARCHAR(15) NOT NULL UNIQUE
 );
 
 
-- Dummy data for Driver Table
 INSERT INTO Drivers (driver_first_name, driver_last_name, driver_email, driver_mobileNo, driver_idNo, driver_licenseNo) VALUES ("Tan","Kim","tk@gmail.com","94363636","T4364364E","GG 4363 F");
INSERT INTO Drivers (driver_first_name, driver_last_name, driver_email, driver_mobileNo, driver_idNo, driver_licenseNo) VALUES ("Ratan","Tata","rt@hotmail.com","90032453","T6782345E","FY 3435 F");

 USE driver_db;
 SELECT * FROM Drivers
 
