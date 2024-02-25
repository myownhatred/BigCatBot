package entities

import "gorm.io/gorm"

type AnimeOpening struct {
	gorm.Model
	Description string `json:"description"`
	Link        string `json:"link"`
}

type FreeMaw struct {
	ID          int    `json:"id"`
	Type        string `json:"typ"`
	Description string `json:"description"`
	Link        string `json:"link"`
}
