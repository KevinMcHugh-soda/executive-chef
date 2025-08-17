package customer_test

import (
	"math/rand"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"executive-chef/internal/customer"
	"executive-chef/internal/ingredient"
)

func TestDeckDraw(t *testing.T) {
	rand.Seed(1)
	gofakeit.Seed(1)
	ingredients := []ingredient.Ingredient{
		{Name: "Chicken", Role: ingredient.Protein},
	}
	d := customer.NewDeck(ingredients)
	require.Len(t, d.Cards, 15)
	drawn := d.Draw(3)
	assert.Len(t, drawn, 3)
	assert.Len(t, d.Cards, 12)
}
