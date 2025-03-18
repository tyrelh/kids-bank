package accounting

import (
	"fmt"
	"kids-bank/database"
)

func GetAccountByName(name string) (Account, error) {
	account, err := getAccountByName(name)
	if err != nil {
		return Account{}, fmt.Errorf("error getting account by name %s: %w", name, err)
	}
	return account, nil
}

func UpdateAccount(account Account) error {
	query := "UPDATE accounts SET interest_rate = ?, interest_frequency = ? WHERE id = ?"
	db := database.Db()
	_, err := db.Exec(query, account.InterestRate, account.InterestFrequency, account.Id)
	if err != nil {
		return fmt.Errorf("error updating account: %w", err)
	}
	return nil
}

func getAccountByName(name string) (Account, error) {
	if name == "" {
		return Account{}, fmt.Errorf("name cannot be empty")
	}
	query := "SELECT * FROM accounts WHERE name = ?"
	db := database.Db()
	row := db.QueryRow(query, name)
	var account Account
	err := row.Scan(
		&account.Id,
		&account.Name,
		&account.InterestRate,
		&account.InterestFrequency,
	)
	if err != nil {
		return Account{}, fmt.Errorf("error querying account by name %s: %w", name, err)
	}
	return account, nil
}
