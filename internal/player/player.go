package player

import "executive-chef/internal/ingredient"

// Player represents a game participant who drafts ingredients.
type Player struct {
	Drafted []ingredient.Ingredient
}

// New creates a player with an empty drafted list.
func New() *Player {
	return &Player{Drafted: []ingredient.Ingredient{}}
}

// Add adds an ingredient to the player's drafted list.
func (p *Player) Add(ing ingredient.Ingredient) {
	p.Drafted = append(p.Drafted, ing)
}
