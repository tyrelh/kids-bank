package accounting

import (
	"database/sql"
	"fmt"
	"kids-bank/database"
	"log"
	"time"
)

func ApplyInterest(accountName string) (Transaction, error) {
	log.Printf("Applying interest to %s account", accountName)
	account, err := getAccountByName(accountName)
	if err != nil {
		return Transaction{}, fmt.Errorf("error getting account by name %s: %w", accountName, err)
	}
	log.Println("Applying interest rate of " + fmt.Sprintf("%.2f", account.InterestRate) + " to " + accountName + " account")

	currentBalance, err := GetCurrentBalanceForAccount(accountName)
	if err != nil {
		return Transaction{}, fmt.Errorf("error getting current balance for account %s: %w", accountName, err)
	}
	interestAmount := currentBalance * account.InterestRate
	interestAmount = RoundFloatToTwoDecimalPlaces(interestAmount)
	newBalance := currentBalance + interestAmount
	newBalance = RoundFloatToTwoDecimalPlaces(newBalance)

	log.Println("Applying interest of $" + fmt.Sprintf("%.2f", interestAmount) + " to " + accountName + " account")
	transaction, err := createTransaction(account.Id, newBalance, interestAmount, INTEREST_TRANSACTION)

	if err != nil {
		return Transaction{}, fmt.Errorf("error creating transaction: %w", err)
	}
	log.Println("New balance is $" + fmt.Sprintf("%.2f", transaction.RollingBalanceDollars))
	return transaction, nil
}

func HasInterestBeenAppliedInPeriod(accountName string) (bool, error) {
	account, err := getAccountByName(accountName)
	if err != nil {
		return false, fmt.Errorf("error getting account by name %s: %w", accountName, err)
	}
	query := "SELECT * FROM transactions WHERE account_id = ? AND type = ? ORDER BY created_at DESC LIMIT 1"
	db := database.Db()
	row := db.QueryRow(query, account.Id, INTEREST_TRANSACTION)
	transactions, err := scanSingleTransactionRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("HasInterestBeenAppliedInPeriod error scanning transactions for account %s: %w", accountName, err)
	}

	if transactions.LocalCreatedAt().AddDate(0, 0, account.InterestFrequency).After(time.Now()) {
		return true, nil
	}
	return false, nil
}
