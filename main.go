package main

import (
	"log"
	"net/http"
	"os"

	"kids-bank/renderer"
)

func lambdaHandler() {
	// Example of using Go to return HTML from a Lambda function
	// https://stackoverflow.com/questions/76430232/how-to-return-html-from-a-go-lambda-function
	log.Println("Running in Lambda mode")
}

func server() {
	log.Println("Running in server mode")
	listenerPort := os.Getenv("KB_PORT")
	if listenerPort == "" {
		listenerPort = "8080"
	}

	handleTest := func(w http.ResponseWriter, r *http.Request) {
		renderer.RenderIndex(w, r)
	}

	http.HandleFunc("/", handleTest)
	log.Println("Listening on port " + listenerPort + "...")
	log.Fatal(http.ListenAndServe(":"+listenerPort, nil))
}

func main() {
	shouldRunServer := os.Getenv("KB_SERVER")
	if shouldRunServer != "" {
		server()
	} else {
		lambdaHandler()
	}
}
