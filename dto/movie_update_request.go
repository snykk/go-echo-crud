package dto

import "github.com/snykk/go-echo-crud/entity"

// request
type MovieUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Publisher   string `json:"publisher"`
	Rating      string `json:"rating"`
	ReleaseYear int    `json:"release_year"`
}

func (m *MovieUpdateRequest) UpdateToEntity() entity.Movie {
	return entity.Movie{
		Title:       m.Title,
		Description: m.Description,
		Publisher:   m.Publisher,
		Rating:      m.Rating,
		ReleaseYear: m.ReleaseYear,
	}
}
