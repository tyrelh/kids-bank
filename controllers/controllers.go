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
	transactions, err := accounting.GetAllTransactionsForAccount(accounting.SAVINGS_ACCOUNT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	balance, err := accounting.GetCurrentBalanceForAccount(accounting.SAVINGS_ACCOUNT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	account, err := accounting.GetAccountByName(accounting.SAVINGS_ACCOUNT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	interestAlreadyApplied, err := accounting.HasInterestBeenAppliedInPeriod(accounting.SAVINGS_ACCOUNT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Transactions    []accounting.Transaction
		Balance         float32
		Account         accounting.Account
		InterestApplied bool
	}{
		Transactions:    transactions,
		Balance:         balance,
		Account:         account,
		InterestApplied: interestAlreadyApplied,
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
	transaction, err := accounting.Deposit(amountFloat, accounting.SAVINGS_ACCOUNT)
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
	transaction, err := accounting.ApplyInterest(accounting.SAVINGS_ACCOUNT)
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

func UpdateInterestRate(w http.ResponseWriter, r *http.Request) {
	rateString := r.FormValue("interest-rate")
	rateFloat64, err := strconv.ParseFloat(rateString, 32)
	if err != nil {
		http.Error(w, "error parsing rate: "+err.Error(), http.StatusBadRequest)
		return
	}
	rateFloat := float32(rateFloat64)
	account, err := accounting.GetAccountByName(accounting.SAVINGS_ACCOUNT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	account.InterestRate = rateFloat
	err = accounting.UpdateAccount(account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return balance
	_, err = w.Write([]byte(fmt.Sprintf("%.2f", rateFloat)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
