package validator

import (
	"golang.org/x/exp/constraints"
)

var AllowedCurrency = []string{"rub", "usd", "eur"}

func ValidateExchangeRequest[T constraints.Integer | constraints.Float](v *Validator, fromCurrency, toCurrency string, amount T) {
	v.Check(IsPermittedValue(fromCurrency, AllowedCurrency...), "from_currency", "supported currencies: rub, usd, eur")
	v.Check(IsPermittedValue(toCurrency, AllowedCurrency...), "to_currency", "supported currencies: rub, usd, eur")
	v.Check(amount > 0, "amount", "cannot be less or equal to zero")
	v.Check(fromCurrency != toCurrency, "currency", "source and target currencies cannot be the same")
}
