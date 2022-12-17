# RideSharingMicroservice

## Description of the Ride-Sharing App
The ride-sharing application adopts a microservice architecture where each REST API represents a distinct microservice. They are broken into Driver API, Passenger API and Trips Management API.
The driver API allows drivers to manage their account details(create and update), view assigned trips, start and end these trips. 
Whereas the Passenger API allows users to register their Passenger Accounts, update their account details, request for trips from one destination to another as well as view their trip history.
##

## Microservice Architecture Diagram
![Architecture-diagram](https://user-images.githubusercontent.com/73008987/208143360-123538d0-443e-42ca-af4a-983a1583f5af.png)

A Component-based front-end is built for drivers and passengers. All communications to microservices will be made from the Driver and Passenger UI over RESTful APIs. The console UIs will be running on PowerShell/ Command Prompt Terminals. Each of the microservices were developed using the Golang. For each microservice, a REST API has been created. All the databases are developed using MySQL engine.

## Documentation
### Domain
![domain](https://user-images.githubusercontent.com/73008987/208224379-008b2021-ce8f-407c-a40b-7c19df7934b8.png)

The domain of the application developed is a Ride-Sharing Platform. The Ride Sharing Platform domain can be separated into 3, loosely coupled sub-domains known as the Driver Management, Passenger Management and Trip Management. The main users for this application are Drivers and Passengers.

### Glossary
Definitions of terminologies widely used in the sub-domains of this application.

- Passenger – An application user requesting for a Ride-Sharing trip to travel from one destination to another.
- Driver – An application user who owns a vehicle, using it to transport passengers from one destination to another. 
- Trip – A journey requested by the passenger from one source to a destination location allocated to a driver. 


### Sub-Domains

Having identified the 3 domains, each of them are supported by a back-end microservice each that supports the functionalities of the Ride Sharing Application.

**1. Driver Management Sub-Domain**
-	Supports the Driver Account Creation process. Each driver registering on the platform has to enter his first name, last name, mobile number, email address, identification number and car license number. 
-	All the above-mentioned credentials (except Identification Number) can be updated.
-	Driver accounts cannot be deleted.

**2. Passenger Management Sub-Domain**
-	Supports the Passenger account creation process. Every passenger registering has to enter his, first name, last name, mobile number, and email address.
-	Any of the above-mentioned passenger credentials can be updated.
-	Similarly, no passenger account can be deleted.

**3. Trip Management Sub-Domain**
-	Applicable for both Drivers and Passengers.
-	Passenger can request for a trip by inputting postal codes of pick-up and drop-off location.
-	An available driver, not driving any passenger, will be assigned one trip requested at any one time. 
-	Driver has ability to initiate start trip and end trip.
-	Passenger can retrieve his trip history in reverse chronological order (latest order first).

### Context Map

![contextMap](https://user-images.githubusercontent.com/73008987/208224523-87241674-9596-4afa-b73a-11700b5dd005.png)

From the 3 subdomains identified, they have been broken down into 3 bounded contexts. The driver Account Management context is responsible for passing the details to the Trips Allocation bounded context for allocating a particular Trip ID to the Car Driver. Whereas for the Passenger Account Management context, it is responsible for passing the Passenger Name and Id to the Trip Allocation context to authorise the passenger as well as to associate a new trip to the Passenger Id. As the passenger needs to retrieve his trip history, the “Trip History Retrieval” feature highly depends on the Trips Allocation feature to log every completed trip.  Due to the tightly coupled dependency, these 2 features are identified to be in the Trip Management bounded context. The arrowhead lines define that communications between the bounded contexts will be sent over APIs.

### Bounded Contexts
**Domain Model for Passenger Account Management** 

![passenger_domain_model](https://user-images.githubusercontent.com/73008987/208224731-a02dd57c-6a37-436c-9e01-112292f96a70.png)

**Domain Model for Driver Account Management** 

![driver_domain_model](https://user-images.githubusercontent.com/73008987/208224755-a057599e-e6c3-4424-a039-0268a2a3dbfb.png)

**Domain Model for Trip Management** 

![trip_domain_model](https://user-images.githubusercontent.com/73008987/208224771-190bea8a-6637-4289-afe8-2ad7aad7b4fb.png)

### Entities and Aggregates

**Passenger Management Sub-Domain**

![passenger_class](https://user-images.githubusercontent.com/73008987/208224823-5c71e61a-67ae-459b-a6d6-f12094a77009.png)

**Driver Management Sub-Domain**

![driver_class](https://user-images.githubusercontent.com/73008987/208224841-a107feb1-0bb1-4dd4-bac8-a50c9e4950b1.png)

**Trips Management Sub-Domain**

![complete_class_diagram](https://user-images.githubusercontent.com/73008987/208224851-5b1f4c0c-d3f6-48fb-9d28-6d3a6209e0c7.png)

One trip is requested by a passenger and driven by a driver. 1 driver and 1 passenger can drive and request multiple trips respectively. Only the PassengerId and DriverId are sent over APIs to associate the trip records with the Driver and Passenger respectively.


## Setting Up Microservices
1.	Download and Execute the MySQL Database Script 
2.	Download and Execute the driver-service.go, passenger-service.go and trip-service.go.
NOTE: Run all of these microservices simultaneously on separate VS code/PowerShell/command prompt terminals
NOTE: Open all the above-mentioned go files from the folder level to import all module dependencies. 
3.	Allow the firewall connections through ports when prompted by Windows Firewall Security:  

| PORT | MICROSERVICE           |
|------|------------------------|
| 5000 | Passenger Microservice |
| 5002 | Driver Microservice    | 
| 5003 | Trip Microservice      |


4.	Download and execute the Driver.go and Passenger.go console files.
NOTE: Open all the above-mentioned go files from the folder level to import all module dependencies. 
5.	Tools Required: 
 -  Postman – Test the APIs
 -  Visual Studio Code – Open the Microservice and console programs
 -  Golang extension on Visual Studio Code– to import dependencies
 -  MySQL Workbench – to run the MySQL database scripts


## Rest API Documentation

### Status Codes
| **HTTP Status**              | **Status Code** | **Response Message**                                          |
|---------------------------------|-----------------|---------------------------------------------------------------|
| http.StatusMethodNotAllowed     | 405             | For invalid API methods e.g. DELETE Request on a PUT Endpoint |
| http.StatusConflict             | 409             | duplicate record found in database                            |
| http.StatusInternalServerError  | 500             | Database query execution error or JSON payload error          |
| http.StatusAccepted             | 201             | successful inserted record in database                        |
| http.StatusOK                   | 200             | successfully retrieve OR updated record from database         |
| http.StatusNotFound             | 404             | record not found in database                                  |
| http.StatusBadRequest           | 400             | invalid request payload e.g. empty inputs provided            |


### Driver API
API Endpoint: http://localhost:5002
Sample Driver Object in JSON 
![Driver](https://user-images.githubusercontent.com/73008987/208151298-b51aea93-c94e-432e-a8cf-987d578e20cc.png)

| **Propetry**   | **Description**                                                                                         |
|-----------------|--------------------------------------------------------------------------------------------------------|
| driver_id       |  unique identifier of a driver *auto-generated by database*                                           |
| First Name      |  first name of a driver  (maximum length of 20 characters)                                             |
| Last Name       |  last name of a driver (maximum length of 20 characters)                                               |
| Mobile Number   |  Singapore-registered 8 digit mobile number *must be unique*                                           |
| Email Address   |  valid email address *must be unique* (maximum length of 320 characters)                               |
| Car License No  |  License plate No. of Singapore-registered Vehicle *must be unique and maximum length of 15 characters |


**1.  Signup As New Driver**

API Request Body

```bash 
   curl  -X POST http://localhost:5002/api/v1/driver --data '{
	"First Name": "Jeff",
	"Last Name" : "Son"
	"Mobile Number" : "91234567",
	"Email Address" : "demo@gmail.com",
	"Identification No": "T0234567A",
	"Car License No": "FA2345D"
}'-H "Content-Type: application/json"
```

API Response

A successful account creation returns a 202 StatusAccepted status.
Any duplicate email Address or mobile number entered by existing drivers will return a 409 Status Conflict response.
```code 
Passenger has been successfully signed up
```

**2.  Update Driver Details** 

API Request

Option 1: Update with existing Mobile Number

```bash
curl – X PUT http://localhost:5002/api/v1/driver?mobileNo={existingMobileNo}
--data '{
	"First Name": "*",
	"Last Name": "*"
	"Mobile Number": "91234567",
	"Email Address": "updateDemo@gmail.com",
	"Car License No": "FA2345D"
}' -H "Content-Type:application/json"
```
Option 2: Update with existing Email Address

```bash
curl – X PUT http://localhost:5002/api/v1/driver?emailAddress={existingEmailAddress}
--data '{
	"First Name": "*",
	"Last Name" : "*"
	"Mobile Number" : "91234567",
	"Email Address" : "updateDemo@gmail.com",
	"Car License No": "FA2345D"
}'
-H "Content-Type:application/json"
```

**existingEmailAddress** OR **existingMobileNo** are required parameters

Note: Driver account’s Identification Number cannot be updated

> The * symbol denotes the details to remain unchanged by user.

API Response

A successful account update will return a 200 OK status.

```code
Driver record successfully updated
```

**3.  Get Available Drivers** 

Return list of driver ids available to accept trip requests

API Request

```bash
curl -X GET "http://localhost:5002/api/v1/driver/all"
-H "Content-Type: application/json"
```
API Response

Returns a 200 OK status with driver IDs

```code
{
  "Driver IDs": [
    1,
    3,
    4,
    5,
    6
  ]
}
```
### Passenger API

API Endpoint - http://localhost:5000

| **Property**   |  **Description**                                                           |
|----------------|----------------------------------------------------------------------------|
| Passenger Id   |  unique identifier of a passenger *auto-generated by database*             |
| First Name     |   first name of a Passenger *maximum length of 20 characters*              |
| Last Name      |  last name of a driver *maximum length of 20 characters*                   |
| Mobile Number  |  Singapore-registered 8-digit mobile number *must be unique*               |
| Email Address  |  valid email address *must be unique and maximum length of 320 characters* |

**1.  Create a New Passenger Account**

Passenger has to enter details to register account.

API Request

```bash
curl -X POST "http://localhost:5000/api/v1/passenger" --data '{
    "First Name": "John",
    "Last Name": "Doe",
    "Mobile Number": "80123456",
    "Email Address": "demoPassenger@gmail.com",
}'
-H 'Content-Type: application/json'
```

API Response
A successful passenger account returns a 202 StatusAccepted code
Any duplicate email Address or Mobile Number used will throw a 409 StatusConflict error.
```code
Passenger has been successfully signed up
```

**2.  Update Passenger Details**
API Request

Option 1: Update with existing mobileNo
```bash
curl -X PUT “http://localhost:5000/api/v1/passenger?mobileNo={existingMobileNo}” 
--data ‘{
"First Name": "*",
    "Last Name": "*",
    "Mobile Number": "80125354",
    "Email Address": "updateDemoPassenger@gmail.com",
}’
– H ‘Content-Type: application/json’
```
Option 2: Update with existing emailAddress
```bash
curl -X PUT “http://localhost:5000/api/v1/passenger?emailAddress={existingEmailAddress}” 
--data '{
"First Name": "*",
    "Last Name": "*",
    "Mobile Number": "80125354",
    "Email Address": "updateDemoPassenger@gmail.com",
}'
– H 'Content-Type: application/json'
```
> The * symbol denotes the details to be remained unchanged by user.

API Response

A successful passenger account update returns a 200 OK status.
```bash
Passenger record successfully updated
```
Any existing mobile number or email address provided by a query string not found will throw a 404 Not Found error.
```bash
Passenger with the mobile number not found
```
```bash
Passenger with the email address not found
```

### Trip Management API

API Endpoint: http://localhost:5003

| **Property**          | **Description**                                                                                                      |
|-----------------------|----------------------------------------------------------------------------------------------------------------------|
| Trip Id (int)         | Unique identifier of Trip (auto-generated by database)                                                                             |
| From (int64)          | 6-digit Singapore Postal Code                                                                                        |
| To (int64)            | 6-digit Singapore Postal Code                                                                                        |
| Trip Status (string)  |  Pending (Driver yet to accept), Started (Driver has started Trip), Ended (Driver Ended Trip upon passenger Dropoff) |
| Start Time(string)    |  Date and Time of Trip Started represented in YYYY-MM-DD hh:mm:ss                                                    |
| End Time(string)      |  Date and Time of Trip Started represented in YYYY-MM-DD hh:mm:ss                                                    |
| Passenger Id(int)     |  Unique identifier of passenger requesting trip                                                                      |
| Driver Id(int)        |  Unique Identifier of Driver transporting the passenger                                                              |


**1.  Passenger Requests New Trip**

API Request

Sends a new Trip object, communicates to the Driver API to assign an available driver and notifies passenger of the trip status. passengerId url parameter auto retrieved from Passenger API.

```bash
curl -X POST "http://localhost:5003/api/v1/trips/passenger/{passengerId}"
--data '{
    "From":452987,
    "To":522619,
    "Trip Status":”Pending”,
    "Passenger Id":5
}'
-H'Content-Type:application/json'
```

API Response

```code
"Created Trip":{
  //trip object
}
```

**2.  Passenger View Trip History**

API Request

Passenger list of trips taken will be returned in JSON format in reverse chronological order. Latest Ride First, earliest ride last. 

```bash
curl -X GET "http://localhost:5003/api/v1/trips/passenger/{passengerId}" -H 'Content-Type:application/json'
```

API Response

Returns a 200 OK. Error 404 denotes passenger has not made any rides.

```code
Trips: {
	//
}
```

**3. Driver View Assigned Trip**

API Request

Driver sees the trip details assigned to him. Valid driverId is required.

API Response

A trip assigned to that driver returns a 200 OK. If no trips have been assigned to that driver,a 404 Not Found response is shown. 
```bash
curl -X GET "http://localhost:5003/api/v1/trips/driver/{driverId}" -H 'Content-Type:application/json'
```
Example Response
```code
{
	//
}
```
**4.  Driver Updates the Trip Status**

API Request

Driver can either Start or End the Trip. Valid DriverId and corresponding tripId is required.

```bash
curl -X PUT "http://localhost:5003/api/v1/trips/driver/{driverId}/{tripId}" --data '**{Status:Started}**'-H'Content-Type:application/json'
```
```bash
curl -X PUT "http://localhost:5003/api/v1/trips/driver/{driverId}/{tripId}" --data '**{Status:Ended}**'-H'Content-Type:application/json'
```

API Response

```code
Trip has been successfully Ended
```

```code
Trip has been successfully Ended
```

**5.  Retrieve the Trip Status**

API Request

Passenger will be shown his trip status whenever driver has initiated a start or end trip accordingly. A valid Trip Id has to be provided in the URL parameter.

```bash
curl -X GET "http://localhost:5003/api/v1/trips/{tripId}"
 -H 'Content-Type:application/json'
```
API Response

Valid trip found with a status returns a 200 OK status while any invalid Trip Id returns a 404 Not Found status.


```code
{“Status”: “Started”}
```

```code
{“Status”: “Ended”}
```
