package postgres

import "github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages"

func (db *DB) GetAll() ([]storages.Exchange, error) {
	var results []storages.Exchange

	err := db.DB.Model(&storages.Exchange{}).
		Select("from_valute.code", "to_valute.code", "exchange.rate").
		Joins("FromValute").Joins("ToValute").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (db *DB) GetRateBetween(fromValute, toValute string) (storages.Exchange, error) {
	var results storages.Exchange

	err := db.DB.Model(&storages.Exchange{}).
		Select("from_valute.code", "to_valute.code", "exchange.rate").
		Joins("FromValute").Joins("ToValute").
		Where("from_valute.code = ? AND to_valute.code = ?", fromValute, toValute).
		Find(&results).Error

	if err != nil {
		return storages.Exchange{}, err
	}

	return results, nil
}
