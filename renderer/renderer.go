package renderer

import (
	"embed"
	"kids-bank/accounting"
	"net/http"
	"text/template"
)

//go:embed *.html
var templates embed.FS

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
	data := struct {
		Transactions []accounting.Transaction
	}{
		Transactions: transactions,
	}

	templ, err := template.ParseFS(templates, "transactions.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
