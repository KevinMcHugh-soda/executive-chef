package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"executive-chef/internal/ingredient"
)

type model struct {
	ingredients []ingredient.Ingredient
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "ctrl+c" || key.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Ingredients:\n"
	for _, ing := range m.ingredients {
		s += fmt.Sprintf("- %s (%s)\n", ing.Name, ing.Role)
	}
	return s
}

func Run(ingredients []ingredient.Ingredient) error {
	p := tea.NewProgram(model{ingredients: ingredients})
	_, err := p.Run()
	return err
}
