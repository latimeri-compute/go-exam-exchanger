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

func (db *DB) GetRateBetween() {
	// TODO implement
}

/* SELECT
    fv.code AS from_valute_code,
    tv.code AS to_valute_code,
    e.rate
FROM
    exchange e
JOIN
    valutes fv ON e.from_valute_id = fv.id
JOIN
    valutes tv ON e.to_valute_id = tv.id;
*/
