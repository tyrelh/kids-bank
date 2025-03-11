package accounting

import (
	"database/sql"
	"fmt"
	"kids-bank/database"
	"log"
	"net/http"
	"strconv"
)

func GetAllTransactionsForAccount(account string) ([]Transaction, error) {
	if account == "" {
		return []Transaction{}, fmt.Errorf("account cannot be empty")
	}
	query := "SELECT * FROM account WHERE account_type = ?"
	db := database.Db()
	rows, err := db.Query(query, account)
	if err != nil {
		return []Transaction{}, fmt.Errorf("error querying transactions for account %s: %w", account, err)
	}
	defer rows.Close()
	transactions, err := scanMultipleTransactionRows(rows)
	if err != nil {
		return []Transaction{}, fmt.Errorf("error scanning transactions for account %s: %w", account, err)
	}
	return transactions, nil
}

func GetCurrentBalanceForAccount(account string) (float32, error) {
	if account == "" {
		return 0, fmt.Errorf("account cannot be empty")
	}
	latestTransaction, err := getMostRecentTransactionForAccount(account)
	if err != nil {
		return 0, fmt.Errorf("error getting most recent transaction for account %s: %w", account, err)
	}
	return latestTransaction.RollingAmountDollars, nil
}

func Deposit(w http.ResponseWriter, r *http.Request) {
	account := "savings"
	amountString := r.FormValue("deposit")
	amountFloat64, err := strconv.ParseFloat(amountString, 32)
	if err != nil {
		http.Error(w, "error parsing deposit amount: "+err.Error(), http.StatusBadRequest)
		return
	}
	amountFloat := float32(amountFloat64)

	latestTransaction, err := getMostRecentTransactionForAccount(account)
	if err != nil {
		http.Error(w, "error getting most recent transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	newBalance := latestTransaction.RollingAmountDollars + amountFloat

	log.Println("Creating transaction for $" + amountString + " deposit to " + account)

	query := "INSERT INTO account (rolling_amount_dollars, change_amount_dollars, transaction_type, account_type) VALUES (?, ?, ?, ?)"
	db := database.Db()
	_, err = db.Exec(query, newBalance, amountFloat, "deposit", account)
	if err != nil {
		http.Error(w, "error inserting deposit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Deposit successful, new balance is $" + fmt.Sprintf("%.2f", newBalance))
}

func getMostRecentTransactionForAccount(account string) (Transaction, error) {
	if account == "" {
		return Transaction{}, fmt.Errorf("account cannot be empty")
	}
	query := "SELECT * FROM account WHERE account_type = ? ORDER BY id DESC LIMIT 1"
	db := database.Db()
	row := db.QueryRow(query, account)
	var transaction Transaction
	err := row.Scan(
		&transaction.Id,
		&transaction.CreatedAt,
		&transaction.RollingAmountDollars,
		&transaction.ChangeAmountDollars,
		&transaction.TransactionType,
		&transaction.Account,
	)
	if err != nil {
		return Transaction{}, fmt.Errorf("error querying most recent transaction for account %s: %w", account, err)
	}
	return transaction, nil
}

func scanMultipleTransactionRows(rows *sql.Rows) ([]Transaction, error) {
	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&transaction.Id,
			&transaction.CreatedAt,
			&transaction.RollingAmountDollars,
			&transaction.ChangeAmountDollars,
			&transaction.TransactionType,
			&transaction.Account,
		)
		if err != nil {
			return []Transaction{}, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
