package main

import (
	"errors"

	"gorm.io/gorm"
)

func getBan(sID string) (b Ban, _ error) {
	tx := db.Session(&gorm.Session{})
	find := tx.Select("steam_id", "reason", "expiry_date").Where("steam_id = ?", sID).Find(&b)
	if find.Error != nil {
		return b, find.Error
	}
	if find.RowsAffected == 0 {
		return b, errNotFound
	}
	return b, nil
}

func addBan(b Ban) (upd bool, _ error) {
	_, err := getBan(b.SteamID)
	if err != nil && !errors.Is(err, errNotFound) {
		return false, err
	} else if err == nil {
		return true, updateBan(b)
	}
	tx := db.Session(&gorm.Session{})
	create := tx.Create(&b)
	if create.Error != nil {
		return false, create.Error
	}
	if create.RowsAffected == 0 {
		return false, errNotInserted
	}
	return false, nil
}

func delBan(sID string) error {
	var b Ban
	tx := db.Session(&gorm.Session{})
	delete := tx.Where("steam_id = ?", sID).Unscoped().Delete(&b)
	if delete.Error != nil {
		return delete.Error
	}
	if delete.RowsAffected == 0 {
		return errNotDeleted
	}
	return nil
}

func updateBan(b Ban) error {
	tx := db.Session(&gorm.Session{})
	update := tx.Model(&b).Where("steam_id = ?", b.SteamID).Updates(&b)
	if update.Error != nil {
		return update.Error
	}
	if update.RowsAffected == 0 {
		return errNotUpdated
	}
	return nil
}
