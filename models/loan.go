package models

var (
	LoanCollection = "loans"
)

type TakeLoanRequest struct {
	Amount   float64 `json:"amount" validate:"required"`
	Duration int     `json:"duration" validate:"duration"`
}
type PayBackLoanRequest struct {
	Amount float64 `json:"amount" validate:"required"`
}

type Loan struct {
	ID          string  `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID      string  `bson:"user_id" json:"user_id"`
	TotalAmount float64 `bson:"total_amount"  json:"total_amount"`
	Balance     float64 `bson:"balance"  json:"balance"`
	Duration    int     `bson:"duration"  json:"duration"`
	TimeTaken   int     `bson:"time_taken"  json:"time_taken"`
	LastPayment int     `bson:"last_payment"  json:"last_payment"`
}
