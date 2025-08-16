package deck

import (
	"math/rand"
	"time"

	"executive-chef/internal/ingredient"
)

// Deck represents a collection of ingredient cards.
type Deck struct {
	Cards []ingredient.Ingredient
}

// New creates a new deck containing 50 cards randomly chosen
// from the provided ingredient list. Ingredients can repeat.
func New(all []ingredient.Ingredient) *Deck {
	rand.Seed(time.Now().UnixNano())
	cards := make([]ingredient.Ingredient, 50)
	for i := 0; i < 50; i++ {
		cards[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	return &Deck{Cards: cards}
}

// Draw removes n cards from the top of the deck and returns them.
func (d *Deck) Draw(n int) []ingredient.Ingredient {
	if n > len(d.Cards) {
		n = len(d.Cards)
	}
	drawn := make([]ingredient.Ingredient, n)
	copy(drawn, d.Cards[:n])
	d.Cards = d.Cards[n:]
	return drawn
}
