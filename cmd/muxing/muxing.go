package main

import (
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
var (
	noHeadesetResp  string = "No headers set Response expected to have"
	contentTypeName string = "Content-Type"
	gotMessageStr   string = "I got message:"
)

// get request with param create struct param
func handleParam(w http.ResponseWriter, r *http.Request) {
	var paramsSlice []string

	// check kind of request
	// Get with param
	if r.Method == http.MethodGet {
		// check param
		paramsSlice = strings.Split(r.URL.Path, "/")
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
			http.Error(w, noHeadesetResp, http.StatusBadRequest)
			return
		}
	} else if r.Method == http.MethodPost { // Methos POST with body
		paramsSlice = strings.Split(r.URL.Path, "/")
		if len(paramsSlice) < 2 {
			http.Error(w, "expect /data or /header in task handler", http.StatusBadRequest)
			return
		} else {
			// take body into handler
			body, err := ioutil.ReadAll(r.Body)
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
					// check conetnt type in Header
					if _, ok := r.Header[contentTypeName]; ok {
						if strings.Split(r.Header[contentTypeName][0], ";")[0] == "application/x-www-form-urlencoded" {
							postWithBodyAsData(body, w)
							return
						} else if strings.Split(strings.Split(r.Header[contentTypeName][0], ";")[0], " ")[0] == "multipart/form-data" {
							postWithBodyAsForm(body, w)
							return
						} else if strings.Split(r.Header[contentTypeName][0], ";")[0] == "application/json" {
							postWithBodyAsJson(body, w)
							return
						} else {
							http.Error(w, "POST data request is unknown requested data", http.StatusBadRequest)
							return
						}
					} else {
						// Contetn-Type there is not in Header
						postWithBodyAsData(body, w)
					}
				case "headers":
					// take a,b
					val_a, ok_a := r.Header["A"]
					val_b, ok_b := r.Header["B"]
					// there are a,b
					if ok_a && ok_b {
						postWithBodyAsJsonWithCalc(w, val_a[0], val_b[0])
					} else {
						http.Error(w, noHeadesetResp, http.StatusBadRequest)
						return
					}
				}
			}
		}
	} else {
		http.Error(w, noHeadesetResp, http.StatusBadRequest)
		return
	}
}

func getWithParam(parameter string, hRes http.ResponseWriter) {
	// handle parameter on error
	hRes.Header().Set(contentTypeName, "text/plain")
	hRes.WriteHeader(http.StatusOK)
	fmt.Fprintf(hRes, "Hello, %v!", parameter)
}

func getWithBad(hRes http.ResponseWriter) {
	http.Error(hRes, "500 Internal Status Error", http.StatusInternalServerError)
}

func postWithBodyAsData(body_Param []byte, hRes http.ResponseWriter) {
	hRes.WriteHeader(http.StatusOK)
	fmt.Fprintf(hRes, "%v\n%v", gotMessageStr, string(body_Param))
}

func postWithBodyAsForm(body_Param []byte, hRes http.ResponseWriter) {
	hRes.WriteHeader(http.StatusOK)
	fmt.Fprintf(hRes, "%v\n%v", gotMessageStr, string(body_Param))
}

func postWithBodyAsJson(body_Param []byte, hRes http.ResponseWriter) {
	hRes.WriteHeader(http.StatusOK)
	fmt.Fprintf(hRes, "%v\n%v", gotMessageStr, string(body_Param))
}

func postWithBodyAsJsonWithCalc(hRes http.ResponseWriter, first_header_value string, second_header_value string) {
	// calculation
	// check on number a,b
	firstHeaderParam, err := strconv.Atoi(first_header_value)
	if err != nil {
		fmt.Fprintf(hRes, "Header with wrong first param:\n%v, error:\n%v", first_header_value, err)
		return
	}
	secondHeaderParam, err := strconv.Atoi(second_header_value)
	if err != nil {
		fmt.Fprintf(hRes, "Header with wrong second param:\n%v, error:\n%v", second_header_value, err)
		return
	}
	// move sum to stirng
	tmpSum := strconv.Itoa(firstHeaderParam + secondHeaderParam)
	// add new header
	hRes.Header().Add("a+b", fmt.Sprintf("%v", tmpSum))
	hRes.WriteHeader(http.StatusOK)
	// have send messages
	fmt.Fprintf(hRes, "%v\n%v", gotMessageStr, hRes.Header().Get("a+b"))
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
