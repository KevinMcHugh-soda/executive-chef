package ui

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"executive-chef/internal/customer"
	"executive-chef/internal/dish"
	"executive-chef/internal/game"
	"executive-chef/internal/ingredient"
)

func stripANSI(s string) string {
	re := regexp.MustCompile("\x1b\\[[0-9;]*m")
	return re.ReplaceAllString(s, "")
}

func TestServiceModeViewShowsAllCravings(t *testing.T) {
	c := customer.Customer{
		Name: "Alice",
		Cravings: []customer.Craving{
			{Ingredients: []ingredient.Ingredient{
				{Name: "Lettuce", Role: ingredient.Vegetable},
				{Name: "Tomato", Role: ingredient.Vegetable},
			}},
			{Ingredients: []ingredient.Ingredient{
				{Name: "Tomato", Role: ingredient.Vegetable},
				{Name: "Cheese", Role: ingredient.Protein},
			}},
		},
	}
	d := &dish.Dish{
		Name: "Salad",
		Ingredients: []ingredient.Ingredient{
			{Name: "Tomato", Role: ingredient.Vegetable},
			{Name: "Cheese", Role: ingredient.Protein},
		},
	}
	sm := serviceMode{current: &game.ServiceResultEvent{Customer: c, Dish: d, Payment: 3}}
	out := stripANSI(sm.View(&model{}))
	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Equal(t, 4, len(lines))
	assert.Equal(t, "Service", strings.TrimSpace(lines[0]))
	assert.Equal(t, "Alice", strings.TrimSpace(lines[1]))
	assert.Equal(t, "Lettuce, Tomato", strings.TrimSpace(lines[2]))
	assert.Equal(t, "Tomato, Cheese -> Salad ($3)", strings.TrimSpace(lines[3]))
}
