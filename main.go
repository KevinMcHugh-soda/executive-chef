package main

import (
	"fmt"
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

	t := game.Turn{Deck: d, Player: p, Events: events, Actions: actions}

	go func() {
		t.DraftPhase()
		t.DesignPhase()
		close(events)
	}()

	if err := ui.Run(events, actions); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Your Dishes:")
	for _, dish := range p.Dishes {
		fmt.Printf("- %s: ", dish.Name)
		for i, ing := range dish.Ingredients {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(ing.Name)
		}
		fmt.Println()
	}
}
