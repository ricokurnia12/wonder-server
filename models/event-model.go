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
	Title              string       `json:"title"`
	EnglishTitle       string       `json:"english_title"`
	Slug               string       `json:"slug" gorm:"uniqueIndex"`
	Description        string       `json:"description"`
	EnglishDescription string       `json:"english_description"`
	Content            string       `json:"content"`
	EnglishContent     string       `json:"englishcontent"`
	Date               time.Time    `json:"date"`
	Location           string       `json:"location"`
	Province           string       `json:"province"`
	Category           CategoryType `json:"category"`
	Image              string       `json:"image"`
	Featured           bool         `json:"featured"`
	MapUrl             string       `json:"map_url"`
}
