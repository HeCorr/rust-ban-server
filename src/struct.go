package main

import (
	"errors"

	"gorm.io/gorm"
)

var (
	errNotFound    = errors.New("not found")
	errNotInserted = errors.New("not inserted")
	errNotDeleted  = errors.New("not deleted")
	errNotUpdated  = errors.New("not updated")
)

// Ban struct
type Ban struct {
	gorm.Model `json:"-"`
	SteamID    string `json:"steamId" gorm:"uniqueIndex"`
	Reason     string `json:"reason"`
	ExpiryDate int64  `json:"expiryDate"`
}
