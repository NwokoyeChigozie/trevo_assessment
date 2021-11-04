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

type FlutterRequestBody struct {
	TxRef          string         `bson:"tx_ref"  json:"tx_ref"`
	Amount         string         `bson:"amount"  json:"amount"`
	Currency       string         `bson:"currency"  json:"currency"`
	RedirectUrl    string         `bson:"redirect_url"  json:"redirect_url"`
	PaymentOptions string         `bson:"payment_options"  json:"payment_options"`
	Customer       Customer       `bson:"customer"  json:"customer"`
	Customizations Customizations `bson:"customizations"  json:"customizations"`
}

type Customer struct {
	Email string `bson:"email"  json:"email"`
	Name  string `bson:"name"  json:"name"`
}

type Customizations struct {
	Title       string `bson:"title"  json:"title"`
	Description string `bson:"description"  json:"description"`
}
