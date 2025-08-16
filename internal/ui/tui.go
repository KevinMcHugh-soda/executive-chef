package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"executive-chef/internal/game"
	"executive-chef/internal/ingredient"
)

// uiMode represents a UI mode for a game phase.
type uiMode interface {
	Init(*model) tea.Cmd
	Update(*model, tea.Msg) (uiMode, tea.Cmd)
	View(*model) string
}

type model struct {
	actions chan<- game.Action
	mode    uiMode
}

func initialModel(actions chan<- game.Action) *model {
	m := &model{actions: actions}
	m.mode = &draftMode{}
	return m
}

func (m *model) Init() tea.Cmd {
	return m.mode.Init(m)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if newMode, cmd := m.mode.Update(m, msg); newMode != nil {
		m.mode = newMode
		initCmd := m.mode.Init(m)
		return m, tea.Batch(cmd, initCmd)
	} else {
		return m, cmd
	}
}

func (m *model) View() string {
	return m.mode.View(m)
}

// ---- Draft Mode ----
type draftMode struct {
	draft  []ingredient.Ingredient
	cursor int
}

func (d *draftMode) Init(m *model) tea.Cmd { return nil }

func (d *draftMode) Update(m *model, msg tea.Msg) (uiMode, tea.Cmd) {
	switch msg := msg.(type) {
	case game.DraftOptionsEvent:
		d.draft = msg.Reveal
		if d.cursor >= len(d.draft) {
			d.cursor = len(d.draft) - 1
		}
	case game.DesignOptionsEvent:
		return &designMode{drafted: msg.Drafted}, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return nil, tea.Quit
		case "up", "k":
			if d.cursor > 0 {
				d.cursor--
			}
		case "down", "j":
			if d.cursor < len(d.draft)-1 {
				d.cursor++
			}
		case "enter", " ":
			if len(d.draft) > 0 {
				m.actions <- game.DraftSelectionAction{Index: d.cursor}
			}
		}
	}
	return nil, nil
}

func (d *draftMode) View(m *model) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Draftable Ingredients:") + "\n")
	for i, ing := range d.draft {
		cursor := " "
		if d.cursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s (%s)\n", cursor, ing.Name, ing.Role))
	}
	return paneStyle.Render(b.String())
}

// ---- Design Mode ----
type designMode struct {
	drafted  []ingredient.Ingredient
	cursor   int
	selected map[int]bool
	name     textinput.Model
	message  string
}

func (d *designMode) Init(m *model) tea.Cmd {
	d.selected = make(map[int]bool)
	d.name = textinput.New()
	d.name.Placeholder = "Dish name"
	d.name.Focus()
	return nil
}

func (d *designMode) Update(m *model, msg tea.Msg) (uiMode, tea.Cmd) {
	var cmd tea.Cmd
	d.name, cmd = d.name.Update(msg)
	switch msg := msg.(type) {
	case game.DishCreatedEvent:
		d.message = fmt.Sprintf("Added dish '%s'!", msg.Dish.Name)
		d.name.SetValue("")
		d.selected = make(map[int]bool)
	case game.ServiceResultEvent:
		return &serviceMode{results: []game.ServiceResultEvent{msg}}, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.actions <- game.FinishDesignAction{}
			return nil, tea.Quit
		case "up", "k":
			if d.cursor > 0 {
				d.cursor--
			}
		case "down", "j":
			if d.cursor < len(d.drafted)-1 {
				d.cursor++
			}
		case " ":
			if d.selected[d.cursor] {
				delete(d.selected, d.cursor)
			} else {
				d.selected[d.cursor] = true
			}
		case "enter":
			name := strings.TrimSpace(d.name.Value())
			if name != "" && len(d.selected) > 0 {
				var indices []int
				for i := range d.drafted {
					if d.selected[i] {
						indices = append(indices, i)
					}
				}
				m.actions <- game.CreateDishAction{Name: name, Indices: indices}
			}
		case "f":
			m.actions <- game.FinishDesignAction{}
			return nil, nil
		}
	}
	return nil, cmd
}

func (d *designMode) View(m *model) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Design Dishes") + "\n")
	for i, ing := range d.drafted {
		cursor := " "
		if d.cursor == i {
			cursor = ">"
		}
		mark := " "
		if d.selected[i] {
			mark = "*"
		}
		b.WriteString(fmt.Sprintf("%s%s %s (%s)\n", cursor, mark, ing.Name, ing.Role))
	}
	b.WriteString("\n" + d.name.View() + "\n")
	if d.message != "" {
		b.WriteString(d.message + "\n")
	}
	b.WriteString("\nspace: select • enter: create dish • f: finish\n")
	return paneStyle.Render(b.String())
}

// ---- Service Mode ----
type serviceMode struct {
	results  []game.ServiceResultEvent
	finished bool
}

func (s *serviceMode) Init(m *model) tea.Cmd { return nil }

func (s *serviceMode) Update(m *model, msg tea.Msg) (uiMode, tea.Cmd) {
	switch msg := msg.(type) {
	case game.ServiceResultEvent:
		s.results = append(s.results, msg)
	case game.ServiceEndEvent:
		s.finished = true
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return nil, tea.Quit
		}
	}
	return nil, nil
}

func (s *serviceMode) View(m *model) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Service Results") + "\n")
	for _, r := range s.results {
		var craving []string
		if len(r.Customer.Cravings) > 0 {
			for _, ing := range r.Customer.Cravings[0].Ingredients {
				craving = append(craving, ing.Name)
			}
		}
		b.WriteString(fmt.Sprintf("%s -> ", strings.Join(craving, ", ")))
		if r.Dish != nil {
			b.WriteString(r.Dish.Name)
		} else {
			b.WriteString("no dish")
		}
		b.WriteString("\n")
	}
	if s.finished {
		b.WriteString("\nPress q to quit\n")
	}
	return paneStyle.Render(b.String())
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	paneStyle  = lipgloss.NewStyle().Padding(0, 1)
)

// Run renders game events and sends player actions back to the game.
func Run(events <-chan game.Event, actions chan<- game.Action) error {
	m := initialModel(actions)
	p := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		for e := range events {
			p.Send(e)
		}
	}()

	_, err := p.Run()
	return err
}
