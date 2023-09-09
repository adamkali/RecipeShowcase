package models

import "bytes"

// Recipe represents the recipe table.
type Recipe struct {
	ID           string      `json:"id"` // Unique identifier for the recipe
	Name         string      `json:"name"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []string    `json:"instructions"`
	PictureURL   string      `json:"picture_url"`
}

func (r Recipe) preview() bytes.Buffer {
    // read from the template using gin 

}
