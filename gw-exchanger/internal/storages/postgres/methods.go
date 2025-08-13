package postgres

func (db *DB) GetAll() {
	// var res []storages.ReturnExchanges
	// var rate storages.ExchangeRate
	// var valute storages.Valute

	// ctx := context.Background()

	tx := db.DB.Select("valutes.code", "").Joins("FROM exchange_rates AS r INNER JOIN valutes AS v USING(v.id)")
	if err := tx.Error; err != nil {
		return
	}

}
