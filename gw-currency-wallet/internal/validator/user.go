package validator

func ValidateUser(v *Validator, email, password string) bool {
	v.Check(Matches(email, EmailRX), "email", "is not a valid email")
	v.Check(len(email) <= 255, "email", "email is too long")

	v.Check(len(password) >= 6, "password", "must be at least 6 characters long")

	return v.Valid()
}
