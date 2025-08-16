package game

import (
	"executive-chef/internal/customer"
	"executive-chef/internal/dish"
	"executive-chef/internal/ingredient"
)

// Event represents something that happened in the game and should be rendered by the UI.
type Event interface {
	EventType() string
}

// Phase represents the current phase of a turn.
type Phase string

const (
	PhaseDraft   Phase = "Draft"
	PhaseDesign  Phase = "Design"
	PhaseService Phase = "Service"
)

// PhaseEvent announces the current turn and phase of the game.
type PhaseEvent struct {
	Turn  int
	Phase Phase
}

func (e PhaseEvent) EventType() string { return "phase" }

// DraftOptionsEvent is sent when a new set of draftable ingredients should be shown.
type DraftOptionsEvent struct {
	Reveal []ingredient.Ingredient
	Picks  int
}

func (e DraftOptionsEvent) EventType() string { return "draft_options" }

// IngredientDraftedEvent announces that an ingredient has been drafted by the player.
type IngredientDraftedEvent struct {
	Ingredient ingredient.Ingredient
}

func (e IngredientDraftedEvent) EventType() string { return "ingredient_drafted" }

// DesignOptionsEvent is sent when the player can design dishes from drafted ingredients.
type DesignOptionsEvent struct {
	Drafted []ingredient.Ingredient
}

func (e DesignOptionsEvent) EventType() string { return "design_options" }

// DishCreatedEvent notifies the UI that a dish has been created.
type DishCreatedEvent struct {
	Dish dish.Dish
}

func (e DishCreatedEvent) EventType() string { return "dish_created" }

// DishDeletedEvent notifies the UI that a dish has been deleted.
type DishDeletedEvent struct {
	Dish dish.Dish
}

func (e DishDeletedEvent) EventType() string { return "dish_deleted" }

// ServiceResultEvent reports which dish a customer selected.
// Dish will be nil if no available dish satisfies the customer's cravings.
type ServiceResultEvent struct {
	Customer customer.Customer
	Dish     *dish.Dish
	Payment  int
	Money    int
}

func (e ServiceResultEvent) EventType() string { return "service_result" }

// ServiceEndEvent signals that all customers have been served.
type ServiceEndEvent struct{}

func (e ServiceEndEvent) EventType() string { return "service_end" }

// Action represents an input from the player relayed by the UI.
type Action interface {
	ActionType() string
}

// DraftSelectionAction is sent by the UI when the player selects an ingredient during drafting.
type DraftSelectionAction struct {
	Index int
}

func (a DraftSelectionAction) ActionType() string { return "draft_selection" }

// CreateDishAction contains information to create a new dish.
type CreateDishAction struct {
	Name    string
	Indices []int
}

func (a CreateDishAction) ActionType() string { return "create_dish" }

// DeleteDishAction removes the most recently created dish.
type DeleteDishAction struct{}

func (a DeleteDishAction) ActionType() string { return "delete_dish" }

// FinishDesignAction signals that the player is done designing dishes.
type FinishDesignAction struct{}

func (a FinishDesignAction) ActionType() string { return "finish_design" }

// ContinueAction advances the game during service or to the next turn.
type ContinueAction struct{}

func (a ContinueAction) ActionType() string { return "continue" }
