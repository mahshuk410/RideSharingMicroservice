package main

import (
	"bytes"
	"database/sql"

	"time"

	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"math/rand"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type Trip struct {
	TripId         int    `json:"Trip Id"`
	FromPostalCode int64  `json:"From"`
	ToPostalCode   int64  `json:"To"`
	TripStatus     string `json:"Trip Status"`
	StartTime      string `json:"Start Time"`
	EndTime        string `json:"End Time"`
	PassengerId    int    `json:"Passenger Id"`
	DriverId       int    `json:"Driver Id"`
}

var (
	onlinePassengerId, onlineDriverId int
	sourceDestinationMap              map[string]int

	client           http.Client //for requests to other services via APIs
	driverServiceUrl string      = "http://localhost:5002/api/v1/driver"
	allDriver        []int       //list maintaining all available drivers
)

func main() {
	router := mux.Router{}
	router.HandleFunc("/api/v1/trips/passenger/{passengerId}", passengerRequestTrip).Methods("POST", "GET")
	router.HandleFunc("/api/v1/trips/driver/{driverId}", driverGetTrip).Methods("GET")
	router.HandleFunc("/api/v1/trips/driver/{driverId}/{tripId}", driverUpdateTripStatus).Methods("PUT")
	router.HandleFunc("/api/v1/trips/{tripId}", getTripStatus).Methods("GET")
	fmt.Println("Listening at port 5003...")
	log.Fatal(http.ListenAndServe(":5003", &router))

}

func passengerRequestTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//initialise database variable, will be passed as parameters in methods
	trip_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/trip_db")
	if err != nil { //check for database opening errors

		log.Fatal(err.Error())
	}
	defer trip_db.Close()

	if r.Method == http.MethodPost {
		if body, err := io.ReadAll(r.Body); err == nil {
			if err := json.Unmarshal(body, &sourceDestinationMap); err == nil {
				//check if Source and destination address is the same
				if sourceDestinationMap["From"] == sourceDestinationMap["To"] {
					w.WriteHeader(http.StatusConflict)
					fmt.Fprintf(w, "%s", bytes.NewBufferString("Departure Address and Destination Address should be different."))
				}
				params := mux.Vars(r)
				onlinePassengerId, _ = strconv.Atoi(params["passengerId"])
				newTrip := Trip{FromPostalCode: int64(sourceDestinationMap["From"]), ToPostalCode: int64(sourceDestinationMap["To"]), TripStatus: "Pending", PassengerId: onlinePassengerId}
				allDriver = requestAllDriverId() //get id of all drivers
				if onlineDriverId = findAvailableDriver(*trip_db); onlineDriverId > 0 {
					newTrip.DriverId = onlineDriverId // assign random-generated Available Driver

					creationStatus, generated_trip_id := createTrip(newTrip, *trip_db)
					switch creationStatus {
					case 1: //trip successfully created once a driver is assigned
						newTrip.TripId = generated_trip_id
						w.WriteHeader(http.StatusAccepted)

						output, _ := json.Marshal(map[string]Trip{"Created Trip": newTrip})
						fmt.Fprintf(w, "%v", bytes.NewBuffer(output))
					case 0:
						w.WriteHeader(http.StatusInternalServerError) //database error inserting trip object
						fmt.Fprintf(w, "%s", bytes.NewBufferString("Error processing Trip request. Please try again"))
						newTrip = Trip{} //reset all values
					}

				} else { //no drivers available
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "%s", bytes.NewBufferString("No Drivers are found at the moment. Please try again."))

				}

			}

		}

	} else if r.Method == http.MethodGet {
		if passengerId, _ := strconv.Atoi(mux.Vars(r)["passengerId"]); passengerId > 0 {
			tripHistory := getPassengerTripHistory(*trip_db, passengerId)
			if len(tripHistory) > 0 {
				tripData, _ := json.Marshal(map[string][]Trip{"Trips": tripHistory})
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "%v", bytes.NewBuffer(tripData))

			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "%v", bytes.NewBufferString("No trips have been made so far."))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%v", bytes.NewBufferString("Please provide a valid Passenger ID"))
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%v", bytes.NewBufferString("Invalid API operation"))
	}
}
func driverGetTrip(w http.ResponseWriter, r *http.Request) {
	driverId, _ := strconv.Atoi(mux.Vars(r)["driverId"])

	//initialise database variable, will be passed as parameters in methods
	trip_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/trip_db")
	if err != nil { //check for database opening errors

		log.Fatal(err.Error())
	}
	defer trip_db.Close()

	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		trip := getDriverAssignedTrip(*trip_db, driverId)
		if trip.TripId > 0 { //check if trip object was successfully retrieved
			output, _ := json.Marshal(trip)

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%v", bytes.NewBuffer(output))

		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%v", bytes.NewBufferString("No trips found"))
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s", bytes.NewBufferString("Invalid API operation"))
	}

}
func driverUpdateTripStatus(w http.ResponseWriter, r *http.Request) {

	//initialise database variable, will be passed as parameters in methods
	trip_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/trip_db")
	if err != nil { //check for database opening errors

		log.Fatal(err.Error())
	}
	defer trip_db.Close()

	var status map[string]string
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPut {
		param := mux.Vars(r)
		if body, err := io.ReadAll(r.Body); err == nil {
			json.Unmarshal(body, &status)
			driverId, _ := strconv.Atoi(param["driverId"])
			tripId, _ := strconv.Atoi(param["tripId"])

			if updateTripStatus(*trip_db, status["Status"], tripId, driverId) {
				w.WriteHeader(http.StatusAccepted)
				switch status["Status"] {
				case "Started":
					fmt.Fprintf(w, "%v", bytes.NewBufferString("Trip has been successfully Started"))
				case "Ended":
					fmt.Fprintf(w, "%v", bytes.NewBufferString("Trip has been successfully Ended"))
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%v", bytes.NewBufferString("Error updating Trip status.Please check the Driver ID and Trip ID"))
			}

		}
	} else { //invalid HTTP method
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s", bytes.NewBufferString("Invalid API operation"))

	}
}

func getTripStatus(w http.ResponseWriter, r *http.Request) {
	//initialise database variable, will be passed as parameters in methods
	trip_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/trip_db")
	if err != nil { //check for database opening errors

		log.Fatal(err.Error())
	}
	defer trip_db.Close()
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		tripId, _ := strconv.Atoi(mux.Vars(r)["tripId"])
		if status := getTripDetail(*trip_db, tripId); status != "" {
			output, _ := json.Marshal(map[string]string{"Status": status})
			resBody := bytes.NewBuffer(output)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%v", resBody)

		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%v", bytes.NewBufferString("Trip Status Not Found"))

		}
	}
}

// trip database functions
func createTrip(trip Trip, trip_db sql.DB) (int, int) {
	results, err := trip_db.Exec("INSERT INTO Trips (fromPostalCode, toPostalCode, tripStatus, passenger_id, driver_id) VALUES (?, ?, ?, ?, ?)", trip.FromPostalCode, trip.ToPostalCode, trip.TripStatus, trip.PassengerId, trip.DriverId)
	if err == nil {
		if trip_id, err := results.LastInsertId(); err == nil {
			trip.TripId = int(trip_id)
			return 1, trip.TripId
		}
	} else {
		fmt.Println(err.Error())
	}
	return 0, 0
}

func findAvailableDriver(trip_db sql.DB) int {
	var busyDrivers []int
	fmt.Println(time.Now())
	results, err := trip_db.Query("SELECT DISTINCT driver_id FROM Trips WHERE (startTime <= now() AND endTime >= now() AND tripStatus IN ('Started','Pending'))") //gets the busy drivers
	if err == nil {
		for results.Next() {
			var id int
			results.Scan(&id)
			busyDrivers = append(busyDrivers, id)
		}
		if len(busyDrivers) > 0 { //filter only if there are busy drivers
			for _, busyDriver := range busyDrivers {
				for index, driver := range allDriver {
					if driver == busyDriver {
						allDriver = append(allDriver[:index], allDriver[index+1:]...) //remove busy drivers from the driver list

					}
				}
			}
		}
		//generate random index number to assign a driver
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(allDriver))
		return allDriver[randomIndex]
	} else {
		fmt.Println(err.Error())
	}
	return 0
}
func getDriverAssignedTrip(trip_db sql.DB, driverId int) Trip {
	var (
		trip                         Trip
		TripId, PassengerId          int
		FromPostalCode, ToPostalCode int64
		TripStatus                   string
	)
	fmt.Println("here", driverId)
	results := trip_db.QueryRow("SELECT trip_id,fromPostalCode,toPostalCode,tripStatus,passenger_id FROM Trips WHERE driver_id = ? AND tripStatus IN('Started', 'Pending')", driverId)
	if results.Err() == nil {
		results.Scan(&TripId, &FromPostalCode, &ToPostalCode, &TripStatus, &PassengerId)
		trip = Trip{TripId: TripId, FromPostalCode: FromPostalCode, ToPostalCode: ToPostalCode, TripStatus: TripStatus, PassengerId: PassengerId, DriverId: driverId}
		return trip
	}
	return Trip{}

}
func updateTripStatus(trip_db sql.DB, status string, tripId int, driverId int) bool {
	var (
		updateSuccessful     bool
		updateQueryStatement string //modify dynamically. if starting trip,  set startTime to current time if ending trip set endTime to current time
	)

	if status == "Started" || status == "Ended" {
		switch status {
		case "Started":
			updateQueryStatement = "UPDATE Trips SET tripStatus=?,startTime = NOW() WHERE trip_id = ? AND driver_id =?"
		case "Ended":
			updateQueryStatement = "UPDATE Trips SET tripStatus=?,endTime = NOW() WHERE trip_id = ? AND driver_id =?"
		}
		results, err := trip_db.Exec(updateQueryStatement, status, tripId, driverId)
		if rows, err := results.RowsAffected(); rows > 0 {
			updateSuccessful = true
		} else {
			fmt.Println(err)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return updateSuccessful
}
func getTripDetail(trip_db sql.DB, trip_id int) string {
	var (
		TripStatus string
	)
	results := trip_db.QueryRow("SELECT tripStatus FROM Trips WHERE trip_id = ? ", trip_id)
	if results.Err() == nil {
		results.Scan(&TripStatus)
		return TripStatus
	}
	return ""
}

func getPassengerTripHistory(trip_db sql.DB, passenger_id int) []Trip {
	var tripHistory []Trip
	results, err := trip_db.Query("SELECT trip_id, fromPostalCode, toPostalCode, tripStatus, startTime, endTime, driver_id FROM Trips WHERE passenger_id=? ORDER BY startTime DESC", passenger_id)
	if err == nil {
		for results.Next() {
			var trip Trip
			results.Scan(&(trip.TripId), &(trip.FromPostalCode), &(trip.ToPostalCode), &(trip.TripStatus), &(trip.StartTime), &(trip.EndTime), &(trip.DriverId))
			trip.PassengerId = passenger_id
			tripHistory = append(tripHistory, trip)
		}
	} else {
		fmt.Println(err.Error())
	}
	return tripHistory
}

// API call to Driver API
func requestAllDriverId() []int {
	var driverIdList map[string][]int
	if req, err := http.NewRequest(http.MethodGet, driverServiceUrl+"/all", nil); err == nil {
		if res, err := client.Do(req); err == nil {
			resBody, err := io.ReadAll(res.Body)
			if err == nil {
				json.Unmarshal(resBody, &driverIdList)

				return driverIdList["Driver IDs"]

			}
		}
	}
	return nil
}
