package main

import "gorm.io/gorm"

// Ban struct
type Ban struct {
	gorm.Model `json:"-"`
	SteamID    string `json:"steamId" gorm:"uniqueIndex"`
	Reason     string `json:"reason"`
	ExpiryDate int64  `json:"expiryDate"`
}
