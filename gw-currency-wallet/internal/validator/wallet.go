package validator

import "strings"

func (v *Validator) CheckBalanceChange(amount float64, currency string) {
	v.Check(IsPermittedValue(strings.ToLower(currency), "usd", "eur", "rub"), "currency", "supported currencies: rub, usd, eur")
	v.Check(amount > 0, "amount", "can't be zero or less")
}
