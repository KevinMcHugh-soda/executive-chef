package game

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"executive-chef/internal/ingredient"
)

func TestHasIngredients(t *testing.T) {
	have := []ingredient.Ingredient{
		{Name: "Chicken", Role: ingredient.Protein},
		{Name: "Rice", Role: ingredient.Carb},
	}
	needed := []ingredient.Ingredient{{Name: "Rice", Role: ingredient.Carb}}
	assert.True(t, hasIngredients(have, needed))

	needed = append(needed, ingredient.Ingredient{Name: "Broccoli", Role: ingredient.Vegetable})
	assert.False(t, hasIngredients(have, needed))
}
