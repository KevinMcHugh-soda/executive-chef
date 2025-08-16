package game

import (
	"executive-chef/internal/deck"
	"executive-chef/internal/dish"
	"executive-chef/internal/ingredient"
	"executive-chef/internal/player"
)

// Turn represents a single turn in the game.
type Turn struct {
	Deck    *deck.Deck
	Player  *player.Player
	Events  chan<- Event
	Actions <-chan Action
}

// DraftPhase performs the drafting phase of a turn. Ten cards are revealed
// and the player may draft three of them.
func (t *Turn) DraftPhase() {
	reveal := t.Deck.Draw(10)
	remaining := 3
	t.Events <- DraftOptionsEvent{Reveal: reveal, Picks: remaining}
	for remaining > 0 && len(reveal) > 0 {
		act := <-t.Actions
		sel, ok := act.(DraftSelectionAction)
		if !ok || sel.Index < 0 || sel.Index >= len(reveal) {
			continue
		}
		chosen := reveal[sel.Index]
		t.Player.Add(chosen)
		reveal = append(reveal[:sel.Index], reveal[sel.Index+1:]...)
		remaining--
		if remaining > 0 && len(reveal) > 0 {
			t.Events <- DraftOptionsEvent{Reveal: reveal, Picks: remaining}
		}
	}
}

// DesignPhase allows the player to combine drafted ingredients into named dishes.
// The player can create multiple dishes until a FinishDesignAction is received.
func (t *Turn) DesignPhase() {
	t.Events <- DesignOptionsEvent{Drafted: t.Player.Drafted}
	for {
		act := <-t.Actions
		switch a := act.(type) {
		case CreateDishAction:
			if a.Name == "" {
				continue
			}
			used := make(map[int]bool)
			var dishIngs []ingredient.Ingredient
			valid := true
			for _, idx := range a.Indices {
				if idx < 0 || idx >= len(t.Player.Drafted) || used[idx] {
					valid = false
					break
				}
				used[idx] = true
				dishIngs = append(dishIngs, t.Player.Drafted[idx])
			}
			if !valid || len(dishIngs) == 0 {
				continue
			}
			d := dish.Dish{Name: a.Name, Ingredients: dishIngs}
			t.Player.AddDish(d)
			t.Events <- DishCreatedEvent{Dish: d}
		case FinishDesignAction:
			return
		}
	}
}
