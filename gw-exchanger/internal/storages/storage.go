package storages

// ИНТЕРФЕЙС
type ExchangerModelInterface interface {
	GetAll() ([]Exchange, error)
	GetRateBetween(fromValute, toValute string) (Exchange, error)
}
