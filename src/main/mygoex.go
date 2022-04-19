package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// this is plain dummy example code only
// not intended to be "good" go code :-)

func readURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	// read the response body on the line below
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// convert the body to string and log
	strBody := string(body)
	log.Println(strBody)
	return strBody
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	responseEnv := os.Getenv("RESPONSE")
	svc1url := os.Getenv("SERVICE1_URL")
	svc2url := os.Getenv("SERVICE2_URL")
	response := "Hello Caller!"
	response += "\n\n\n"
	if len(response) == 0 {
		log.Println("No response value in env configured")
		response += "No response value in env configured"
	} else {
		log.Print("Response value in env: ")
		log.Println(responseEnv)
		response += responseEnv
		response += "'"
		response += readURL("https://www.google.de")
	}
	response += "\n\n\n"
	if len(svc1url) == 0 {
		log.Println("No service-1 url in env configured")
	} else {
		log.Print("Calling service on URL: ")
		log.Println(svc1url)
		response += "Result from Service-1:\n"
		response += readURL(svc1url)
	}
	response += "\n\n\n"
	if len(svc2url) == 0 {
		log.Println("No service-2 url in env configured")
	} else {
		log.print("Calling service on URL: ")
		log.Println(svc2url)
		response += "Result from Service-2:\n"
		response += readURL(svc2url)
	}

	fmt.Fprintln(w, response)
	log.Println("Servicing request.")
}

func listenAndServe(port string) {
	fmt.Printf("serving on %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func main() {
	http.HandleFunc("/", helloHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	go listenAndServe(port)

	port = os.Getenv("SECOND_PORT")
	if len(port) == 0 {
		port = "8888"
	}
	go listenAndServe(port)

	select {}
}
