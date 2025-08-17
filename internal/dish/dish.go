package dish

import "executive-chef/internal/ingredient"

// MaxIngredients is the maximum number of ingredients allowed in a dish.
const MaxIngredients = 3

// Dish represents a named combination of ingredients.
type Dish struct {
	Name        string
	Ingredients []ingredient.Ingredient
}
