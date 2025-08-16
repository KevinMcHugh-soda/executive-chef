package customer_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"executive-chef/internal/customer"
	"executive-chef/internal/ingredient"
)

func TestRandomCravingUniqueness(t *testing.T) {
	rand.Seed(1)
	ingredients := []ingredient.Ingredient{
		{Name: "Chicken", Role: ingredient.Protein},
		{Name: "Chicken", Role: ingredient.Protein},
		{Name: "Rice", Role: ingredient.Carb},
	}
	cr := customer.RandomCraving(ingredients)
	require.NotEmpty(t, cr.Ingredients)

	allowed := []ingredient.Ingredient{
		{Name: "Chicken", Role: ingredient.Protein},
		{Name: "Rice", Role: ingredient.Carb},
	}
	seen := map[ingredient.Ingredient]bool{}
	for _, ing := range cr.Ingredients {
		assert.Contains(t, allowed, ing)
		assert.False(t, seen[ing])
		seen[ing] = true
	}
}
