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

func TestDesignPhaseRejectsDishesWithMoreThanThreeIngredients(t *testing.T) {
	p := player.New()
	p.Drafted = []ingredient.Ingredient{
		{Name: "Ing1", Role: ingredient.Protein},
		{Name: "Ing2", Role: ingredient.Protein},
		{Name: "Ing3", Role: ingredient.Protein},
		{Name: "Ing4", Role: ingredient.Protein},
	}
	events := make(chan Event, 10)
	actions := make(chan Action, 10)
	g := New(nil, nil, p, events, actions)
	turn := Turn{Number: 1, Game: g}

	done := make(chan struct{})
	go func() {
		turn.DesignPhase()
		close(done)
	}()

	actions <- CreateDishAction{Name: "TooMany", Indices: []int{0, 1, 2, 3}}
	actions <- FinishDesignAction{}
	<-done

	assert.Empty(t, p.Dishes)
}

func TestServicePhaseRejectsNonMatchingDish(t *testing.T) {
	p := player.New()
	chicken := ingredient.Ingredient{Name: "Chicken", Role: ingredient.Protein}
	p.Drafted = []ingredient.Ingredient{chicken}
	p.Dishes = []dish.Dish{{Name: "Chicken Dish", Ingredients: []ingredient.Ingredient{chicken}}}

	cust := customer.Customer{
		Name: "Customer",
		Cravings: []customer.Craving{
			{Ingredients: []ingredient.Ingredient{{Name: "Tomato", Role: ingredient.Vegetable}}},
		},
	}
	customers := &customer.Deck{Cards: []customer.Customer{cust}}

	events := make(chan Event, 3)
	actions := make(chan Action, 1)
	actions <- ContinueAction{}

	g := New(nil, customers, p, events, actions)
	turn := Turn{Number: 1, Game: g}
	turn.ServicePhase()

	<-events // phase event
	sr := (<-events).(ServiceResultEvent)
	assert.Nil(t, sr.Dish)
	assert.Equal(t, 0, sr.Payment)
	assert.Equal(t, 0, p.Money)
}
