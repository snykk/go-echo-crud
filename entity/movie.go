package entity

import (
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Publisher   string `json:"publisher"`
	Rating      string `json:"rating"`
	ReleaseYear int    `json:"release_year"`
}
