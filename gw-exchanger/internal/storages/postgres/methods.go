package postgres

import "github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages"

func (db *DB) GetAll() ([]storages.ReturnExchanges, error) {
	var results []storages.ReturnExchanges

	err := db.DB.Model(&storages.Exchange{}).
		Select("from_valute.code as from_valute_code, to_valute.code as to_valute_code, exchanges.rate").
		Joins("JOIN valutes as from_valute ON exchanges.from_valute_id = from_valute.id").
		Joins("JOIN valutes as to_valute ON exchanges.to_valute_id = to_valute.id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (db *DB) GetRateBetween(fromValute, toValute string) (storages.ReturnExchanges, error) {
	var results storages.ReturnExchanges

	err := db.DB.Model(&storages.Exchange{}).
		Select("from_valute.code as from_valute_code, to_valute.code as to_valute_code, exchanges.rate").
		Joins("JOIN valutes as from_valute ON exchanges.from_valute_id = from_valute.id").
		Joins("JOIN valutes as to_valute ON exchanges.to_valute_id = to_valute.id").
		Where("from_valute.code = ? AND to_valute.code = ?", fromValute, toValute).
		First(&results).Error

	if err != nil {
		return storages.ReturnExchanges{}, err
	}

	return results, nil
}
