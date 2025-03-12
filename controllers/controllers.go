package controllers

import (
	"embed"
	"fmt"
	"kids-bank/accounting"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

//go:embed *.html
var templates embed.FS

// Create a function map with a formatting function
var funcMap = template.FuncMap{
	"formatMoney": func(amount float32) string {
		return fmt.Sprintf("%.2f", amount)
	},
}

func RenderAdmin(w http.ResponseWriter, r *http.Request) {
	transactions, err := accounting.GetAllTransactionsForAccount("savings")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	balance, err := accounting.GetCurrentBalanceForAccount("savings")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Transactions []accounting.Transaction
		Balance      float32
	}{
		Transactions: transactions,
		Balance:      balance,
	}

	templ, err := template.New("admin.html").Funcs(funcMap).ParseFS(templates, "admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Deposit(w http.ResponseWriter, r *http.Request) {
	amountString := r.FormValue("deposit")
	amountFloat64, err := strconv.ParseFloat(amountString, 32)
	if err != nil {
		http.Error(w, "error parsing deposit amount: "+err.Error(), http.StatusBadRequest)
		return
	}
	amountFloat := float32(amountFloat64)
	amountFloat = accounting.RoundFloatToTwoDecimalPlaces(amountFloat)
	transaction, err := accounting.Deposit(amountFloat, accounting.SavingsAccount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return balance
	_, err = w.Write([]byte(fmt.Sprintf("$%.2f", transaction.RollingBalanceDollars)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ApplyInterest(w http.ResponseWriter, r *http.Request) {
	log.Println("Applying interest")
	transaction, err := accounting.ApplyInterest(accounting.SavingsAccount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return balance
	_, err = w.Write([]byte(fmt.Sprintf("$%.2f", transaction.RollingBalanceDollars)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
