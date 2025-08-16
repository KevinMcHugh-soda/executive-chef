package main

import (
	"log"

	"executive-chef/internal/ingredient"
	"executive-chef/internal/ui"
)

func main() {
	ingredients, err := ingredient.LoadFromFile("ingredients.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if err := ui.Run(ingredients); err != nil {
		log.Fatal(err)
	}
}
