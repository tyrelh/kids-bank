package accounting

import "time"

type Transaction struct {
	Id                   int
	CreatedAt            time.Time
	RollingAmountDollars float32
	ChangeAmountDollars  float32
	TransactionType      string
	Account              string
}
