package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"executive-chef/internal/game"
	"executive-chef/internal/ingredient"
)

type model struct {
	draft   []ingredient.Ingredient
	cursor  int
	actions chan<- game.Action
}

func initialModel(actions chan<- game.Action) model {
	return model{actions: actions}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case game.DraftOptionsEvent:
		m.draft = msg.Reveal
		if m.cursor >= len(m.draft) {
			m.cursor = len(m.draft) - 1
		}
	case game.DesignOptionsEvent:
		return m, tea.Quit
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
				m.actions <- game.DraftSelectionAction{Index: m.cursor}
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
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		paneStyle.Render(draftView),
	)
}

// Run renders game events and sends player actions back to the game.
func Run(events <-chan game.Event, actions chan<- game.Action) error {
	m := initialModel(actions)
	p := tea.NewProgram(m, tea.WithAltScreen())

	designEventCh := make(chan game.DesignOptionsEvent, 1)

	go func() {
		for e := range events {
			if de, ok := e.(game.DesignOptionsEvent); ok {
				designEventCh <- de
				p.Send(e)
				p.Quit()
				return
			}
			p.Send(e)
		}
	}()

	if _, err := p.Run(); err != nil {
		return err
	}

	var design game.DesignOptionsEvent
	select {
	case design = <-designEventCh:
	default:
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter dish name (or 'done' to finish): ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if strings.EqualFold(name, "done") || name == "" {
			actions <- game.FinishDesignAction{}
			return nil
		}
		fmt.Println("Available ingredients:")
		for i, ing := range design.Drafted {
			fmt.Printf("%d) %s (%s)\n", i+1, ing.Name, ing.Role)
		}
		fmt.Print("Choose ingredient numbers separated by spaces: ")
		input, _ := reader.ReadString('\n')
		fields := strings.Fields(input)
		var indices []int
		for _, f := range fields {
			idx, err := strconv.Atoi(f)
			if err != nil || idx < 1 || idx > len(design.Drafted) {
				indices = nil
				break
			}
			indices = append(indices, idx-1)
		}
		if len(indices) == 0 {
			fmt.Println("Invalid selection.")
			continue
		}
		actions <- game.CreateDishAction{Name: name, Indices: indices}
		if evt, ok := <-events; ok {
			if created, ok := evt.(game.DishCreatedEvent); ok {
				fmt.Printf("Added dish '%s'!\n", created.Dish.Name)
			}
		}
	}
}
