package customer

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"executive-chef/internal/ingredient"
)

// Craving represents a combination of ingredients a customer wants.
type Craving struct {
	Ingredients []ingredient.Ingredient
}

// Customer represents a single customer with ordered cravings and a name.
type Customer struct {
	Name       string
	Cravings   []Craving
	Constraint *ingredient.Ingredient // ingredient the customer refuses, nil if none
}

// RandomCraving returns a Craving made of random ingredients.
// Each ingredient in the resulting craving will be unique even if the
// provided slice contains duplicates.
func RandomCraving(ingredients []ingredient.Ingredient) Craving {
	if len(ingredients) == 0 {
		return Craving{}
	}

	// Deduplicate input ingredients so cravings never contain repeats.
	unique := make([]ingredient.Ingredient, 0, len(ingredients))
	seen := make(map[ingredient.Ingredient]bool)
	for _, ing := range ingredients {
		if !seen[ing] {
			seen[ing] = true
			unique = append(unique, ing)
		}
	}

	n := rand.Intn(len(unique)) + 1
	idxs := rand.Perm(len(unique))[:n]
	combo := make([]ingredient.Ingredient, 0, n)
	for _, i := range idxs {
		combo = append(combo, unique[i])
	}
	return Craving{Ingredients: combo}
}

// RandomCustomer generates a Customer with the given number of cravings.
// Cravings are ordered from most to least desired.
func RandomCustomer(ingredients []ingredient.Ingredient, numCravings int) Customer {
	if numCravings <= 0 {
		numCravings = 1
	}
	cravings := make([]Craving, numCravings)
	for i := 0; i < numCravings; i++ {
		cravings[i] = RandomCraving(ingredients)
	}

	// Choose a constraint from ingredients not already in cravings with 50% chance.
	var constraint *ingredient.Ingredient
	if len(ingredients) > 0 {
		used := make(map[ingredient.Ingredient]bool)
		for _, cr := range cravings {
			for _, ing := range cr.Ingredients {
				used[ing] = true
			}
		}
		var candidates []ingredient.Ingredient
		for _, ing := range ingredients {
			if !used[ing] {
				candidates = append(candidates, ing)
			}
		}
		if len(candidates) > 0 && rand.Intn(2) == 0 {
			c := candidates[rand.Intn(len(candidates))]
			constraint = &c
		}
	}

	return Customer{Name: gofakeit.Name(), Cravings: cravings, Constraint: constraint}
}

// RandomCustomers generates the specified number of customers.
func RandomCustomers(ingredients []ingredient.Ingredient, count int) []Customer {
	customers := make([]Customer, count)
	for i := 0; i < count; i++ {
		numCravings := rand.Intn(3) + 1
		customers[i] = RandomCustomer(ingredients, numCravings)
	}
	return customers
}

func init() {
	rand.Seed(time.Now().UnixNano())
	gofakeit.Seed(time.Now().UnixNano())
}
