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

// DraftOptionsEvent is sent when a new set of draftable ingredients should be shown.
type DraftOptionsEvent struct {
	Reveal []ingredient.Ingredient
	Picks  int
}

func (e DraftOptionsEvent) EventType() string { return "draft_options" }

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

// ServiceResultEvent reports which dish a customer selected.
// Dish will be nil if no available dish satisfies the customer's cravings.
type ServiceResultEvent struct {
	Customer customer.Customer
	Dish     *dish.Dish
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

// FinishDesignAction signals that the player is done designing dishes.
type FinishDesignAction struct{}

func (a FinishDesignAction) ActionType() string { return "finish_design" }
