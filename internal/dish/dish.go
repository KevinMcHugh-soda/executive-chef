package dish

import "executive-chef/internal/ingredient"

// Dish represents a named combination of ingredients.
type Dish struct {
	Name        string
	Ingredients []ingredient.Ingredient
}
