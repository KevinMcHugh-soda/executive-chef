package deck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"executive-chef/internal/deck"
	"executive-chef/internal/ingredient"
)

func TestNewDeck(t *testing.T) {
	all := []ingredient.Ingredient{
		{Name: "Chicken", Role: ingredient.Protein},
		{Name: "Rice", Role: ingredient.Carb},
		{Name: "Broccoli", Role: ingredient.Vegetable},
	}
	d := deck.New(all)
	require.NotNil(t, d)
	assert.Equal(t, 50, len(d.Cards))
	for _, card := range d.Cards {
		assert.Contains(t, all, card)
	}
}

func TestDraw(t *testing.T) {
	all := []ingredient.Ingredient{
		{Name: "Chicken", Role: ingredient.Protein},
		{Name: "Rice", Role: ingredient.Carb},
	}
	d := deck.New(all)
	drawn := d.Draw(10)
	assert.Len(t, drawn, 10)
	assert.Len(t, d.Cards, 40)
	drawn = d.Draw(100)
	assert.Len(t, drawn, 40)
	assert.Empty(t, d.Cards)
}
