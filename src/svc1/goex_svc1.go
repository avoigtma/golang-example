package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// this is plain dummy example code only
// not intended to be "good" go code :-)
func dumpHeaders(r *http.Request) {
	headers := []string{
		"x-request-id",
		"x-b3-traceid",
		"x-b3-spanid",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"x-b3-flags",
		"x-ot-span-context",
	}
	for _, header := range headers {
		// debug output
		headerval := r.Header.Get(header)
		log.Printf("Got header %v (%v)", header, headerval)
		// debug output
	}
}

func svc1Handler(w http.ResponseWriter, r *http.Request) {
	//	ctx := r.Context()
	dumpHeaders(r)

	response := os.Getenv("RESPONSE")
	if len(response) == 0 {
		response = "Hello OpenShift!"
	}

	fmt.Fprintln(w, response)
	log.Println("Servicing request.")
}

func listenAndServe(port string) {
	log.Printf("serving on %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func main() {
	http.HandleFunc("/", svc1Handler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	go listenAndServe(port)

	select {}
}
