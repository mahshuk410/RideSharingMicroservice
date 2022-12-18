package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type Driver struct {
	driver_id    int
	FirstName    string `json:"First Name"`
	LastName     string `json:"Last Name"`
	MobileNumber string `json:"Mobile Number"`
	EmailAddress string `json:"Email Address"`
	IdNo         string `json:"Identification No"`
	CarLicenseNo string `json:"Car License No"`
}

var driverList []Driver

// API Status Codes
// http.StatusMethodNotAllowed -- 405 -- For invalid API methods e.g. DELETE Request on a PUT Endpoint
// http.StatusConflict -- 409 -- duplicate record found in database
// http.StatusInternalServerError -- 500 -- Database query execution error or JSON payload error
// http.StatusAccepted - 201 -- successful inserted record in database
// http.StatusOK -- 200 - successfully updated record in
// http.StatusNotFound -- 404 -- driver record not found in database
// http.StatusBadRequest -- 400 -- invalid request payload e.g. empty inputs provided

func main() {

	router := mux.Router{}
	router.HandleFunc("/api/v1/driver", driverSignup).Methods("POST")
	router.HandleFunc("/api/v1/driver", updateDriver).Methods("PUT")
	router.HandleFunc("/api/v1/driver/all", requestAvailableDrivers).Methods("GET")
	router.HandleFunc("/api/v1/driver/{driverId}", getSpecificDriver).Methods("GET")
	fmt.Println("Listening at port 5002...")
	log.Fatal(http.ListenAndServe(":5002", &router))
}

func driverSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		searchEmailAddress, lookupMobileNo string //will be used to perform search operations
		invalidDriver                      bool
		newDriverData                      Driver
	)

	//initialise database variable, will be passed as parameters in methods
	driver_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/driver_db")
	if err != nil { //check for database opening errors
		driver_db.Close()
		log.Fatal(err.Error())
	}
	defer driver_db.Close()

	//synchronise driverList with latest copy of Driver records
	getAllDriversRecords(*driver_db)

	//handle signup function
	if r.Method == "POST" {
		if body, err := io.ReadAll(r.Body); err == nil {
			if err := json.Unmarshal(body, &newDriverData); err == nil {

				searchEmailAddress, lookupMobileNo = newDriverData.EmailAddress, newDriverData.MobileNumber //will be used to perform search operations

				for _, obj := range driverList {
					// search condition 1: Check for existing Matching email address
					if searchEmailAddress == obj.EmailAddress {
						invalidDriver = true
						w.WriteHeader(http.StatusConflict) // duplicate driver profile found
						fmt.Fprintf(w, "%s", bytes.NewBufferString("Email Address "+searchEmailAddress+" used by another Driver"))

						break
					}
					// search condition 2: Check for existing Matching Mobile Number
					if lookupMobileNo == obj.MobileNumber {
						invalidDriver = true
						w.WriteHeader(http.StatusConflict) // duplicate driver profile found
						fmt.Fprintf(w, "%s", bytes.NewBufferString("Mobile Number "+lookupMobileNo+" used by another Driver"))
						break
					}
				}
				if !invalidDriver { // normal flow --> all unique values

					//insert driver record to database
					status, newDriverId := insertDriverRecord(newDriverData, *driver_db)
					getAllDriversRecords(*driver_db) //refresh driver List after insert operation
					// status outcome of insert operation is to be notified to driver
					switch status {
					case 0:
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, "%s", bytes.NewBufferString("Database error signing up driver "+newDriverData.FirstName+" "+newDriverData.LastName))

					case 1:
						data := map[string]int64{"New Driver Id": newDriverId}
						jsonBody, _ := json.Marshal(data)
						w.WriteHeader(http.StatusAccepted)
						fmt.Fprintf(w, "%s", bytes.NewBuffer(jsonBody))

					}
				}

			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Error signing up driver : %v", err)
			fmt.Fprintf(w, "%s", bytes.NewBufferString("Error signing up driver"))
		}

	} else { // restrict API operations by other HTTP protocols

		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Invalid API Operation, only POST requests Allowed")
	}
}

func updateDriver(w http.ResponseWriter, r *http.Request) {

	//initialise database variable, will be passed as parameters in methods
	driver_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/driver_db")
	if err != nil { //check for database opening errors
		driver_db.Close()
		log.Fatal(err.Error())
	}
	defer driver_db.Close()

	w.Header().Set("Content-Type", "application/json")

	//refresh driverList records to get latest data
	getAllDriversRecords(*driver_db)
	var (
		driverFound    bool
		query          = r.URL.Query()
		newMap, oldMap map[string]string
		//get query string arguments
		lookupMobileNo     string = query.Get("mobileNo")     // search for existing Driver record by mobile number
		lookupEmailAddress string = query.Get("emailAddress") // search for existing Driver record by email address

	)

	if r.Method == "PUT" {
		if body, err := io.ReadAll(r.Body); err == nil {

			json.Unmarshal(body, &newMap)

			for _, driver := range driverList {

				if (lookupMobileNo == driver.MobileNumber) || (lookupEmailAddress == driver.EmailAddress) { // mobile number found OR email address found
					jsonString, _ := json.Marshal(driver)
					json.Unmarshal(jsonString, &oldMap) //retrieve old driver details

					driverFound = true
					for k, _ := range newMap {
						if strings.Contains(newMap[k], "*") {
							newMap[k] = oldMap[k] //for unchanged details denoted as * by driver, existing value will be copied over
						}
					}
					//driver_id is updated for the new Driver map since driver_id of driver is initialised to 0
					//id is required for "WHERE" clause parameter during update operation
					newMap["Driver Id"] = strconv.Itoa(driver.driver_id)

					status := updateDriverRecord(newMap, *driver_db) //update driver record in database

					switch status { //format http response body based on search and update status
					case 0:
						w.WriteHeader(http.StatusInternalServerError)

						fmt.Fprintf(w, "%s", bytes.NewBufferString("Database error updating driver "+driver.FirstName+" "+driver.LastName+"'s record")) //display the old First Name and Last Name to inform Client whose record could'nt be updated

					case 1:
						w.WriteHeader(http.StatusOK)
						fmt.Fprintf(w, "%s", bytes.NewBufferString("Driver record successfully updated"))

					}
					break
				}

			}
			//format http response body based on invalid result type
			if !driverFound && lookupMobileNo != "" { //invalid mobile number provided
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "%s", bytes.NewBufferString("Driver with the mobile number "+lookupMobileNo+" not found"))
			} else if !driverFound && lookupEmailAddress != "" {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "%s", bytes.NewBufferString("Driver with the email address "+lookupEmailAddress+" not found"))
			}

		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed) //error 405 for invalid operation
		fmt.Fprintf(w, "%s", "Invalid API Operation")
	}
}
func getSpecificDriver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//initialise database variable, will be passed as parameters in methods
	driver_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/driver_db")
	if err != nil { //check for database opening errors
		driver_db.Close()
		log.Fatal(err.Error())
	}
	defer driver_db.Close()
	if r.Method == "GET" {
		driver_id, _ := strconv.Atoi(mux.Vars(r)["driverId"])
		if driver := getDriverRecord(*driver_db, driver_id); driver.FirstName != "" { //check if valid driver has been found
			w.WriteHeader(http.StatusOK)

			resBody, _ := json.Marshal(driver)
			fmt.Fprintf(w, "%v", bytes.NewBuffer(resBody))
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%v", bytes.NewBufferString("No Driver Found"))
		}
	} else { // restrict API operations by other HTTP protocols

		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Invalid API Operation, only POST requests Allowed")
	}
}
func requestAvailableDrivers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		driverIdList []int
	)
	//initialise database variable, will be passed as parameters in methods
	driver_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/driver_db")
	if err != nil { //check for database opening errors
		driver_db.Close()
		log.Fatal(err.Error())
	}
	defer driver_db.Close()

	if r.Method == "GET" {
		driverIdList = getAllDriverIds(*driver_db)
		w.WriteHeader(http.StatusOK)
		output, _ := json.Marshal(map[string][]int{"Driver IDs": driverIdList})

		fmt.Fprintf(w, "%s", output)
	}
}

// func loginDriver(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type","application/json")
// 	var (
// 		query                     = r.URL.Query()
// 		lookupMobileNo     string = query.Get("mobileNo")
// 		searchEmailAddress string = query.Get("emailAddress")
// 		driverFound     bool
// 	)

// 	if r.Method == "GET" {

// 		for _, driver := range driverList {
// 			if lookupMobileNo != "" {
// 				//write database code to query the database
// 				if passenger.MobileNumber == lookupMobileNo {

// 					driverFound = true
// 					w.WriteHeader(http.StatusOK)
// 					fmt.Fprintf(w, "Driver %v is logged in successfully.", passenger.FirstName+passenger.LastName)
// 				}
// 			} else if searchEmailAddress != "" {
// 				if passenger.EmailAddress == searchEmailAddress {

// 					driverFound = true
// 					w.WriteHeader(http.StatusOK)
// 					fmt.Fprintf(w, "Driver %v is logged in successfully.", passenger.FirstName+passenger.LastName)
// 				}
// 			}

// 		}
// 		if !driverFound {
// 			w.WriteHeader(http.StatusNotFound)
// 			fmt.Fprintf(w, "Invalid Driver Mobile Number")
// 		}
// 	}
// }

//database functions

func getAllDriversRecords(driver_db sql.DB) {

	results, err := driver_db.Query("SELECT * FROM Drivers")
	if err != nil {
		log.Fatal(err.Error())
	}
	for results.Next() {
		var d Driver
		results.Scan(&d.driver_id, &d.FirstName, &d.LastName, &d.EmailAddress, &d.MobileNumber, &d.IdNo, &d.CarLicenseNo)
		driverList = append(driverList, d)
	}
}
func getDriverRecord(driver_db sql.DB, driver_id int) Driver {
	var d Driver
	if results := driver_db.QueryRow("SELECT driver_first_name,driver_last_name, driver_email, driver_mobileNo,driver_licenseNo FROM Drivers WHERE driver_id = ?", driver_id); results.Err() == nil {
		results.Scan(&d.FirstName, &d.LastName, &d.EmailAddress, &d.MobileNumber, &d.CarLicenseNo)
		d.driver_id = driver_id
	}
	return d
}
func insertDriverRecord(d Driver, driver_db sql.DB) (int, int64) {
	var (
		status int = 0 //1 represents successful insertion, 0 represents failure

	)
	result, err := driver_db.Exec("INSERT INTO Drivers (driver_first_name, driver_last_name, driver_email, driver_mobileNo, driver_idNo, driver_licenseNo) VALUES(?,?,?,?,?,?)", strings.TrimSpace(d.FirstName), strings.TrimSpace(d.LastName), strings.ToLower(strings.TrimSpace(d.EmailAddress)), strings.TrimSpace(d.MobileNumber), strings.TrimSpace(d.IdNo), strings.ToUpper(strings.TrimSpace(d.CarLicenseNo)))
	if err != nil {
		fmt.Println(err.Error())

	}

	driverId, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		driverId = 0
	}
	if driverId > 0 {
		status = 1
	}
	return status, driverId

}

func updateDriverRecord(driver map[string]string, driver_db sql.DB) int {
	var updateStatus int //defaults to 0 --> denotes failed update query execution
	id, _ := strconv.Atoi(driver["Driver Id"])
	results, err := driver_db.Exec("UPDATE Drivers SET driver_first_name = ?,driver_last_name = ?,driver_email = ?,driver_mobileNo = ?,driver_licenseNo = ? WHERE driver_id = ?", strings.TrimSpace(driver["First Name"]), strings.TrimSpace(driver["Last Name"]), strings.ToLower(strings.TrimSpace(driver["Email Address"])), strings.TrimSpace(driver["Mobile Number"]), strings.ToUpper(strings.TrimSpace(driver["Car License No"])), id)
	if err != nil {
		log.Fatal(err.Error())

	}
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		log.Fatal(err.Error())
	}
	if rowsAffected > 0 {
		updateStatus = 1 // 1 denotes successful Update operation
	}
	return updateStatus
}

func getAllDriverIds(driver_db sql.DB) []int {
	var (
		idList []int
		id     int
	)
	if results, err := driver_db.Query("SELECT driver_id from Drivers"); err == nil {
		for results.Next() {

			results.Scan(&id)
			idList = append(idList, id)
		}
	}
	return idList
}
