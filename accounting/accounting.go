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

func ApplyInterest(w http.ResponseWriter, r *http.Request) {
	log.Println("Applying interest to savings account")
	rateType := "interest"
	rate, err := getInterestRateByType(rateType)
	if err != nil {
		http.Error(w, "error getting interest rate: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Applying interest rate of " + fmt.Sprintf("%.2f", rate.Rate) + " to savings account")

	currentBalance, err := GetCurrentBalanceForAccount("savings")
	if err != nil {
		http.Error(w, "error getting current balance: "+err.Error(), http.StatusInternalServerError)
		return
	}
	interestAmount := currentBalance * rate.Rate
	interestAmount = roundFloatToTwoDecimalPlaces(interestAmount)
	newBalance := currentBalance + interestAmount
	newBalance = roundFloatToTwoDecimalPlaces(newBalance)

	log.Println("Applying interest of $" + fmt.Sprintf("%.2f", interestAmount) + " to savings account")
	err = createTransaction("savings", interestAmount, newBalance, "interest")
	log.Println("New balance is $" + fmt.Sprintf("%.2f", newBalance))
	if err != nil {
		http.Error(w, "error creating transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

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
	amountFloat = roundFloatToTwoDecimalPlaces(amountFloat)

	latestTransaction, err := getMostRecentTransactionForAccount(account)
	if err != nil {
		http.Error(w, "error getting most recent transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	newBalance := latestTransaction.RollingAmountDollars + amountFloat
	newBalance = roundFloatToTwoDecimalPlaces(newBalance)

	log.Println("Creating transaction for $" + amountString + " deposit to " + account)

	err = createTransaction(account, amountFloat, newBalance, "deposit")
	if err != nil {
		http.Error(w, "error creating transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Deposit successful, new balance is $" + fmt.Sprintf("%.2f", newBalance))
}

func createTransaction(account string, changeAmount float32, newBalance float32, transactionType string) error {

	query := "INSERT INTO account (rolling_amount_dollars, change_amount_dollars, transaction_type, account_type) VALUES (?, ?, ?, ?)"
	db := database.Db()
	_, err := db.Exec(query, newBalance, changeAmount, transactionType, account)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}
	return nil
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

func getInterestRateByType(rateType string) (Rate, error) {
	if rateType == "" {
		return Rate{}, fmt.Errorf("rate type cannot be empty")
	}
	query := "SELECT * FROM rates WHERE rate_type = ?"
	db := database.Db()
	row := db.QueryRow(query, rateType)
	var rate Rate
	err := row.Scan(
		&rate.Id,
		&rate.Rate,
		&rate.RateType,
		&rate.Frequency,
		&rate.PreviousRate,
		&rate.DateModified,
	)
	if err != nil {
		return Rate{}, fmt.Errorf("error querying rate for type %s: %w", rateType, err)
	}
	return rate, nil
}

func roundFloatToTwoDecimalPlaces(input float32) float32 {
	return float32(int(input*100)) / 100
}
