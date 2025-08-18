package ui

import (
	"testing"

	"executive-chef/internal/ingredient"
)

func TestDefaultDishNameVegetableSalad(t *testing.T) {
	drafted := []ingredient.Ingredient{
		{Name: "Lettuce", Role: ingredient.Vegetable},
		{Name: "Tomato", Role: ingredient.Vegetable},
	}
	selected := map[int]bool{0: true, 1: true}
	got := defaultDishName(selected, drafted)
	want := "Lettuce Salad"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestDefaultDishNameNonVegetable(t *testing.T) {
	drafted := []ingredient.Ingredient{
		{Name: "Chicken", Role: ingredient.Protein},
		{Name: "Rice", Role: ingredient.Carb},
	}
	selected := map[int]bool{0: true, 1: true}
	got := defaultDishName(selected, drafted)
	want := "Chicken and Rice"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
