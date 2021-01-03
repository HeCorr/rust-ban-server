package main

import (
	"gorm.io/gorm"
)

func getBan(sID string) (b Ban, _ error) {
	tx := db.Session(&gorm.Session{})
	find := tx.Where("steam_id = ?", sID).Find(&b)
	if find.Error != nil {
		return b, find.Error
	}
	if find.RowsAffected == 0 {
		return b, errNotFound
	}
	return b, nil
}

func addBan(sID, reason string, expiry int64) error {
	tx := db.Session(&gorm.Session{})
	create := tx.Create(&Ban{SteamID: sID, Reason: reason, ExpiryDate: expiry})
	if create.Error != nil {
		return create.Error
	}
	if create.RowsAffected == 0 {
		return errNotInserted
	}
	return nil
}
