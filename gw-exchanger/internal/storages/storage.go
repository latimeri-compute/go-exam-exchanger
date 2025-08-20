package storages

// ИНТЕРФЕЙС
type ExchangerModelInterface interface {
	GetAll() ([]ReturnExchanges, error)
	GetRateBetween(fromValute, toValute string) (ReturnExchanges, error)
}
