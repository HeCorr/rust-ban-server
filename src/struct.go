package main

import "gorm.io/gorm"

// Ban struct
type Ban struct {
	gorm.Model
	SteamID    string `json:"steamId"`
	Reason     string `json:"reason"`
	ExpiryDate int64  `json:"expiryDate"`
}
