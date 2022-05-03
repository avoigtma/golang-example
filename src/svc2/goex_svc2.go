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

	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinpropagation "github.com/openzipkin/zipkin-go/propagation/b3"
	zipkinreporter "github.com/openzipkin/zipkin-go/reporter"
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
		hval := srcreq.Header.Get(header)clear
		log.Printf("Got header %v (%v) - adding to new request.", header, hval)
		// debug
	}
}

func readURL(client http.Client, inreq *http.Request, ctx context.Context, url string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	// endpoint, err := zipkin.NewEndpoint("svc-2", url)
	// if err != nil {
	// 	log.Fatalf("unable to create endpoint: %+v\n", err)
	// }
	// reporter := logreporter.NewReporter(log.New(os.Stderr, "", log.LstdFlags))
	// defer reporter.Close()
	reporter := zipkinreporter.NewNoopReporter()
	defer reporter.Close()
	// tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	tracer, err := zipkin.NewTracer(reporter)
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}
	span := zipkin.SpanFromContext(req.Context())
	//span.Tag("custom_key", "some value")

	//var zclient *zipkinhttp.Client
	zclient, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	zctx := zipkin.NewContext(req.Context(), span)
	req = req.WithContext(zctx)

	zipkinpropagation.InjectHTTP(req)

	propagateHeaders(inreq, req)
	//&res, err := client.Do(req)
	res, err := zclient.DoWithAppSpan(req, "svc-2")
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

	span.Finish()
	return strBody

}

func randomOutput(r *http.Request, ctx context.Context, url string) string {
	const maxIterations int = 10

	rand.Seed(time.Now().UnixNano())

	iterations := rand.Intn(maxIterations)

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
		log.Printf("Iteration %d: Sleeping %d seconds, then adding next string fragment to output\n", i, randSleep)
		time.Sleep(time.Duration(randSleep) * time.Millisecond)

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

func helloHandler(w http.ResponseWriter, r *http.Request) {
	zipkinpropagation.ExtractHTTP(r)
	ctx := r.Context()

	svc3url := os.Getenv("SERVICE3_URL")
	response := ""
	if len(response) == 0 {
		response = randomOutput(r, ctx, svc3url)
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
