package accounting

import (
	"database/sql"
	"fmt"
	"kids-bank/database"
	"log"
)

var (
	SAVINGS_ACCOUNT        = "savings"
	INTEREST_TRANSACTION   = "interest"
	DEPOSIT_TRANSACTION    = "deposit"
	WITHDRAWAL_TRANSACTION = "withdrawal"
)

func GetAllTransactionsForAccount(accountName string) ([]Transaction, error) {
	if accountName == "" {
		return []Transaction{}, fmt.Errorf("account cannot be empty")
	}
	account, err := getAccountByName(accountName)
	if err != nil {
		return []Transaction{}, fmt.Errorf("error getting account by name %s: %w", accountName, err)
	}
	query := "SELECT * FROM transactions WHERE account_id = ?"
	db := database.Db()
	rows, err := db.Query(query, account.Id)
	if err != nil {
		return []Transaction{}, fmt.Errorf("error querying transactions for account %s: %w", account.Name, err)
	}
	defer rows.Close()
	transactions, err := scanMultipleTransactionRows(rows)
	if err != nil {
		return []Transaction{}, fmt.Errorf("GetAllTransactionsForAccount error scanning transactions for account %s: %w", account.Name, err)
	}
	return transactions, nil
}

func GetCurrentBalanceForAccount(accountName string) (float32, error) {
	if accountName == "" {
		return 0, fmt.Errorf("account cannot be empty")
	}
	account, err := getAccountByName(accountName)
	if err != nil {
		return 0, fmt.Errorf("error getting account by name %s: %w", accountName, err)
	}
	balance, err := getCurrentBalanceForAccount(account.Id)
	if err != nil {
		return 0, fmt.Errorf("error getting most recent transaction for account %s: %w", account.Name, err)
	}
	return balance, nil
}

func Deposit(amount float32, accountName string) (Transaction, error) {
	account, err := getAccountByName(accountName)
	if err != nil {
		return Transaction{}, fmt.Errorf("error getting account by name %s: %w", accountName, err)
	}
	currentBalance, err := getCurrentBalanceForAccount(account.Id)
	if err != nil {
		return Transaction{}, fmt.Errorf("error getting current balance for account %s: %w", accountName, err)
	}
	newBalance := currentBalance + amount
	newBalance = RoundFloatToTwoDecimalPlaces(newBalance)
	log.Printf("Creating transaction for %f deposit to %s account", amount, accountName)
	transaction, err := createTransaction(account.Id, newBalance, amount, DEPOSIT_TRANSACTION)
	if err != nil {
		return Transaction{}, fmt.Errorf("error creating transaction: %w", err)
	}
	log.Printf("New balance is %f", transaction.RollingBalanceDollars)
	return transaction, nil
}

///////////////////////////////////////////////////////////////////////////////
// Private functions

func createTransaction(accountId int, newBalance float32, transactionAmount float32, transactionType string) (Transaction, error) {
	query := "INSERT INTO transactions (account_id, rolling_balance_dollars, amount_dollars, type) VALUES (?, ?, ?, ?)"
	db := database.Db()
	result, err := db.Exec(query, accountId, newBalance, transactionAmount, transactionType)
	if err != nil {
		return Transaction{}, fmt.Errorf("error inserting transaction: %w", err)
	}
	transactionId, err := result.LastInsertId()
	if err != nil {
		return Transaction{}, fmt.Errorf("error getting last insert id: %w", err)
	}
	transaction, err := getTransactionById(int(transactionId))
	if err != nil {
		return Transaction{}, fmt.Errorf("error getting transaction by id %d: %w", transactionId, err)
	}
	return transaction, nil
}

func getCurrentBalanceForAccount(accountId int) (float32, error) {
	if accountId == 0 {
		return 0, fmt.Errorf("account cannot be 0")
	}
	query := "SELECT rolling_balance_dollars FROM transactions WHERE account_id = ? ORDER BY id DESC LIMIT 1"
	db := database.Db()
	row := db.QueryRow(query, accountId)
	var balance float32
	err := row.Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("error querying current balance for account %d: %w", accountId, err)
	}
	return balance, nil
}

func getTransactionById(transactionId int) (Transaction, error) {
	if transactionId == 0 {
		return Transaction{}, fmt.Errorf("transaction cannot be 0")
	}
	query := "SELECT * FROM transactions WHERE id = ? LIMIT 1"
	db := database.Db()
	row := db.QueryRow(query, transactionId)
	transaction, err := scanSingleTransactionRow(row)
	if err != nil {
		return Transaction{}, fmt.Errorf("getTransactionById error scanning transaction by id %d: %w", transactionId, err)
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
			&transaction.AccountId,
			&transaction.RollingBalanceDollars,
			&transaction.AmountDollars,
			&transaction.Type,
		)
		if err != nil {
			return []Transaction{}, fmt.Errorf("scanMultipleTransactionRows error scanning transaction row: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func scanSingleTransactionRow(row *sql.Row) (Transaction, error) {
	var transaction Transaction
	err := row.Scan(
		&transaction.Id,
		&transaction.CreatedAt,
		&transaction.AccountId,
		&transaction.RollingBalanceDollars,
		&transaction.AmountDollars,
		&transaction.Type,
	)
	if err != nil {
		return Transaction{}, err
	}
	return transaction, nil
}

func RoundFloatToTwoDecimalPlaces(input float32) float32 {
	return float32(int(input*100)) / 100
}
