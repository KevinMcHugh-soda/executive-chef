package player

import (
	"executive-chef/internal/dish"
	"executive-chef/internal/ingredient"
)

// Player represents a game participant who drafts ingredients and designs dishes.
type Player struct {
	Drafted []ingredient.Ingredient
	Dishes  []dish.Dish
	Money   int
}

// New creates a player with empty drafted and dish lists.
func New() *Player {
	return &Player{Drafted: []ingredient.Ingredient{}, Dishes: []dish.Dish{}, Money: 0}
}

// Add adds an ingredient to the player's drafted list.
func (p *Player) Add(ing ingredient.Ingredient) {
	p.Drafted = append(p.Drafted, ing)
}

// AddDish adds a dish to the player's designed dishes.
func (p *Player) AddDish(d dish.Dish) {
	p.Dishes = append(p.Dishes, d)
}

// RemoveLastDish removes and returns the most recently added dish.
// The second return value is false if there are no dishes to remove.
func (p *Player) RemoveLastDish() (dish.Dish, bool) {
	if len(p.Dishes) == 0 {
		return dish.Dish{}, false
	}
	d := p.Dishes[len(p.Dishes)-1]
	p.Dishes = p.Dishes[:len(p.Dishes)-1]
	return d, true
}

// RemoveDish removes and returns the dish at the given index.
// The second return value is false if the index is out of range.
func (p *Player) RemoveDish(i int) (dish.Dish, bool) {
	if i < 0 || i >= len(p.Dishes) {
		return dish.Dish{}, false
	}
	d := p.Dishes[i]
	p.Dishes = append(p.Dishes[:i], p.Dishes[i+1:]...)
	return d, true
}

// AddMoney increases the player's money by the given amount.
func (p *Player) AddMoney(amount int) {
	p.Money += amount
}

// ResetTurn clears drafted ingredients for a new turn while keeping dishes.
func (p *Player) ResetTurn() {
	p.Drafted = nil
}
