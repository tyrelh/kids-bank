package accounting

type Transaction struct {
	Id                   int
	CreatedAt            string
	RollingAmountDollars float32
	ChangeAmountDollars  float32
	TransactionType      string
	Account              string
}
