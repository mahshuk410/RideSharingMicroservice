package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TwiN/go-color"
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

type Trip struct {
	TripId         int    `json:"Trip Id"`
	FromPostalCode int64  `json:"From"`
	ToPostalCode   int64  `json:"To"`
	TripStatus     string `json:"Trip Status"`
	PassengerId    int    `json:"Passenger Id"`
	DriverId       int    `json:"Driver Id"`
}

var (
	driverUrl      string = "http://localhost:5002/api/v1/driver"
	tripServiceUrl string = "http://localhost:5003/api/v1/trips/driver"
)

func main() {
	//create new http client object

	client := http.Client{}

	for {
		// Display Driver SignIn Menu
		displayDriverMenu(client)

	}
}
func displayDriverMenu(client http.Client) {
	var option int

	fmt.Println(color.Ize(color.Cyan, "\nRide Sharing App[Driver]\n===========================\n[1] Create a New Account\n[2] Update Existing Account Details\n[3] Go Online / Start Trip / Resume Trip\n[4] End Trip\n[5] Exit")) //line to display driver's menu
	fmt.Print("Enter an option:")
	fmt.Scanf("%d", &option)
	// execute functions based on user input
	//all time.sleep functions is to slow down program execution and regualate pace of asynchronous API requests
	switch option {
	case 1:
		driverCreateNewAccount(client)
	case 2:
		driverUpdateAccount(client)
	case 3:
		searchTrip(client)
	case 4:
		endTrip(client)

	case 5:
		exit()
	}
}
func driverCreateNewAccount(client http.Client) {
	var (
		newDriver    Driver = Driver{}
		responseData map[string]int64
	)

	//create a new driver object after gathering user inputs
	fmt.Print("Enter your First Name:")
	fmt.Scan(&(newDriver.FirstName))

	fmt.Print("\nEnter your Last Name:")
	fmt.Scan(&(newDriver.LastName))

	fmt.Print("\nEnter your Mobile Number:")
	fmt.Scan(&(newDriver.MobileNumber))

	fmt.Print("\nEnter your Email Address:")
	fmt.Scan(&(newDriver.EmailAddress))

	fmt.Print("\nEnter your Identification Number:")
	fmt.Scan(&(newDriver.IdNo))

	fmt.Print("\nEnter your car license Number:")
	fmt.Scan(&(newDriver.CarLicenseNo))
	//POST the Driver data over API

	//convert to byte array

	postBody, _ := json.Marshal(newDriver)
	reqBody := bytes.NewBuffer(postBody)
	if req, err := http.NewRequest(http.MethodPost, driverUrl, reqBody); err == nil {
		if res, err := client.Do(req); err == nil {
			loadingBar("Creating Driver Account..")
			if body, err := io.ReadAll(res.Body); err == nil {
				if res.StatusCode == http.StatusAccepted { //successful insertion
					json.Unmarshal(body, &responseData)
					fmt.Printf("Driver account has been successfully created.\nNote your Driver Id:%v\n", responseData["New Driver Id"])
				} else {
					response := string(body) // restrained from json.Unmarshal() method as conversion does not read response body string
					fmt.Println(response)

				}
				time.Sleep(3 * time.Second)
			}

		} else {

			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
}

func driverUpdateAccount(client http.Client) {

	var (
		lookupEmailAddress string
		lookupMobileNo     string
		updateDriver       Driver = Driver{}
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
	fmt.Println(color.Colorize(color.Red, "ENTER '*' IF UNCHANGED"))

	fmt.Print("Enter your First Name:")
	fmt.Scan(&(updateDriver.FirstName))

	fmt.Print("\nEnter your Last Name:")
	fmt.Scan(&(updateDriver.LastName))

	fmt.Print("\nEnter your Mobile Number:")
	fmt.Scan(&(updateDriver.MobileNumber))

	fmt.Print("\nEnter your Email Address:")
	fmt.Scan(&(updateDriver.EmailAddress))
	//idNo field left out to restrict updating of idNo attribute
	fmt.Print("\nEnter your car license Number:")
	fmt.Scan(&(updateDriver.CarLicenseNo))

	postBody, _ := json.Marshal(updateDriver)
	reqBody := bytes.NewBuffer(postBody)
	if req, err := http.NewRequest(http.MethodPut, driverUrl, reqBody); err == nil {
		//set query strings to search for existing Driver records
		query := req.URL.Query()
		query.Add("mobileNo", lookupMobileNo)
		query.Add("emailAddress", lookupEmailAddress)
		req.URL.RawQuery = query.Encode()

		if res, err := client.Do(req); err == nil { //display output based on status outcome of API operation
			if body, err := io.ReadAll(res.Body); err == nil {
				loadingBar("Updating Driver Account Details..")
				response = string(body) //byte array converted directly to display response in string format as JSON.unmarshal cannot decode the string
				fmt.Println(response)
				time.Sleep(2 * time.Second)
			}
		}
	}

}
func searchTrip(client http.Client) {
	var (
		driver_id   int
		newTrip     Trip
		driverInput string
	)

	fmt.Print("Enter Your Driver Id:")
	fmt.Scan(&driver_id)
	// request trip function calls Trip API
	if req, err := http.NewRequest(http.MethodGet, tripServiceUrl+"/"+strconv.Itoa(driver_id), nil); err == nil && driver_id > 0 {
		if res, err := client.Do(req); err == nil {
			loadingBar("Searching For Trips")
			if res.StatusCode == 200 { //successfully retrieved trip

				body, _ := io.ReadAll(res.Body)
				json.Unmarshal(body, &newTrip)
				fmt.Printf("Yay! A trip has been assigned!\n=======================\nTripId: %v\nFrom: %v\nTo: %v\nTrip Status: %v\nPassengerId: %v\n", newTrip.TripId, newTrip.FromPostalCode, newTrip.ToPostalCode, newTrip.TripStatus, newTrip.PassengerId)
				if passengerData, found := getPassengerDetail(client, newTrip.PassengerId); found { //passenger found
					fmt.Println("Passenger Details\n=================")
					requestedPassenger := passengerData["Passenger Record"]
					fmt.Printf("Passenger Name: %v\nPassenger Email: %v\nPassenger Mobile:%v\n", requestedPassenger["First Name"]+" "+requestedPassenger["Last Name"], requestedPassenger["Email Address"], requestedPassenger["Mobile Number"])
				}
				if newTrip.TripStatus == "Started" { //for Ongoing trips after accidental logouts
					fmt.Println("You have picked up the passsenger and en route to destination..")
				} else { //yet to trigger start trip
					for !(strings.Contains(strings.ToLower(driverInput), "start")) { //wait until driver enters start
						fmt.Print("Enter 'Start' to start the trip:")
						fmt.Scan(&driverInput)
					}
					//start trip
					newTrip.TripStatus = "Started" //update cache
					//call Trips API to update status
					fmt.Println("Loading..")
					time.Sleep(2 * time.Second)

					updateTripStatus(client, strconv.Itoa(newTrip.DriverId), strconv.Itoa(newTrip.TripId), newTrip.TripStatus) //successfully started trip
					fmt.Println("Your trip has started..")

				}
				time.Sleep(2 * time.Second)
			} else if res.StatusCode == 404 { //no trips has been assigned yet
				// var input string
				fmt.Println("No trips found. Please try again later.")
				time.Sleep(2 * time.Second)
				fmt.Println("Returning to main menu...")
				time.Sleep(2 * time.Second)
			}
		}
	}
}
func endTrip(client http.Client) {
	var (
		currentDriverId int
		endInput        string
		ongoingTrip     Trip
	)
	fmt.Print("Enter Your Driver ID: ")
	fmt.Scan(&currentDriverId)

	if req, err := http.NewRequest(http.MethodGet, tripServiceUrl+"/"+strconv.Itoa(currentDriverId), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 200 { //successfully retrieved trip

				body, _ := io.ReadAll(res.Body)
				json.Unmarshal(body, &ongoingTrip)

				if ongoingTrip.TripStatus == "Started" && ongoingTrip.DriverId == currentDriverId { //validate that driver has a proper id

					fmt.Printf("Ongoing Trip\n===========\nTripId: %v\nFrom: %v\nTo: %v\nTrip Status: %v\nPassengerId: %v\n", ongoingTrip.TripId, ongoingTrip.FromPostalCode, ongoingTrip.ToPostalCode, ongoingTrip.TripStatus, ongoingTrip.PassengerId)

					fmt.Print("Enter 'End' to complete the trip:")
					fmt.Scan(&endInput)
					if strings.Contains(strings.ToLower(endInput), "end") {
						ongoingTrip.TripStatus = "Ended"
						updateTripStatus(client, strconv.Itoa(ongoingTrip.DriverId), strconv.Itoa(ongoingTrip.TripId), ongoingTrip.TripStatus) //successfully ended trip
						loadingBar("Ending The Trip..")                                                                                        //loading animation
						fmt.Println("Your trip has successfuly ended.")
						fmt.Println("Redirecting to Main Menu...")
						time.Sleep(2 * time.Second)

					}
				}
			} else {
				fmt.Println("No trips have been started yet")
			}
		}
	} else {
		fmt.Println(err)
	}

}
func updateTripStatus(client http.Client, driverId string, tripId string, status string) {
	var errorMessage string
	postBody := map[string]string{"Status": status}
	jsonBody, _ := json.Marshal(postBody)
	reqBody := bytes.NewBuffer(jsonBody)
	if req, err := http.NewRequest(http.MethodPut, tripServiceUrl+"/"+driverId+"/"+tripId, reqBody); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {

			} else if res.StatusCode == 500 {
				jsonErrorMessage, _ := io.ReadAll(res.Body)
				json.Unmarshal(jsonErrorMessage, &errorMessage)
				fmt.Println(errorMessage)
			}
		}
	}
}

func getPassengerDetail(client http.Client, passengerId int) (map[string]map[string]string, bool) {
	var (
		passengerData  map[string]map[string]string
		passengerFound bool
	)
	if req, err := http.NewRequest("GET", "http://localhost:5000/api/v1/passenger/details?passengerId="+strconv.Itoa(passengerId), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			resBody, _ := io.ReadAll(res.Body)
			if res.StatusCode == 200 {
				json.Unmarshal(resBody, &passengerData)
				passengerFound = true
			}

		}
	}
	return passengerData, passengerFound
}
func exit() {
	os.Exit(1)
	fmt.Println("Exiting Ride-Sharing App....")
}

//loading animation

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
