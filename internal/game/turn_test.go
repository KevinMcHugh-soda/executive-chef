package game

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"executive-chef/internal/customer"
	"executive-chef/internal/deck"
	"executive-chef/internal/dish"
	"executive-chef/internal/ingredient"
	"executive-chef/internal/player"
)

func TestHasIngredients(t *testing.T) {
	have := []ingredient.Ingredient{
		{Name: "Chicken", Role: ingredient.Protein},
		{Name: "Rice", Role: ingredient.Carb},
	}
	needed := []ingredient.Ingredient{{Name: "Rice", Role: ingredient.Carb}}
	assert.True(t, hasIngredients(have, needed))

	needed = append(needed, ingredient.Ingredient{Name: "Broccoli", Role: ingredient.Vegetable})
	assert.False(t, hasIngredients(have, needed))
}

func TestDraftPhaseAllowsFivePicksAfterFirstTurn(t *testing.T) {
	reveal := []ingredient.Ingredient{
		{Name: "Ing1", Role: ingredient.Protein},
		{Name: "Ing2", Role: ingredient.Protein},
		{Name: "Ing3", Role: ingredient.Protein},
		{Name: "Ing4", Role: ingredient.Protein},
		{Name: "Ing5", Role: ingredient.Protein},
		{Name: "Ing6", Role: ingredient.Protein},
		{Name: "Ing7", Role: ingredient.Protein},
		{Name: "Ing8", Role: ingredient.Protein},
		{Name: "Ing9", Role: ingredient.Protein},
		{Name: "Ing10", Role: ingredient.Protein},
	}
	d := &deck.Deck{Cards: reveal}
	p := player.New()
	events := make(chan Event, 20)
	actions := make(chan Action, 5)
	for i := 0; i < 5; i++ {
		actions <- DraftSelectionAction{Index: 0}
	}
	g := New(d, nil, p, events, actions)
	turn := Turn{Number: 2, Game: g}
	turn.DraftPhase()
	assert.Len(t, p.Drafted, 5)
}

func TestServicePhaseWaitsForFinalContinue(t *testing.T) {
	ing := ingredient.Ingredient{Name: "Chicken", Role: ingredient.Protein}
	cust := customer.Customer{
		Name:     "Patron",
		Cravings: []customer.Craving{{Ingredients: []ingredient.Ingredient{ing}}},
	}
	cdeck := &customer.Deck{Cards: []customer.Customer{cust}}
	p := player.New()
	p.Drafted = []ingredient.Ingredient{ing}
	p.Dishes = []dish.Dish{{Name: "Dish", Ingredients: []ingredient.Ingredient{ing}}}
	events := make(chan Event, 10)
	actions := make(chan Action, 1)
	g := New(nil, cdeck, p, events, actions)
	turn := Turn{Number: 1, Game: g}

	go turn.ServicePhase()

	<-events // PhaseEvent
	if _, ok := (<-events).(ServiceResultEvent); !ok {
		t.Fatal("expected service result event")
	}

	select {
	case e := <-events:
		t.Fatalf("unexpected event before continue: %T", e)
	case <-time.After(10 * time.Millisecond):
	}

	actions <- ContinueAction{}
	if _, ok := (<-events).(ServiceEndEvent); !ok {
		t.Fatal("expected service end event")
	}
}
