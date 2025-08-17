package main

import (
	"log"

	"executive-chef/internal/customer"
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
	c := customer.NewDeck(ingredients)
	p := player.New()

	events := make(chan game.Event)
	actions := make(chan game.Action)

	g := game.New(d, c, p, events, actions)

	go g.Play()

	if err := ui.Run(events, actions); err != nil {
		log.Fatal(err)
	}

	// Game ended
}
