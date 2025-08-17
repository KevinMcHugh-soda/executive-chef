package game

import (
	"sort"

	"executive-chef/internal/dish"
	"executive-chef/internal/ingredient"
)

// Turn represents a single turn in the game.
type Turn struct {
	Number int
	Game   *Game
}

// DraftPhase performs the drafting phase of a turn. Ten cards are revealed and
// the player may draft three of them in the first turn and five thereafter.
func (t *Turn) DraftPhase() {
	t.Game.Events <- PhaseEvent{Turn: t.Number, Phase: PhaseDraft}
	reveal := t.Game.Deck.Draw(10)
	roleOrder := map[ingredient.Role]int{
		ingredient.Protein:   0,
		ingredient.Vegetable: 1,
		ingredient.Carb:      2,
	}
	sort.Slice(reveal, func(i, j int) bool {
		ri := roleOrder[reveal[i].Role]
		rj := roleOrder[reveal[j].Role]
		if ri != rj {
			return ri < rj
		}
		return reveal[i].Name < reveal[j].Name
	})
	remaining := 3
	if t.Number > 1 {
		remaining = 5
	}
	t.Game.Events <- DraftOptionsEvent{Reveal: reveal, Picks: remaining}
	for remaining > 0 && len(reveal) > 0 {
		act := <-t.Game.Actions
		sel, ok := act.(DraftSelectionAction)
		if !ok || sel.Index < 0 || sel.Index >= len(reveal) {
			continue
		}
		chosen := reveal[sel.Index]
		t.Game.Player.Add(chosen)
		t.Game.Events <- IngredientDraftedEvent{Ingredient: chosen}
		reveal = append(reveal[:sel.Index], reveal[sel.Index+1:]...)
		remaining--
		if remaining > 0 && len(reveal) > 0 {
			t.Game.Events <- DraftOptionsEvent{Reveal: reveal, Picks: remaining}
		}
	}
}

// DesignPhase allows the player to combine drafted ingredients into named dishes.
// The player can create up to two dishes this turn and may have up to ten dishes
// overall. The phase ends when a FinishDesignAction is received.
func (t *Turn) DesignPhase() {
	t.Game.Events <- PhaseEvent{Turn: t.Number, Phase: PhaseDesign}
	t.Game.Events <- DesignOptionsEvent{Drafted: t.Game.Player.Drafted}
	created := []int{}

	for {
		act := <-t.Game.Actions
		switch a := act.(type) {
		case CreateDishAction:
			if a.Name == "" || len(created) >= 2 || len(t.Game.Player.Dishes) >= 10 {
				continue
			}
			used := make(map[int]bool)
			var dishIngs []ingredient.Ingredient
			valid := true
			for _, idx := range a.Indices {
				if idx < 0 || idx >= len(t.Game.Player.Drafted) || used[idx] {
					valid = false
					break
				}
				used[idx] = true
				dishIngs = append(dishIngs, t.Game.Player.Drafted[idx])
			}
			if !valid || len(dishIngs) == 0 {
				continue
			}
			d := dish.Dish{Name: a.Name, Ingredients: dishIngs}
			t.Game.Player.AddDish(d)
			created = append(created, len(t.Game.Player.Dishes)-1)
			t.Game.Events <- DishCreatedEvent{Dish: d}
		case DeleteDishAction:
			if a.Index < 0 || a.Index >= len(t.Game.Player.Dishes) {
				continue
			}
			d, ok := t.Game.Player.RemoveDish(a.Index)
			if ok {
				for i, idx := range created {
					if idx == a.Index {
						created = append(created[:i], created[i+1:]...)
						break
					}
				}
				for i := range created {
					if created[i] > a.Index {
						created[i]--
					}
				}
				t.Game.Events <- DishDeletedEvent{Dish: d, Index: a.Index}
			}
		case FinishDesignAction:
			return
		}
	}
}

// ServicePhase presents dishes to customers who choose based on their cravings.
func (t *Turn) ServicePhase() {
	t.Game.Events <- PhaseEvent{Turn: t.Number, Phase: PhaseService}
	customers := t.Game.Customers.Draw(3)
	var available []dish.Dish
	for _, d := range t.Game.Player.Dishes {
		if hasIngredients(t.Game.Player.Drafted, d.Ingredients) {
			available = append(available, d)
		}
	}
	for i, c := range customers {
		bestIdx := -1
		bestScore := 0
		bestCraving := -1
		for i, d := range available {
			if c.Constraint != nil {
				rejected := false
				for _, ing := range d.Ingredients {
					if ing == *c.Constraint {
						rejected = true
						break
					}
				}
				if rejected {
					continue
				}
			}
			score := 0
			cravingIdx := -1
			for j, cr := range c.Cravings {
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
					cravingIdx = j
				}
			}
			if score > bestScore {
				bestScore = score
				bestIdx = i
				bestCraving = cravingIdx
			}
		}
		var chosen *dish.Dish
		payment := 0
		if bestIdx >= 0 {
			d := available[bestIdx]
			chosen = &d
			switch bestCraving {
			case 0:
				payment = 5
			case 1:
				payment = 3
			case 2:
				payment = 1
			}
			t.Game.Player.AddMoney(payment)
		}
		t.Game.Events <- ServiceResultEvent{Customer: c, Dish: chosen, Payment: payment, Money: t.Game.Player.Money}
		if i < len(customers)-1 {
			for {
				if _, ok := (<-t.Game.Actions).(ContinueAction); ok {
					break
				}
			}
		}
	}
	t.Game.Events <- ServiceEndEvent{}
	for {
		if _, ok := (<-t.Game.Actions).(ContinueAction); ok {
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
