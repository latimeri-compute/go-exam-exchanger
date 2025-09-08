package swagger

import "github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"

// uhhhhhhhh don't worry about it

type ErrorResponse struct {
	Err string `json:"error"`
}

type ErrorInvalidUserPassword struct {
	Err string `json:"error" example:"Invalid username or password"`
}

type ErrUserEmailExists struct {
	Err string `json:"error" example:"Username or email already exists"`
}

type ErrorUnauthorizedResponse struct {
	Err string `json:"error" example:"Unauthorized"`
}

type ErrorInvalidCurrencies struct {
	Err string `json:"error" example:"Insufficient funds or invalid currencies"`
}

type ErrorInsufficientFunds struct {
	Err string `json:"error" example:"Insufficient funds or invalid amount"`
}

type ErrInvalidCurrencyAmount struct {
	Err string `json:"error" example:"Invalid amount or currency"`
}

type InternalErrorResponse struct {
	Err string `json:"error" example:"Internal error"`
}

type MessageResponse struct {
	Mes string `json:"message"`
}

type ExampleFailedExchange struct {
	Err string `json:"error" example:"Failed to retrieve exchange rates"`
}

type ExampleExchangeRates struct {
	Rates ExchangeRates `json:"rates"`
}

type ExchangeRates struct {
	RU utils.Currency `json:"rub->usd" example:"0.35"`
	UR utils.Currency `json:"usd->rub" example:"1.35"`
	UE utils.Currency `json:"usd->eur" example:"2.35"`
	EU utils.Currency `json:"eur->usd" example:"4.35"`
	RE utils.Currency `json:"rub->eur" example:"5.35"`
	ER utils.Currency `json:"eur->rub" example:"6.35"`
}

type ReturnToken struct {
	Token string `json:"token" example:"JWT_TOKEN"`
}

type ExampleUserCreated struct {
	Mes string `json:"message" example:"User registered successfully"`
}

type ReturnBalance struct {
	Balance balance `json:"balance"`
}

type ExampleDeposit struct {
	Message string `json:"message" example:"Account topped up successfully"`
}
type ExampleWithdraw struct {
	Message string `json:"message" example:"Account topped up successfully"`
}

type balance struct {
	USD utils.Currency `json:"USD" example:"120.00"`
	EUR utils.Currency `json:"EUR" example:"10.00"`
	RUB utils.Currency `json:"RUB" example:"45.50"`
}
