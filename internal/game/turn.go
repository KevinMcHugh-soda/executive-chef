package game

import (
	"executive-chef/internal/customer"
	"executive-chef/internal/deck"
	"executive-chef/internal/dish"
	"executive-chef/internal/ingredient"
	"executive-chef/internal/player"
)

// Turn represents a single turn in the game.
type Turn struct {
	Number  int
	Deck    *deck.Deck
	Player  *player.Player
	Events  chan<- Event
	Actions <-chan Action
}

// DraftPhase performs the drafting phase of a turn. Ten cards are revealed
// and the player may draft three of them.
func (t *Turn) DraftPhase() {
	t.Events <- PhaseEvent{Turn: t.Number, Phase: PhaseDraft}
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
		t.Events <- IngredientDraftedEvent{Ingredient: chosen}
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
	t.Events <- PhaseEvent{Turn: t.Number, Phase: PhaseDesign}
	t.Events <- DesignOptionsEvent{Drafted: t.Player.Drafted}
	for {
		act := <-t.Actions
		switch a := act.(type) {
		case CreateDishAction:
			if len(t.Player.Dishes) >= 2 || a.Name == "" {
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

// ServicePhase presents dishes to customers who choose based on their cravings.
func (t *Turn) ServicePhase() {
	t.Events <- PhaseEvent{Turn: t.Number, Phase: PhaseService}
	customers := customer.RandomCustomers(t.Player.Drafted, 3)
	var available []dish.Dish
	for _, d := range t.Player.Dishes {
		if hasIngredients(t.Player.Drafted, d.Ingredients) {
			available = append(available, d)
		}
	}
	for i, c := range customers {
		bestIdx := -1
		bestScore := 0
		for i, d := range available {
			score := 0
			for _, cr := range c.Cravings {
				count := 0
				for _, ing := range cr.Ingredients {
					for _, ding := range d.Ingredients {
						if ing == ding {
							count++
							break
						}
					}
				}
				if count > score {
					score = count
				}
			}
			if score > bestScore {
				bestScore = score
				bestIdx = i
			}
		}
		var chosen *dish.Dish
		if bestIdx >= 0 {
			d := available[bestIdx]
			chosen = &d
			available = append(available[:bestIdx], available[bestIdx+1:]...)
		}
		t.Events <- ServiceResultEvent{Customer: c, Dish: chosen}
		if i < len(customers)-1 {
			for {
				if _, ok := (<-t.Actions).(ContinueAction); ok {
					break
				}
			}
		}
	}
	t.Events <- ServiceEndEvent{}
	for {
		if _, ok := (<-t.Actions).(ContinueAction); ok {
			break
		}
	}
}

func hasIngredients(have []ingredient.Ingredient, needed []ingredient.Ingredient) bool {
	for _, n := range needed {
		found := false
		for _, h := range have {
			if h == n {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
