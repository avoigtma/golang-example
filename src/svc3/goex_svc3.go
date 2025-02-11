package main

import (
	"fmt"
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

func randomOutput() string {
	const outputPart string = "this-is-a-dummy-output-line-which-will-be-concatenated"
	const maxIterations int = 20

	rand.Seed(time.Now().UnixNano())

	iterations := rand.Intn(maxIterations) + 1

	var sb strings.Builder

	log.Printf("Creating output %d iterations.", iterations)

	for i := 1; i < iterations; i++ {
		minSleep := 10
		maxSleep := 100
		randSleep := rand.Intn(maxSleep-minSleep+1) + minSleep
		log.Printf("Iteration %d: Sleeping %d seconds, then adding next string fragment to output\n", i, randSleep)
		timeSleep := time.Duration(randSleep) * time.Millisecond
		time.Sleep(timeSleep)

		sb.WriteString("- ")
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(" - (sleep delay: ")
		sb.WriteString(strconv.FormatInt(int64(timeSleep)/1000000, 10))
		sb.WriteString("ms) - ")
		sb.WriteString(outputPart)
		sb.WriteString(" | ")
	}

	result := sb.String()
	log.Println(result)
	return result
}

func svc3Handler(w http.ResponseWriter, r *http.Request) {
	// debug
	dumpHeaders(r)
	// debug

	response := randomOutput()

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
	http.HandleFunc("/", svc3Handler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	go listenAndServe(port)

	select {}
}
