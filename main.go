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

	t := game.Turn{Deck: d, Player: p}
	t.DraftPhase()

	if err := ui.Run(p.Drafted); err != nil {
		log.Fatal(err)
	}
}
