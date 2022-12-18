package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/TwiN/go-color"
)

type Passenger struct {
	Passenger_id int    `json:"Passenger Id"`
	FirstName    string `json:"First Name"`
	LastName     string `json:"Last Name"`
	Mobilenumber string `json:"Mobile Number"`
	EmailAddress string `json:"Email Address"`
}

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
	passengerServiceUrl string = "http://localhost:5000/api/v1/passenger"
	tripServiceUrl      string = "http://localhost:5003/api/v1/trips"
	driverServiceUrl    string = "http://localhost:5002/api/v1/driver"
)

func main() {
	//create new http client object

	client := http.Client{}
	for {
		// Display Passenger Menu
		displayPassengerMenu(client)

	}

}
func displayPassengerMenu(client http.Client) {
	var option int
	fmt.Println(color.Ize(color.Yellow, "\nRide Sharing App[passenger]\n===========================\n[1] Create a New Account\n[2] Update Existing Account Details\n[3] Request New Trip / Continue Ongoing Trip \n[4] View Trip History \n[5] Exit"))
	fmt.Print("Enter an option:")
	fmt.Scanf("%d", &option)
	// execute functions based on user input
	//all time.sleep functions is to slow down program execution and regualate pace of asynchronous API requests
	switch option {
	case 1:
		passengerCreateNewAccount(client)
	case 2:
		passengerUpdateAccount(client)
	case 3:
		requestTrip(client)
	case 4:
		viewTripHistory(client)
	case 5:
		exit()
	}

}

func passengerCreateNewAccount(client http.Client) {
	var (
		newPassenger Passenger
		//data string
		//input variables to create passenger object
		FirstName, LastName, Mobilenumber, EmailAddress string
	)
	//create a new passenger object after gathering user inputs
	fmt.Print("Enter your First Name:")

	fmt.Scan(&FirstName)
	fmt.Print("\nEnter your Last Name:")
	fmt.Scan(&LastName)
	fmt.Print("\nEnter your Mobile Number:")
	fmt.Scan(&Mobilenumber)
	fmt.Print("\nEnter your Email Address:")
	fmt.Scan(&EmailAddress)
	newPassenger = Passenger{FirstName: FirstName, LastName: LastName, Mobilenumber: Mobilenumber, EmailAddress: EmailAddress}
	//POST the Passenger data over API
	//convert to byte array
	postBody, _ := json.Marshal(newPassenger)
	reqBody := bytes.NewBuffer(postBody)
	if req, err := http.NewRequest(http.MethodPost, passengerServiceUrl, reqBody); err == nil {
		if res, err := client.Do(req); err == nil {
			if body, err := io.ReadAll(res.Body); err == nil {
				loadingBar("Creating Passenger Account..")
				response := string(body) // restrained from json.Unmarshal() method as conversion does not read response body string

				fmt.Println(response)
				time.Sleep(2 * time.Second)
			}
			// //if-else block prints out response messages to display outcome of passenger sign-up
			// if (res.StatusCode == http.StatusAccepted){ //successful
			// 	// if body,err := io.ReadAll(res.Body); err ==nil{

			// 	// 	json.Unmarshal(body,&response)

			// 	// 	fmt.Println(response)
			// 	// }

			// }else if (res.StatusCode == http.StatusConflict){ //failure probably due to duplicated information e.g. re-used email address or mobile Number
			// 	var output string
			// 	if body,err := io.ReadAll(res.Body);err ==nil{

			// 		json.Unmarshal(body,&output)
			// 		fmt.Println(output)
			// 	}

			// }
		} else {

			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
}

func passengerExists(client http.Client) (bool, int) {
	var (
		passengerFound                         bool
		passengerId, loginOption               int
		searchMobileNumber, searchEmailAddress string
	)
	fmt.Print(color.Colorize(color.Blue, "\n=============\nLogin Methods\n=============\nLogin by Mobile Number[1]\nLogin by Email Address[2]\nOption:"))

	fmt.Scan(&loginOption)

	switch loginOption {
	case 1:
		fmt.Print("Enter Mobile Number:")
		fmt.Scan(&searchMobileNumber)
	case 2:
		fmt.Print("Enter Email Address:")
		fmt.Scan(&searchEmailAddress)
	}

	if req, err := http.NewRequest(http.MethodGet, passengerServiceUrl+"/details", nil); err == nil {
		query := req.URL.Query()
		query.Add("emailAddress", searchEmailAddress)
		query.Add("mobileNo", searchMobileNumber)
		req.URL.RawQuery = query.Encode()
		if res, err := client.Do(req); err == nil {
			body, _ := io.ReadAll(res.Body)
			if res.StatusCode == 400 || res.StatusCode == 404 { //unsuccessful login

				fmt.Printf("%s", body) //output error message
				exit()
			} else if res.StatusCode == 200 {
				var mapOutput map[string]Passenger
				json.Unmarshal(body, &mapOutput)
				passengerId = mapOutput["Passenger Record"].Passenger_id
				passengerFound = true
			}
		}
	}
	return passengerFound, passengerId
}
func passengerUpdateAccount(client http.Client) {

	var (
		lookupEmailAddress string
		lookupMobileNo     string
		updatePassenger    Passenger = Passenger{}
		searchMethod       int64
		response           string
	)
	for !(searchMethod == 1 || searchMethod == 2) {
		fmt.Print("[1]Search by Email Address\n[2]Search By Mobile Number\n=========================\nOption: ")
		fmt.Scan(&searchMethod)
	}

	switch searchMethod {
	case 1: // email option selected
		fmt.Print("\nEnter your existing Email Address:")
		fmt.Scan(&lookupEmailAddress)
	case 2: // mobile option selected
		fmt.Print("\nEnter your existing Mobile Number:")
		fmt.Scan(&lookupMobileNo)
	default:
		fmt.Println("Invalid Option")
	}
	// get existing passenger details for passenger's reference
	if req, err := http.NewRequest(http.MethodGet, passengerServiceUrl+"/details", nil); err == nil {
		query := req.URL.Query()
		query.Add("emailAddress", lookupEmailAddress)
		query.Add("mobileNo", lookupMobileNo)
		req.URL.RawQuery = query.Encode()
		if res, err := client.Do(req); err == nil {
			body, _ := io.ReadAll(res.Body)
			if res.StatusCode == 400 || res.StatusCode == 404 { //unsuccessful login

				fmt.Printf("%s", body)                 //output error message
				loadingBar("Returning to Main Menu..") //exit program if invalid record
			} else if res.StatusCode == 200 {
				var mapOutput map[string]Passenger
				json.Unmarshal(body, &mapOutput)
				updatePassenger = mapOutput["Passenger Record"]
			}
		}
	}

	fmt.Println("\nExisting Passenger Details\n==========================")
	fmt.Printf("First Name: %v\nLast Name: %v\nMobile Number:%v\nEmail Address:%v\n", updatePassenger.FirstName, updatePassenger.LastName, updatePassenger.Mobilenumber, updatePassenger.EmailAddress)
	//===========================================================
	fmt.Println(color.Colorize(color.Red, "ENTER '*' IF UNCHANGED"))
	fmt.Print("Enter your new First Name:")

	fmt.Scan(&(updatePassenger.FirstName))
	fmt.Print("\nEnter your new Last Name:")
	fmt.Scan(&(updatePassenger.LastName))
	fmt.Print("\nEnter your new Mobile Number:")
	fmt.Scan(&(updatePassenger.Mobilenumber))
	fmt.Print("\nEnter your new Email Address:")
	fmt.Scan(&(updatePassenger.EmailAddress))

	postBody, _ := json.Marshal(updatePassenger)
	reqBody := bytes.NewBuffer(postBody)
	if req, err := http.NewRequest(http.MethodPut, passengerServiceUrl, reqBody); err == nil {
		//set query strings to search for existing Passenger records
		query := req.URL.Query()
		query.Add("mobileNo", lookupMobileNo)
		query.Add("emailAddress", lookupEmailAddress)
		req.URL.RawQuery = query.Encode()

		if res, err := client.Do(req); err == nil { //display output based on status outcome of API operation
			if body, err := io.ReadAll(res.Body); err == nil {
				response = string(body) //byte array converted directly to display response in string format as JSON.unmarshal cannot decode the string
				fmt.Println(response)
				time.Sleep(2 * time.Second)
			}
		}
	}

}

// request trip function calls Trip API
func requestTrip(client http.Client) {
	var (
		fromPostalCode, toPostalCode int64
		tripDetails                  map[string]Trip
	)

	if passengerFound, passengerId := passengerExists(client); passengerFound {
		if tripData, tripFound := getLatestOngoingTrip(client, passengerId); tripFound { //check whether user has any ongoing rides in case of accidental logout
			ongoingTrip := tripData["Trip"]
			fmt.Printf("Trip Start Time: %v\nFrom: %v\nTo: %v\n", ongoingTrip.StartTime, ongoingTrip.FromPostalCode, ongoingTrip.ToPostalCode)
			fmt.Println("You are on your way to your destination..")

			checkTripForEnded(client, ongoingTrip.TripId)
		} else { //no ongoing rides

			for fromPostalCode == toPostalCode { //validate different From and To Postal Codes
				fmt.Print("From(postal code):")
				fmt.Scan(&fromPostalCode)
				fmt.Print("\nTo(postal code):")
				fmt.Scan(&toPostalCode)
			}
			postalCodeData := map[string]int64{"From": fromPostalCode, "To": toPostalCode} //postal codes to be sent as request body
			postBody, _ := json.Marshal(postalCodeData)
			reqBody := bytes.NewBuffer(postBody)

			if req, err := http.NewRequest(http.MethodPost, tripServiceUrl+"/passenger/"+strconv.Itoa(passengerId), reqBody); err == nil {

				if res, err := client.Do(req); err == nil {
					loadingBar("Searching For a Driver..")
					resBody, _ := io.ReadAll(res.Body)
					if res.StatusCode == 202 { //ride has been found
						json.Unmarshal(resBody, &tripDetails)
						driverInfo := getDriverDetails(tripDetails["Created Trip"].DriverId, client)

						fmt.Println("A driver has been found. Please wait for his arrival :)")
						fmt.Printf("===============\nDriver Details\n===============\n Driver Name: %v\n Driver Mobile Number: %v\n Driver Email: %v\n Driver Car No.: %v\n", driverInfo["First Name"]+driverInfo["Last Name"], driverInfo["Mobile Number"], driverInfo["Email Address"], driverInfo["Car License No"])
						checkTripForStarted(client, tripDetails["Created Trip"].TripId)

						fmt.Println("Enjoy your Ride!") //display once ride has started

						checkTripForEnded(client, tripDetails["Created Trip"].TripId)

					} else { //ride not found
						var errorMessage string
						json.Unmarshal(resBody, &errorMessage)
						fmt.Println(errorMessage)
					}
				} else { //response error
					fmt.Println(err.Error())
				}
			} else { //request error
				fmt.Println(err.Error())
			}

		}

	} else {
		fmt.Println("Passenger Not Found. Please try again.")
	}
	loadingBar("Returning to Main Menu...")
}
func checkTripForStarted(client http.Client, tripId int) { //Trip API function to continuously check from db if driver has started trip
	var tripStarted bool
	for !tripStarted { //repeatedly make API Calls every 3 seconds to validate if trip has been started by driver
		if tripStatus := getTripStatus(client, tripId); tripStatus == "Started" {
			fmt.Println("Woohoo!Your trip has started.")
			tripStarted = true
		} else {
			time.Sleep(3 * time.Second)
		}
	}

}
func checkTripForEnded(client http.Client, tripId int) { //Trip API function to continuously check from trip db if driver has ended trip
	var tripEnded bool
	for !tripEnded { //repeatedly make API Calls every 3 seconds to validate if trip has been end by driver
		if tripStatus := getTripStatus(client, tripId); tripStatus == "Ended" {
			fmt.Println("Yay!Your trip has ended..")
			tripEnded = true

		}
		time.Sleep(3 * time.Second)
	}

}
func getTripStatus(client http.Client, trip_id int) string {
	var statusMap map[string]string
	if req, err := http.NewRequest(http.MethodGet, tripServiceUrl+"/"+strconv.Itoa(trip_id), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 200 {
				resBody, _ := io.ReadAll(res.Body)
				json.Unmarshal(resBody, &statusMap)
				return statusMap["Status"]
			}
		}
	}
	return ""
}
func getLatestOngoingTrip(client http.Client, passengerId int) (map[string]Trip, bool) {
	var tripHistory map[string][]Trip
	var tripFound bool
	var onGoingTrip map[string]Trip
	if req, err := http.NewRequest(http.MethodGet, tripServiceUrl+"/passenger/"+strconv.Itoa(passengerId), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 200 { //trips history is found
				tripData, _ := io.ReadAll(res.Body)
				json.Unmarshal(tripData, &tripHistory)
				for _, trip := range tripHistory["Trips"] {
					// tripStartTime,_ := time.Parse("YYYY--MM-DD hh:mm:ss",trip.StartTime)
					if trip.TripStatus == "Started" {
						onGoingTrip = map[string]Trip{"Trip": trip}
						tripFound = true
						break
					}

				}

			}

		}
	}
	return onGoingTrip, tripFound

}
func viewTripHistory(client http.Client) {
	var tripHistory map[string][]Trip
	if passengerFound, passengerId := passengerExists(client); passengerFound { //valid passenger found
		if req, err := http.NewRequest(http.MethodGet, tripServiceUrl+"/passenger/"+strconv.Itoa(passengerId), nil); err == nil {
			if res, err := client.Do(req); err == nil {
				if res.StatusCode == 200 { //trips history is found
					tripData, _ := io.ReadAll(res.Body)
					json.Unmarshal(tripData, &tripHistory)
					fmt.Print("############\nTrip History\n############")
					for tripNo, trip := range tripHistory["Trips"] {
						fmt.Printf("\n=======Trip %v=======\n", tripNo+1)
						fmt.Printf("TripId: %v\nFrom:%v\nTo:%v\nTrip Status:%v\nStart Time:%v\nEnd Time:%v\n\n", trip.TripId, trip.FromPostalCode, trip.ToPostalCode, trip.TripStatus, trip.StartTime, trip.EndTime)

					}

					loadingBar("Returning to Main Menu...")

				} else if res.StatusCode == 404 { // no trips found

					errorMessage, _ := io.ReadAll(res.Body)

					fmt.Println(string(errorMessage)) //directly convert the byteArray to string for user to see error Message
				}

			}
		}
	} else { //passenger not found
		fmt.Println("Invalid Passenger ID")
	}

}
func exit() {
	os.Exit(1)
	fmt.Print("Exiting Ride-Sharing App....")
}

//Driver API function --> retrieve Assigned Driver Details

func getDriverDetails(driverId int, client http.Client) map[string]string {
	var (
		driverDetails map[string]string
	)
	if req, err := http.NewRequest("GET", driverServiceUrl+"/"+strconv.Itoa(driverId), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 200 {
				resBody, _ := io.ReadAll(res.Body)
				json.Unmarshal(resBody, &driverDetails)
			}
		}
	}
	return driverDetails
}

// loading animation function
func loadingBar(message string) {
	var loadingBar = []string{
		"00%: [                    ]",
		"20%: [####                ]",
		"40%: [########            ]",
		"60%: [############        ]",
		"80%: [################    ]",
		"100%:[####################]\n",
	}
	// Print Driver Search Progress Bar
	fmt.Println(message)

	for _, progress := range loadingBar {
		fmt.Printf("\r \a%s", color.Colorize(color.Green, progress))
		time.Sleep(1 * time.Second)
	}
}
