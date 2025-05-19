package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	Title    string `json:"title"`
	FilePath string `json:"file_path"`
}
