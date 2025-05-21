package models

import (
	"time"

	"gorm.io/gorm"
)

type CategoryType string

const (
	CategoryMusic   CategoryType = "music"
	CategoryArt     CategoryType = "art"
	CategoryCulture CategoryType = "culture"
)

type Event struct {
	gorm.Model
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Content        string       `json:"content"`
	EnglishContent string       `json:"englishcontent"`
	Date           time.Time    `json:"date"`
	Location       string       `json:"location"`
	Province       string       `json:"province"`
	Category       CategoryType `json:"category"`
	Image          string       `json:"image"`
	Featured       bool         `json:"featured"`
}
