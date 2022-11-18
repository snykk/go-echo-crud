package dto

import "github.com/snykk/go-echo-crud/entity"

// request
type Movie struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Publisher   string `json:"publisher" validate:"required"`
	Rating      string `json:"rating" validate:"required"`
	ReleaseYear int    `json:"release_year" validate:"required"`
}

func (m *Movie) ToEntity() entity.Movie {
	return entity.Movie{
		Title:       m.Title,
		Description: m.Description,
		Publisher:   m.Publisher,
		Rating:      m.Rating,
		ReleaseYear: m.ReleaseYear,
	}
}
