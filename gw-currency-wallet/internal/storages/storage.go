package storages

type UserModelInterface interface {
	CreateUser(user *User) error
	FindUser(user *User) error
}

type WalletModelInterface interface {
	ChangeBalance(id uint, amount int, currency string) (Wallet, error)
	GetBalance(id uint) (Wallet, error)
}
