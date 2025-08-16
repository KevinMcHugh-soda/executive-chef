package customer

import (
	"math/rand"
	"time"

	"executive-chef/internal/ingredient"
)

// Craving represents a combination of ingredients a customer wants.
type Craving struct {
	Ingredients []ingredient.Ingredient
}

// Customer represents a single customer with ordered cravings.
type Customer struct {
	Cravings []Craving
}

// RandomCraving returns a Craving made of random ingredients.
func RandomCraving(ingredients []ingredient.Ingredient) Craving {
	if len(ingredients) == 0 {
		return Craving{}
	}
	n := rand.Intn(len(ingredients)) + 1
	idxs := rand.Perm(len(ingredients))[:n]
	combo := make([]ingredient.Ingredient, 0, n)
	for _, i := range idxs {
		combo = append(combo, ingredients[i])
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
	return Customer{Cravings: cravings}
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
}
