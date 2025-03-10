package renderer

import (
	"embed"
	"fmt"
	"kids-bank/accounting"
	"net/http"
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

func RenderIndex(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFS(templates, "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderTransactions(w http.ResponseWriter, r *http.Request) {

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

	templ, err := template.New("transactions.html").Funcs(funcMap).ParseFS(templates, "transactions.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
