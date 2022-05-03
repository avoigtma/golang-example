package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func readURL(client http.Client, inreq *http.Request, ctx context.Context, url string) string {
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

func randomOutput(r *http.Request, url string) string {
	const maxIterations int = 10

	rand.Seed(time.Now().UnixNano())

	iterations := rand.Intn(maxIterations)

	ctx := r.Context()

	var sb strings.Builder

	log.Printf("Creating output %d iterations.", iterations)
	log.Printf("Service URL to call: '%s' ", url)

	client := http.Client{
		Timeout: 60 * time.Second,
	}

	for i := 1; i < iterations; i++ {
		minSleep := 100
		maxSleep := 500
		randSleep := rand.Intn(maxSleep-minSleep+1) + minSleep
		log.Printf("Iteration %d: Sleeping %d milliseconds, then adding next string fragment to output\n", i, randSleep)
		timeSleep := time.Duration(randSleep) * time.Millisecond
		time.Sleep(timeSleep)
		sb.WriteString("Sleeping")
		sb.WriteString(strconv.FormatInt(int64(timeSleep), 10))
		sb.WriteString(" ms\n")

		/// service 3
		if len(url) == 0 {
			log.Println("No service-3 url in env configured")
			sb.WriteString("No service-3 url in env configured\n")
		} else {
			log.Print("Calling service on URL: ")
			log.Println(url)
			sb.WriteString("Result from Service-3, iteration #")
			sb.WriteString(strconv.Itoa(i))
			sb.WriteString(":\n")
			sb.WriteString(readURL(client, r, ctx, url))
			sb.WriteString("\n\n\n")
		}

	}

	result := sb.String()
	log.Println(result)
	return result
}

func svc2Handler(w http.ResponseWriter, r *http.Request) {

	svc3url := os.Getenv("SERVICE3_URL")
	response := ""
	if len(response) == 0 {
		response = randomOutput(r, svc3url)
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
	http.HandleFunc("/", svc2Handler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	go listenAndServe(port)

	select {}
}
