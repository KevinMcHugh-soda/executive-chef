package player_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"executive-chef/internal/dish"
	"executive-chef/internal/ingredient"
	"executive-chef/internal/player"
)

func TestPlayerLifecycle(t *testing.T) {
	p := player.New()
	require.NotNil(t, p)
	assert.Empty(t, p.Drafted)
	assert.Empty(t, p.Dishes)
	assert.Equal(t, 0, p.Money)

	ing := ingredient.Ingredient{Name: "Chicken", Role: ingredient.Protein}
	p.Add(ing)
	assert.Equal(t, []ingredient.Ingredient{ing}, p.Drafted)

	d := dish.Dish{Name: "Chicken Dish", Ingredients: []ingredient.Ingredient{ing}}
	p.AddDish(d)
	assert.Equal(t, []dish.Dish{d}, p.Dishes)

	d2 := dish.Dish{Name: "Veggie Dish", Ingredients: []ingredient.Ingredient{ing}}
	p.AddDish(d2)
	removed, ok := p.RemoveDish(0)
	assert.True(t, ok)
	assert.Equal(t, d, removed)
	assert.Equal(t, []dish.Dish{d2}, p.Dishes)
	_, ok = p.RemoveDish(5)
	assert.False(t, ok)

	p.AddMoney(5)
	assert.Equal(t, 5, p.Money)

	p.ResetTurn()
	assert.Empty(t, p.Drafted)
}
