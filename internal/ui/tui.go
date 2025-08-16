package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"executive-chef/internal/ingredient"
)

type model struct {
	draft  []ingredient.Ingredient
	hand   []ingredient.Ingredient
	cursor int
}

func initialModel(ingredients []ingredient.Ingredient) model {
	return model{draft: ingredients}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.draft)-1 {
				m.cursor++
			}
		case "enter", " ":
			if len(m.draft) > 0 {
				selected := m.draft[m.cursor]
				m.hand = append(m.hand, selected)
				m.draft = append(m.draft[:m.cursor], m.draft[m.cursor+1:]...)
				if m.cursor >= len(m.draft) && m.cursor > 0 {
					m.cursor--
				}
			}
		}
	}
	return m, nil
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	paneStyle  = lipgloss.NewStyle().Padding(0, 1)
)

func (m model) View() string {
	draftView := titleStyle.Render("Draftable Ingredients:") + "\n"
	for i, ing := range m.draft {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		draftView += fmt.Sprintf("%s %s (%s)\n", cursor, ing.Name, ing.Role)
	}

	handView := titleStyle.Render("Your Hand:") + "\n"
	if len(m.hand) == 0 {
		handView += " (empty)\n"
	} else {
		for _, ing := range m.hand {
			handView += fmt.Sprintf("- %s (%s)\n", ing.Name, ing.Role)
		}
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		paneStyle.Render(draftView),
		paneStyle.Render(handView),
	)
}

func Run(ingredients []ingredient.Ingredient) error {
	p := tea.NewProgram(initialModel(ingredients), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
