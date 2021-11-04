package models

import "time"

// models collections
var (
	UserCollectionName = "users"
)

type User struct {
	ID                 string              `bson:"_id,omitempty" json:"_id,omitempty"`
	FullName           string              `bson:"full_name" validate:"required,min=2,max=100" json:"full_name"`
	Email              string              `bson:"email" validate:"email,required" json:"email"`
	Phone              string              `bson:"phone" json:"phone"`
	Password           string              `bson:"password" json:"password" validate:"required,min=6"`
	CreatedAt          time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time           `bson:"updated_at" json:"updated_at"`
	IsVerified         bool                `bson:"isverified" json:"isverified"`
	PortfolioPositions []PortfolioPosition `bson:"portfolio_positions" json:"portfolio_positions"`
}

type PortfolioPosition struct {
	Symbol        string  `bson:"symbol" json:"symbol"`
	TotalQuantity float64 `bson:"total_quantity" json:"total_quantity"`
	EquityValue   float64 `bson:"equity_value" json:"equity_value"`
	PricePerShare float64 `bson:"price_per_share" json:"price_per_share"`
}

type AuthCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserUpdate struct {
	FullName string `bson:"full_name" validate:"required,min=2,max=100" json:"full_name"`
	Phone    string `bson:"phone" validate:"required" json:"phone"`
}

type PasswordUpdate struct {
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

// custom positions data
var position1 = PortfolioPosition{Symbol: "AAPL", TotalQuantity: 20, EquityValue: 2500.0, PricePerShare: 125.0}
var position2 = PortfolioPosition{Symbol: "TSLA", TotalQuantity: 5.0, EquityValue: 3000.0, PricePerShare: 600.0}
var position3 = PortfolioPosition{Symbol: "AMZN", TotalQuantity: 30, EquityValue: 4500.0, PricePerShare: 150.0}
var Positions = []PortfolioPosition{position1, position2, position3}
