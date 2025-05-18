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
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	Excerpt    string `json:"excerpt"`
	Content    string `json:"content"`
	Date       string `json:"date"`
	ReadTime   int    `json:"readTime"`
	Category   string `json:"category"`
	CoverImage string `json:"coverImage"`
	Featured   bool   `json:"featured"`
	AuthorID   uint   `json:"authorId"`
	Author     Author `json:"author" gorm:"foreignKey:AuthorID"`
}
