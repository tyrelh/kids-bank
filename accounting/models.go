package accounting

import "time"

type Transaction struct {
	Id                    int
	CreatedAt             time.Time
	AccountId             int
	RollingBalanceDollars float32
	AmountDollars         float32
	Type                  string
}

type Account struct {
	Id                int
	Name              string
	InterestRate      float32
	InterestFrequency int
}
