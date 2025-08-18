package storages

// ИНТЕРФЕЙС
type ExchangerModelInterface interface {
	GetAll() ([]Exchange, error)
	GetRateBetween()
}
