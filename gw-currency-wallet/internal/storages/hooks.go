package storages

import "gorm.io/gorm"

func (w *Wallet) AfterSave(tx *gorm.DB) (err error) {
	if w.EurBalance < 0 || w.RubBalance < 0 || w.UsdBalance < 0 {
		return ErrLessThanZero
	}

	return nil
}
