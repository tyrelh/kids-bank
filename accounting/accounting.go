package accounting

import (
	"database/sql"
	"fmt"
	"kids-bank/database"
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
	query := "SELECT rolling_amount_dollars FROM account WHERE account_type = ? ORDER BY id DESC LIMIT 1"
	db := database.Db()
	row := db.QueryRow(query, account)
	var balance float32
	err := row.Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("error querying balance for account %s: %w", account, err)
	}
	return balance, nil
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
