package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// this is plain dummy example code only
// not intended to be "good" go code :-)
func dumpHeaders(srcreq *http.Request) {
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
		// debug
		hval := srcreq.Header.Get(header)
		log.Printf("Got header %v (%v)", header, hval)
		// debug
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/", helloHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	go listenAndServe(port)

	select {}
}
