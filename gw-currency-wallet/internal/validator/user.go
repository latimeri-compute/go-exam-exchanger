package validator

func (v *Validator) ValidateUserLogin(username, password string) bool {
	v.Check(len(username) > 0, "username", "cannot be empty")

	return v.Valid()
}

func (v *Validator) ValidateUserRegister(email, username, password string) bool {
	v.Check(Matches(email, EmailRX), "email", "is not a valid email")
	v.Check(len(email) <= 255, "email", "email is too long")
	v.Check(len(username) > 0, "username", "cannot be empty")
	v.Check(len(username) <= 255, "username", "username is too long")
	v.Check(len(password) >= 6, "password", "must be at least 6 characters long")

	return v.Valid()
}
