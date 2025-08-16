package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"executive-chef/internal/deck"
	"executive-chef/internal/dish"
	"executive-chef/internal/ingredient"
	"executive-chef/internal/player"
)

// Turn represents a single turn in the game.
type Turn struct {
	Deck   *deck.Deck
	Player *player.Player
}

// DraftPhase performs the drafting phase of a turn. Ten cards are revealed
// and the player may draft three of them.
func (t *Turn) DraftPhase() {
	reveal := t.Deck.Draw(10)
	reader := bufio.NewReader(os.Stdin)
	for picks := 0; picks < 3 && len(reveal) > 0; picks++ {
		fmt.Println("Choose an ingredient to draft:")
		for i, card := range reveal {
			fmt.Printf("%d) %s (%s)\n", i+1, card.Name, card.Role)
		}
		fmt.Print("Selection: ")
		input, _ := reader.ReadString('\n')
		n, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || n < 1 || n > len(reveal) {
			fmt.Println("Invalid selection")
			picks--
			continue
		}
		chosen := reveal[n-1]
		t.Player.Add(chosen)
		reveal = append(reveal[:n-1], reveal[n:]...)
	}
}

// DesignPhase allows the player to combine drafted ingredients into named dishes.
// The player can create multiple dishes until they enter "done" as the dish name.
func (t *Turn) DesignPhase() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter dish name (or 'done' to finish): ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if strings.EqualFold(name, "done") || name == "" {
			break
		}

		if len(t.Player.Drafted) == 0 {
			fmt.Println("No ingredients available to create dishes.")
			break
		}

		fmt.Println("Available ingredients:")
		for i, ing := range t.Player.Drafted {
			fmt.Printf("%d) %s (%s)\n", i+1, ing.Name, ing.Role)
		}
		fmt.Print("Choose ingredient numbers separated by spaces: ")
		input, _ := reader.ReadString('\n')
		fields := strings.Fields(input)

		var dishIngs []ingredient.Ingredient
		used := make(map[int]bool)
		valid := true
		for _, f := range fields {
			idx, err := strconv.Atoi(f)
			if err != nil || idx < 1 || idx > len(t.Player.Drafted) || used[idx] {
				fmt.Println("Invalid ingredient selection.")
				valid = false
				break
			}
			used[idx] = true
			dishIngs = append(dishIngs, t.Player.Drafted[idx-1])
		}
		if !valid || len(dishIngs) == 0 {
			fmt.Println("No dish created.")
			continue
		}

		t.Player.AddDish(dish.Dish{Name: name, Ingredients: dishIngs})
		fmt.Printf("Added dish '%s'!\n", name)
	}
}
