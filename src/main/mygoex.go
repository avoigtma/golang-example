package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinpropagation "github.com/openzipkin/zipkin-go/propagation/b3"
	zipkinreporter "github.com/openzipkin/zipkin-go/reporter"
)

// this is plain dummy example code only
// not intended to be "good" go code :-)

func readURL(client http.Client, ctx context.Context, url string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	// endpoint, err := zipkin.NewEndpoint("main", url)
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
	//req = req.WithContext(zctx)

	zipkinpropagation.InjectHTTP(req)

	//	res, err := client.Do(req)
	res, err := zclient.DoWithAppSpan(req, "main")
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
	// otel instrumentation
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	ctx := context.Background()
	// zipkin
	zipkinpropagation.ExtractHTTP(r)
	// zipkin

	log.Println("Servicing request.")

	responseEnv := os.Getenv("RESPONSE")
	svc1url := os.Getenv("SERVICE1_URL")
	svc2url := os.Getenv("SERVICE2_URL")
	response := "Hello Caller!"
	response += "\n\n\n"
	if len(responseEnv) == 0 {
		log.Println("No response value in env configured")
		response += "No response value in env configured\n"
	} else {
		log.Print("Response value in env: ")
		log.Println(responseEnv)
		response += responseEnv
		response += "'"
		response += "\n\n\n"
	}

	/// service 1
	if len(svc1url) == 0 {
		log.Println("No service-1 url in env configured")
		response += "No service-1 url in env configured\n"
	} else {
		log.Print("Calling service on URL: ")
		log.Println(svc1url)
		response += "Result from Service-1:\n"
		response += readURL(client, ctx, svc1url)
		response += "\n\n\n"
	}

	/// service 2
	if len(svc2url) == 0 {
		log.Println("No service-2 url in env configured")
		response += "No service-2 url in env configured\n"
	} else {
		log.Print("Calling service on URL: ")
		log.Println(svc2url)
		response += "Result from Service-2:\n"
		response += readURL(client, ctx, svc2url)
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
