package main

import (
	"log"

	"executive-chef/internal/deck"
	"executive-chef/internal/game"
	"executive-chef/internal/ingredient"
	"executive-chef/internal/player"
	"executive-chef/internal/ui"
)

func main() {
	ingredients, err := ingredient.LoadFromFile("ingredients.yaml")
	if err != nil {
		log.Fatal(err)
	}

	d := deck.New(ingredients)
	p := player.New()

	events := make(chan game.Event)
	actions := make(chan game.Action)

	go func() {
		turn := 1
		for {
			t := game.Turn{Number: turn, Deck: d, Player: p, Events: events, Actions: actions}
			t.DraftPhase()
			t.DesignPhase()
			t.ServicePhase()
			p.ResetTurn()
			turn++
		}
	}()

	if err := ui.Run(events, actions); err != nil {
		log.Fatal(err)
	}

	// Game ended
}
