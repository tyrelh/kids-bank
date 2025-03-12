package main

import (
	"log"
	"net/http"
	"os"

	"kids-bank/accounting"
	"kids-bank/database"
	"kids-bank/renderer"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func lambdaHandler() {
	// Example of using Go to return HTML from a Lambda function
	// https://stackoverflow.com/questions/76430232/how-to-return-html-from-a-go-lambda-function
	log.Println("Running in Lambda mode")
}

func server() {
	log.Println("Running in server mode")

	db := database.Db()
	defer db.Close()

	listenerPort := os.Getenv("KB_PORT")
	if listenerPort == "" {
		listenerPort = "8080"
	}

	http.HandleFunc("/", renderer.RenderIndex)
	http.HandleFunc("/admin", renderer.RenderAdmin)
	http.HandleFunc("/deposit", accounting.Deposit)
	http.HandleFunc("/apply-interest", accounting.ApplyInterest)

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
