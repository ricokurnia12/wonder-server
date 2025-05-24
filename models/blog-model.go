package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name   string `json:"name"`
	Role   string `json:"role"`
	Avatar string `json:"avatar"`
}

type BlogPost struct {
	gorm.Model
	Title          string `json:"title"`
	EnglishTitle   string `json:"english_title"`
	Slug           string `json:"slug" gorm:"uniqueIndex"`
	Excerpt        string `json:"excerpt"`
	EnglishExcerpt string `json:"english_excerpt"`
	Content        string `json:"content"`
	EnglishContent string `json:"englishcontent"`
	Date           string `json:"date"`
	ReadTime       int    `json:"readTime"`
	Category       string `json:"category"`
	CoverImage     string `json:"coverImage"`
	Featured       bool   `json:"featured"`
}
