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

// LocalCreatedAt returns the CreatedAt time converted to the local timezone
func (t Transaction) LocalCreatedAt() time.Time {
	return t.CreatedAt.In(time.Local)
}

// FormattedLocalTime returns the CreatedAt time as a formatted string in local timezone
func (t Transaction) FormattedLocalTime(format string) string {
	return t.CreatedAt.In(time.Local).Format(format)
}

type Account struct {
	Id                int
	Name              string
	InterestRate      float32
	InterestFrequency int
}
