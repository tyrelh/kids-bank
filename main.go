package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
)

func lambdaHandler() {
	log.Println("Hello from Lambda!")
}

func server() {
	log.Println("Hello from Server!")
	listenerPort := os.Getenv("KB_PORT")
	if listenerPort == "" {
		listenerPort = "8080"
	}

	handleTest := func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("templates/index.html"))
		templ.Execute(w, nil)
	}
	http.HandleFunc("/", handleTest)
	log.Println("Listening on port " + listenerPort + "...")
	http.ListenAndServe(":"+listenerPort, nil)
}

func main() {
	shouldRunServer := os.Getenv("KB_SERVER")
	if shouldRunServer != "" {
		server()
	} else {
		lambdaHandler()
	}
}
