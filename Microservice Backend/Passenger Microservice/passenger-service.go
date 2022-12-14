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

type Passenger struct {
	Passenger_id int    `json:"Passenger Id"`
	FirstName    string `json:"First Name"`
	LastName     string `json:"Last Name"`
	MobileNumber string `json:"Mobile Number"`
	EmailAddress string `json:"Email Address"`
}

var passengerList []Passenger

// API Status Codes
// http.StatusMethodNotAllowed -- 405 -- For invalid API methods e.g. DELETE Request on a PUT Endpoint
// http.StatusConflict -- 409 -- duplicate record found in database
// http.StatusInternalServerError -- 500 -- Database query execution error or JSON payload error
// http.StatusAccepted - 201 -- successful inserted record in database
// http.StatusOK -- 200 - successfully updated record in
// http.StatusNotFound -- 404 -- passenger record not found in database
// http.StatusBadRequest -- 400 -- invalid request payload e.g. empty inputs provided

func main() {

	router := mux.Router{}
	router.HandleFunc("/api/v1/passenger", passengerSignup).Methods("POST")
	router.HandleFunc("/api/v1/passenger/details", passengerLogin).Methods("GET")
	router.HandleFunc("/api/v1/passenger", passengerUpdate).Methods("PUT")
	fmt.Println("Listening at port 5000...")
	log.Fatal(http.ListenAndServe(":5000", &router))
}

func passengerSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		searchEmailAddress, lookupMobileNo string //will be used to perform search operations
		invalidPassenger                   bool
		newPassengerData                   Passenger //object format
	)

	//initialise database variable, will be passed as parameters in methods
	passenger_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/passenger_db")
	if err != nil { //check for database opening errors
		passenger_db.Close()
		log.Fatal(err.Error())
	}
	defer passenger_db.Close()

	//synchronise passengerList with latest copy of Passenger records
	getAllPassengersRecord(*passenger_db)

	//handle signup function
	if r.Method == "POST" {
		if body, err := io.ReadAll(r.Body); err == nil {

			if err := json.Unmarshal(body, &newPassengerData); err == nil {

				searchEmailAddress, lookupMobileNo = newPassengerData.EmailAddress, newPassengerData.MobileNumber //will be used to perform search operations

				for _, obj := range passengerList {
					// search condition 1: Check for existing Matching email address
					if searchEmailAddress == obj.EmailAddress {
						invalidPassenger = true
						w.WriteHeader(http.StatusConflict) // duplicate passenger profile found
						fmt.Fprintf(w, "%s", bytes.NewBufferString("Email Address "+searchEmailAddress+" used by another Passenger"))

						break
					}
					// search condition 2: Check for existing Matching Mobile Number
					if lookupMobileNo == obj.MobileNumber {
						invalidPassenger = true
						w.WriteHeader(http.StatusConflict) // duplicate passenger profile found
						fmt.Fprintf(w, "%s", bytes.NewBufferString("Mobile Number "+lookupMobileNo+" used by another Passenger"))
						break
					}
				}
				if !invalidPassenger { // normal flow --> all unique values

					//insert passenger record to database
					status := insertPassengerRecord(newPassengerData, *passenger_db)
					getAllPassengersRecord(*passenger_db) //refresh passenger List after insert operation
					// status outcome of insert operation is to be notified to passenger
					switch status {
					case 0:
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, "%s", bytes.NewBufferString("Database error signing up passenger "+" "+newPassengerData.FirstName+" "+newPassengerData.LastName))

					case 1:
						w.WriteHeader(http.StatusAccepted)
						fmt.Fprintf(w, "%s", bytes.NewBufferString("Passenger "+newPassengerData.FirstName+" "+newPassengerData.LastName+" has been successfully signed up"))

					}
				}
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Error signing up passenger : %v", err)
			fmt.Fprintf(w, "%s", bytes.NewBufferString("Error signing up passenger"))
		}

	} else { // restrict API operations by other HTTP protocols

		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Invalid API Operation, only POST requests Allowed")
	}
}
func passengerUpdate(w http.ResponseWriter, r *http.Request) {

	//initialise database variable, will be passed as parameters in methods
	passenger_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/passenger_db")
	if err != nil { //check for database opening errors
		passenger_db.Close()
		log.Fatal(err.Error())
	}
	defer passenger_db.Close()

	w.Header().Set("Content-Type", "application/json")

	//refresh passengerList records to get latest data
	getAllPassengersRecord(*passenger_db)
	var (
		updatedMap, oldMap map[string]string //updatedMap store new details, oldMap store old details
		passengerFound     bool
		query              = r.URL.Query()

		//get query string arguments
		lookupMobileNo     string = query.Get("mobileNo")     // search for existing Passenger record by mobile number
		lookupEmailAddress string = query.Get("emailAddress") // search for existing Passenger record by email address

	)

	if r.Method == "PUT" {
		if body, err := io.ReadAll(r.Body); err == nil {

			json.Unmarshal(body, &updatedMap)
			for _, passenger := range passengerList {

				if (lookupMobileNo == passenger.MobileNumber) || (lookupEmailAddress == passenger.EmailAddress) { // mobile number found OR email address found
					p, _ := json.Marshal(passenger)
					json.Unmarshal(p, &oldMap)
					passengerFound = true
					//passenger_id is updated for the updatedMap object (new Passenger object) as updatedMap passenger_id is initialised to 0
					//id is required for "WHERE" clause parameter during update operation

					//identify * fields --> fields that remain unchanged will be assigned values of old passenger object
					for k, v := range updatedMap {
						if strings.Contains(v, "*") {
							updatedMap[k] = oldMap[k]
						}
					}
					updatedMap["Passenger Id"] = strconv.Itoa(passenger.Passenger_id)
					status := passengerUpdateRecord(updatedMap, *passenger_db) //update passenger record in database

					switch status { //format http response body based on search and update status
					case 0:
						w.WriteHeader(http.StatusInternalServerError)

						fmt.Fprintf(w, "%s", bytes.NewBufferString("Database error updating passenger "+passenger.FirstName+" "+passenger.LastName)) //display the old First Name and Last Name to inform Client whose record could'nt be updated

					case 1:
						w.WriteHeader(http.StatusOK)
						fmt.Fprintf(w, "%s", bytes.NewBufferString("Passenger record successfully updated"))

					}
					break
				}

			}
			//format http response body based on invalid result type
			if !passengerFound && lookupMobileNo != "" { //invalid mobile number provided
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "%s", bytes.NewBufferString("Passenger with the mobile number "+lookupMobileNo+" not found"))
			} else if !passengerFound && lookupEmailAddress != "" {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "%s", bytes.NewBufferString("Passenger with the email address "+lookupEmailAddress+" not found"))
			}
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed) //error 405 for invalid operation
		fmt.Fprintf(w, "%s", "Invalid API Operation")
	}
}

func passengerLogin(w http.ResponseWriter, r *http.Request) {

	//initialise database variable, will be passed as parameters in methods
	passenger_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/passenger_db")
	if err != nil { //check for database opening errors
		passenger_db.Close()
		log.Fatal(err.Error())
	}
	defer passenger_db.Close()

	w.Header().Set("Content-Type", "application/json")
	var (
		query                     = r.URL.Query()
		searchMobileNo     string = query.Get("mobileNo")
		searchEmailAddress string = query.Get("emailAddress")
	)
	searchId, _ := strconv.Atoi(query.Get("passengerId"))
	if r.Method == "GET" {
		if searchMobileNo != "" || searchEmailAddress != "" || searchId > 0 { //verify that at least one search value is passed in
			if passenger, found := passengerExists(searchMobileNo, searchEmailAddress, searchId, *passenger_db); passenger.Passenger_id > 0 || found { // verify passenger id found from database
				w.WriteHeader(http.StatusOK)
				if passenger.Passenger_id == 0 { //formatting json output for requesting passenger details from driver console
					passengerMap := map[string]string{"First Name": passenger.FirstName, "Last Name": passenger.LastName, "Mobile Number": passenger.MobileNumber, "Email Address": passenger.EmailAddress}
					out, _ := json.Marshal(map[string]map[string]string{"Passenger Record": passengerMap})

					fmt.Fprintf(w, "%v", bytes.NewBuffer(out))
				} else { //formatting json output for requesting passenger details from passenger microservice
					out, _ := json.Marshal(map[string]Passenger{"Passenger Record": passenger})

					fmt.Fprintf(w, "%v", bytes.NewBuffer(out))
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "%s", bytes.NewBufferString("Passenger Not Found"))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s", bytes.NewBufferString("Invalid Input. Passenger's Email Address or mobile number has to be provided."))

		}
	}
}

// database functions
func passengerExists(mobileNo string, emailAddress string, id int, passenger_db sql.DB) (Passenger, bool) {
	var (
		found     bool
		results   *sql.Row
		passenger Passenger
	)
	if id > 0 { //passengerID provided
		results = passenger_db.QueryRow("Select passenger_first_name, passenger_last_name, passenger_email, passenger_mobileNo FROM Passengers WHERE passenger_id = ?", id)
		if results.Err() == nil {

			if err := results.Scan(&(passenger.FirstName), &(passenger.LastName), &(passenger.EmailAddress), &(passenger.MobileNumber)); err == nil {
				found = true
			}
		}

	} else { //email or MobileNumber provided
		results = passenger_db.QueryRow("Select * FROM Passengers WHERE passenger_mobileNo =? OR passenger_email = ? OR passenger_id = ?", mobileNo, emailAddress, id)
		if results.Err() == nil {

			if err := results.Scan(&(passenger.Passenger_id), &(passenger.FirstName), &(passenger.LastName), &(passenger.EmailAddress), &(passenger.MobileNumber)); err == nil {
				found = true
			}
		}

	}
	return passenger, found
}
func getAllPassengersRecord(passenger_db sql.DB) {

	results, err := passenger_db.Query("SELECT * FROM Passengers")
	if err != nil {
		log.Fatal(err.Error())
	}
	for results.Next() {
		var p Passenger
		results.Scan(&p.Passenger_id, &p.FirstName, &p.LastName, &p.EmailAddress, &p.MobileNumber)
		passengerList = append(passengerList, p)
	}
}

func insertPassengerRecord(p Passenger, passenger_db sql.DB) int {
	var status int = 0 //1 represents successful insertion, 0 represents failure
	result, err := passenger_db.Exec("INSERT INTO Passengers (passenger_first_name,passenger_last_name,passenger_email,passenger_mobileNo) VALUES(?,?,?,?)", strings.TrimSpace(p.FirstName), strings.TrimSpace(p.LastName), strings.ToLower(strings.TrimSpace(p.EmailAddress)), strings.TrimSpace(p.MobileNumber))
	if err != nil {
		fmt.Println(err.Error())

	}

	passengerId, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())

	}
	if passengerId > 0 {
		status = 1
	}
	return status

}

func passengerUpdateRecord(p map[string]string, passenger_db sql.DB) int {
	var updateStatus int //defaults to 0 --> denotes failed update query execution
	id, _ := strconv.Atoi(p["Passenger Id"])
	results, err := passenger_db.Exec("UPDATE Passengers SET passenger_first_name = ?,passenger_last_name = ?,passenger_email = ?,passenger_mobileNo = ? WHERE passenger_id = ?", strings.TrimSpace(p["First Name"]), strings.TrimSpace(p["Last Name"]), strings.ToLower(strings.TrimSpace(p["Email Address"])), strings.TrimSpace(p["Mobile Number"]), id)
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
