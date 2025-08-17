package game

import (
	"executive-chef/internal/customer"
	"executive-chef/internal/deck"
	"executive-chef/internal/player"
)

type Game struct {
	Deck      *deck.Deck
	Customers *customer.Deck
	Player    *player.Player
	Events    chan<- Event
	Actions   <-chan Action
}

func New(d *deck.Deck, c *customer.Deck, p *player.Player, events chan<- Event, actions <-chan Action) *Game {
	return &Game{Deck: d, Customers: c, Player: p, Events: events, Actions: actions}
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
