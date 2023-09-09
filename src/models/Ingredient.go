package models

// Ingredient represents the ingredients table.
type Ingredient struct {
	ID          string    `json:"id"` // Unique identifier for the ingredient
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
