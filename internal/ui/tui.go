package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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
	// Status returns a short help prompt for the status bar.
	Status(*model) string
}

type model struct {
	actions chan<- game.Action
	mode    uiMode
	events  []string
	vp      viewport.Model
	width   int
	turn    int
	phase   game.Phase
}

func initialModel(actions chan<- game.Action) *model {
	m := &model{actions: actions}
	m.mode = &draftMode{}
	m.vp = viewport.New(28, 7)
	return m
}

func (m *model) Init() tea.Cmd {
	return m.mode.Init(m)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.vp.Width = 28
		m.vp.Height = msg.Height - 4
	case game.Event:
		m.events = append(m.events, eventString(msg))
		m.vp.SetContent(strings.Join(m.events, "\n"))
		m.vp.GotoBottom()
		if info, ok := msg.(game.PhaseEvent); ok {
			m.turn = info.Turn
			m.phase = info.Phase
		}
	}

	var vpCmd tea.Cmd
	m.vp, vpCmd = m.vp.Update(msg)

	newMode, modeCmd := m.mode.Update(m, msg)
	if newMode != nil {
		m.mode = newMode
		initCmd := m.mode.Init(m)
		return m, tea.Batch(modeCmd, initCmd, vpCmd)
	}

	return m, tea.Batch(vpCmd, modeCmd)
}

func (m *model) View() string {
	mainView := m.mode.View(m)
	mainWidth := m.width - 30
	if mainWidth < 0 {
		mainWidth = 0
	}
	main := lipgloss.NewStyle().Width(mainWidth).Render(mainView)

	info := paneStyle.Render(titleStyle.Render("Game Info") + "\n" + fmt.Sprintf("Turn: %d\nPhase: %s", m.turn, m.phase))
	logView := paneStyle.Width(30).Render(titleStyle.Render("Events") + "\n" + m.vp.View())

	content := lipgloss.JoinHorizontal(lipgloss.Top, main, logView)
	status := statusStyle.Render(m.mode.Status(m))
	return lipgloss.JoinVertical(lipgloss.Left, info, content, status)
}

func eventString(e game.Event) string {
	switch e := e.(type) {
	case game.PhaseEvent:
		return fmt.Sprintf("Turn %d: %s phase", e.Turn, e.Phase)
	case game.DraftOptionsEvent:
		var names []string
		for _, ing := range e.Reveal {
			names = append(names, ing.Name)
		}
		return fmt.Sprintf("Draft: %s", strings.Join(names, ", "))
	case game.DesignOptionsEvent:
		return "Design phase begins"
	case game.DishCreatedEvent:
		return fmt.Sprintf("Dish created: %s", e.Dish.Name)
	default:
		return e.EventType()
	}
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
		line := fmt.Sprintf("%s %s (%s)", cursor, ing.Name, ing.Role)
		if d.cursor == i {
			line = selectedStyle.Render(line)
		}
		b.WriteString(line + "\n")
	}
	return paneStyle.Render(b.String())
}

func (d *draftMode) Status(m *model) string {
	return "up/down: move • enter/space: draft • q: quit"
}

// ---- Design Mode ----
type designMode struct {
	drafted  []ingredient.Ingredient
	cursor   int
	selected map[int]bool
	name     textinput.Model
	message  string
	dishes   []string
}

func (d *designMode) Init(m *model) tea.Cmd {
	d.selected = make(map[int]bool)
	d.name = textinput.New()
	d.name.Placeholder = "Dish name"
	d.name.Focus()
	d.dishes = []string{}
	return nil
}

func (d *designMode) Update(m *model, msg tea.Msg) (uiMode, tea.Cmd) {
	var cmd tea.Cmd
	d.name, cmd = d.name.Update(msg)
	switch msg := msg.(type) {
	case game.DishCreatedEvent:
		d.dishes = append(d.dishes, msg.Dish.Name)
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
			if len(d.dishes) >= 2 {
				d.message = "Maximum of 2 dishes reached"
				break
			}
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
		line := fmt.Sprintf("%s%s %s (%s)", cursor, mark, ing.Name, ing.Role)
		if d.cursor == i || d.selected[i] {
			line = selectedStyle.Render(line)
		}
		b.WriteString(line + "\n")
	}
	if len(d.dishes) > 0 {
		b.WriteString("\nDishes:\n")
		for _, name := range d.dishes {
			b.WriteString("- " + name + "\n")
		}
	}
	b.WriteString("\n" + d.name.View() + "\n")
	if d.message != "" {
		b.WriteString(d.message + "\n")
	}
	return paneStyle.Render(b.String())
}

func (d *designMode) Status(m *model) string {
	return "up/down: move • space: select • enter: create dish • f: finish • q: quit"
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
		b.WriteString(fmt.Sprintf("%s: %s -> ", r.Customer.Name, strings.Join(craving, ", ")))
		if r.Dish != nil {
			b.WriteString(r.Dish.Name)
		} else {
			b.WriteString("no dish")
		}
		b.WriteString("\n")
	}
	return paneStyle.Render(b.String())
}

func (s *serviceMode) Status(m *model) string {
	if s.finished {
		return "q: quit"
	}
	return ""
}

var (
	titleStyle    = lipgloss.NewStyle().Bold(true)
	paneStyle     = lipgloss.NewStyle().Padding(0, 1)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	statusStyle   = lipgloss.NewStyle().Padding(0, 1)
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
