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

type Rate struct {
	Id           int
	Rate         float32
	RateType     string
	Frequency    int
	PreviousRate float32
	DateModified time.Time
}
