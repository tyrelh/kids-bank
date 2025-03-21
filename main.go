package main

import (
	"log"
	"net/http"
	"os"

	"kids-bank/controllers"
	"kids-bank/database"
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

	http.HandleFunc("/admin", controllers.RenderAdmin)
	http.HandleFunc("/deposit", controllers.Deposit)
	http.HandleFunc("/applyInterest", controllers.ApplyInterest)
	http.HandleFunc("/updateInterestRate", controllers.UpdateInterestRate)
	http.HandleFunc("/updateInterestFrequency", controllers.UpdateInterestFrequency)

	log.Printf("Listening on http://localhost:%s\n", listenerPort)
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
