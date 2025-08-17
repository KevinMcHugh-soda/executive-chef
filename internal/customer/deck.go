package customer

import (
	"math/rand"

	"executive-chef/internal/ingredient"
)

// Deck represents a collection of customer cards.
type Deck struct {
	Cards []Customer
}

// NewDeck creates a deck containing 15 random customers.
// Customers are shuffled upon creation.
func NewDeck(ingredients []ingredient.Ingredient) *Deck {
	cards := RandomCustomers(ingredients, 15)
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	return &Deck{Cards: cards}
}

// Draw removes up to n customers from the top of the deck and returns them.
func (d *Deck) Draw(n int) []Customer {
	if n > len(d.Cards) {
		n = len(d.Cards)
	}
	drawn := make([]Customer, n)
	copy(drawn, d.Cards[:n])
	d.Cards = d.Cards[n:]
	return drawn
}
