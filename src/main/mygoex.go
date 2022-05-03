package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// this is plain dummy example code only
// not intended to be "good" go code :-)

func propagateHeaders(srcreq *http.Request, dstreq *http.Request) {
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
		dstreq.Header.Add(header, srcreq.Header.Get(header))
		// debug
		hval := srcreq.Header.Get(header)
		log.Printf("Got header %v (%v) - adding to outbound request.", header, hval)
		// debug
	}
}

func readURL(client http.Client, inreq *http.Request, url string) string {
	ctx := context.Background()

	outreq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	propagateHeaders(inreq, outreq)

	res, err := client.Do(outreq)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	strBody := string(body)
	log.Println(strBody)
	return strBody

}

func MainServiceHandler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	log.Println("Servicing request.")

	svc1url := os.Getenv("SERVICE1_URL")
	svc2url := os.Getenv("SERVICE2_URL")

	response := "Hello Caller!"
	response += "\n\n\n"

	/// service 1
	if len(svc1url) == 0 {
		log.Println("No service-1 url in env configured")
		response += "No service-1 url in env configured\n"
	} else {
		log.Print("Calling service on URL: ")
		log.Println(svc1url)
		response += "\n---------------------\n"
		response += "Result from Service-1:\n"
		response += readURL(client, r, svc1url)
		response += "\n---------------------\n"
		response += "\n\n\n"
	}

	/// service 2
	if len(svc2url) == 0 {
		log.Println("No service-2 url in env configured")
		response += "No service-2 url in env configured\n"
	} else {
		log.Print("Calling service on URL: ")
		log.Println(svc2url)
		response += "\n---------------------\n"
		response += "Result from Service-2:\n"
		response += readURL(client, r, svc2url)
		response += "\n---------------------\n"
	}

	fmt.Fprintln(w, response)

}

func listenAndServe(port string) {
	log.Printf("serving on %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func main() {
	http.HandleFunc("/", MainServiceHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	go listenAndServe(port)

	select {}
}
