package game

import (
	"executive-chef/internal/deck"
	"executive-chef/internal/ingredient"
	"executive-chef/internal/player"
)

type Game struct {
	Deck           *deck.Deck
	Player         *player.Player
	Events         chan<- Event
	Actions        <-chan Action
	AllIngredients []ingredient.Ingredient
}

func New(all []ingredient.Ingredient, d *deck.Deck, p *player.Player, events chan<- Event, actions <-chan Action) *Game {
	return &Game{AllIngredients: all, Deck: d, Player: p, Events: events, Actions: actions}
}

func (g *Game) Play() {
	turn := 1
	for {
		t := Turn{Number: turn, Game: g}
		t.DraftPhase()
		t.DesignPhase()
		t.ServicePhase()
		g.Player.ResetTurn()
		turn++
	}
}
