package main

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func addBan(b Ban) error {
	tx := db.Session(&gorm.Session{})
	create := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&b)
	if create.Error != nil {
		return create.Error
	}
	if create.RowsAffected == 0 {
		return errNotInserted
	}
	return nil
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
