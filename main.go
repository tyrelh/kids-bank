package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
)

func lambdaHandler() {
	// Example of using Go to return HTML from a Lambda function
	// https://stackoverflow.com/questions/76430232/how-to-return-html-from-a-go-lambda-function
	log.Println("Hello from Lambda!")
}

func server() {
	log.Println("Hello from Server!")
	listenerPort := os.Getenv("KB_PORT")
	if listenerPort == "" {
		listenerPort = "8080"
	}

	handleTest := func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("index.html"))
		err := templ.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
