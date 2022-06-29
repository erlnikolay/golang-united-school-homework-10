package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

/**
Please note Start functions is a placeholder for you to start your own solution.
Feel free to drop gorilla.mux if you want and use any other solution available.

main function reads host/port from env just for an example, flavor it following your taste
*/

type CalculationData struct {
	A string
	B string
}

// get request with param create struct param
func handleParam(w http.ResponseWriter, r *http.Request) {
	var paramsSlice []string

	// check kind of request
	// Get with param
	if r.Method == http.MethodGet {
		// check param
		paramsSlice = strings.Split(r.URL.Path, "/")
		fmt.Println("GET")
		fmt.Println(paramsSlice)
		fmt.Println(len(paramsSlice))
		if len(paramsSlice) < 2 {
			http.Error(w, "expect /name/{param} in task handler", http.StatusBadRequest)
			return
		}
		switch paramsSlice[1] {
		case "name":
			getWithParam(paramsSlice[len(paramsSlice)-1], w)
			return
		case "bad":
			getWithBad(w)
			return
		default:
			http.Error(w, "No headers set Response expected to have", http.StatusBadRequest)
			return
		}
	} else if r.Method == http.MethodPost { // Methos POST with body
		paramsSlice = strings.Split(r.URL.Path, "/")
		fmt.Println("POST")
		fmt.Println(paramsSlice)
		fmt.Println(len(paramsSlice))
		if len(paramsSlice) < 2 {
			http.Error(w, "expect /data or /header in task handler", http.StatusBadRequest)
			return
		} else {
			// take body into handler
			body, err := ioutil.ReadAll(r.Body)
			fmt.Printf("Request body: %v\n", string(body))
			fmt.Printf("Lenght of body: %v\n", len(body))
			if err != nil {
				fmt.Fprintf(w, "err %v %v\n", err, err.Error())
				return
			} else {
				// handle of post request
				// there is not a body
				if len(body) <= 0 {
					fmt.Fprintf(w, "No body set\n")
					return
				}
				switch paramsSlice[1] {
				case "data":
					if strings.Split(r.Header["Content-Type"][0], ";")[0] == "application/x-www-form-urlencoded" {
						postWithBodyAsData(body, w)
						return
					} else if strings.Split(strings.Split(r.Header["Content-Type"][0], ";")[0], " ")[0] == "multipart/form-data" {
						postWithBodyAsForm(body, w)
						return
					} else if strings.Split(r.Header["Content-Type"][0], ";")[0] == "application/json" {
						postWithBodyAsJson(body, w)
						return
					} else {
						http.Error(w, "POST data request is unknown requested data", http.StatusBadRequest)
						return
					}
				case "headers":
					if strings.Split(r.Header["Content-Type"][0], ";")[0] == "application/json" {
						postWithBodyAsJsonWithCalc(body, w)
						return
					} else {
						http.Error(w, "POST headers request is unknown headers data", http.StatusBadRequest)
						return
					}
					//fmt.Println(string(body))
					//postWithBody(body, w)
				}
			}
		}
	} else {
		http.Error(w, "No headers set Response expected to have", http.StatusBadRequest)
		return
	}
}

func getWithParam(parameter string, hRes http.ResponseWriter) {
	// handle parameter on error
	hRes.Header().Set("Content-Type", "text/plain")
	hRes.WriteHeader(http.StatusOK)
	fmt.Fprintf(hRes, "Hello, %v", parameter)
}

func getWithBad(hRes http.ResponseWriter) {
	http.Error(hRes, "500 Internal Status Error", http.StatusInternalServerError)
}

func postWithBodyAsData(bodyParam []byte, hRes http.ResponseWriter) {
	//fmt.Println("Parameters Data")
	hRes.WriteHeader(http.StatusOK)
	fmt.Fprintf(hRes, "I got message:\n%v\n", string(bodyParam))
}

func postWithBodyAsForm(bodyParam []byte, hRes http.ResponseWriter) {
	//fmt.Println("Form Data")
	hRes.WriteHeader(http.StatusOK)
	fmt.Fprintf(hRes, "I got message:\n%v\n", string(bodyParam))
}

func postWithBodyAsJson(bodyParam []byte, hRes http.ResponseWriter) {
	fmt.Println("JSON Data")
	hRes.WriteHeader(http.StatusOK)
	fmt.Fprintf(hRes, "I got message:\n%v\n", string(bodyParam))
}

func postWithBodyAsJsonWithCalc(bodyParam []byte, hRes http.ResponseWriter) {
	var bodyParsData CalculationData

	fmt.Println("JSON Data for headers")
	err := json.Unmarshal(bodyParam, &bodyParsData)
	if err != nil {
		fmt.Fprintf(hRes, "Header with wrong data:\n%v, error:\n%v", string(bodyParam), err)
		return
	} else {
		// calculation
		firstHeaderParam, err := strconv.Atoi(bodyParsData.A)
		if err != nil {
			fmt.Fprintf(hRes, "Header with wrong first param:\n%v, error:\n%v", bodyParsData.A, err)
			return
		}
		secondHeaderParam, err := strconv.Atoi(bodyParsData.B)
		if err != nil {
			fmt.Fprintf(hRes, "Header with wrong first param:\n%v, error:\n%v", bodyParsData.B, err)
			return
		}
		fmt.Println(bodyParsData.A)
		fmt.Println(bodyParsData.B)
		tmpSum := strconv.Itoa(firstHeaderParam + secondHeaderParam)
		hRes.Header().Add("a+b", tmpSum)
		hRes.WriteHeader(http.StatusOK)
		fmt.Println(hRes.Header().Values("a+b"))
		fmt.Println(hRes.Header())
		fmt.Fprintf(hRes, "I got message:\n%v\n%v\n", string(bodyParam), hRes.Header())
	}
}

// Start /** Starts the web server listener on given host and port.
func Start(host string, port int) {
	router := mux.NewRouter()

	router.HandleFunc("/name/{param}", handleParam)
	router.HandleFunc("/bad", handleParam)
	router.HandleFunc("/data", handleParam)
	router.HandleFunc("/headers", handleParam)

	log.Println(fmt.Printf("Starting API server on %s:%d\n", host, port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router); err != nil {
		log.Fatal(err)
	}
}

//main /** starts program, gets HOST:PORT param and calls Start func.
func main() {
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8081
	}
	Start(host, port)
}
