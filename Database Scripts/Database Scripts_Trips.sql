-- DATABASE SCRIPT FOR TRIPS DATABASE

CREATE DATABASE trip_db;
USE trip_db;

CREATE TABLE Trips (
 trip_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY UNIQUE,
 fromPostalCode BIGINT NOT NULL ,
 toPostalCode BIGINT NOT NULL,
 tripStatus VARCHAR(8) NOT NULL,
startTime DATETIME NULL,
 endTime DATETIME NULL,
passenger_id INT NOT NULL,
driver_id INT NOT NULL,
 CHECK (fromPostalCode >=100000),
 CHECK (toPostalCode >= 100000 AND fromPostalCode <> toPostalCode),
 CHECK(tripStatus IN ("Pending","Started","Ended"))
 );

-- Dummy data for Trips Table
 INSERT INTO Trips(fromPostalCode, toPostalCode, tripStatus, startTime, endTime, passenger_id, driver_id)
VALUES(453253,984242,"Started",NOW(),'2022-12-3 13:15:00',12,6)
 INSERT INTO Trips(fromPostalCode, toPostalCode, tripStatus, startTime, endTime, passenger_id, driver_id)
VALUES(553253,334882,"Started",NOW(),'2022-12-3 19:15:00',5,4)
INSERT INTO Trips(fromPostalCode,toPostalCode,tripStatus,startTime,endTime,passenger_id,driver_id)
VALUES(346433,666644,'Ended','2022-12-17 15:46:33','2022-12-17 15:47:14',3,4)

 INSERT INTO Trips(fromPostalCode, toPostalCode, tripStatus, startTime, endTime, passenger_id, driver_id)
VALUES(452987,522619,'Ended','2022-12-17 19:31:46','2022-12-17 19:31:52',3,1)
 INSERT INTO Trips(fromPostalCode, toPostalCode, tripStatus, startTime, endTime, passenger_id, driver_id)
VALUES(436534,224455,'Ended','2022-12-17 19:32:04','2022-12-17 19:33:00',3,1)
 INSERT INTO Trips(fromPostalCode, toPostalCode, tripStatus, startTime, endTime, passenger_id, driver_id)
VALUES(452987,522619,'Ended','2022-12-18 10:59:39','2022-12-18 10:59:54',3,6)

USE trip_db;
SELECT * FROM Trips WHERE passenger_id = 3




