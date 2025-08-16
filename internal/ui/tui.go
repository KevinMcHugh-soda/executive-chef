package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"executive-chef/internal/dish"
	"executive-chef/internal/game"
	"executive-chef/internal/ingredient"
)

const logWidth = 30

// uiMode represents a UI mode for a game phase.
type uiMode interface {
	Init(*model) tea.Cmd
	Update(*model, tea.Msg) (uiMode, tea.Cmd)
	View(*model) string
	// Status returns a short help prompt for the status bar.
	Status(*model) string
}

type model struct {
	actions     chan<- game.Action
	mode        uiMode
	events      []string
	vp          viewport.Model
	turn        int
	phase       game.Phase
	dishes      []dish.Dish
	ingredients []ingredient.Ingredient
	message     string
	width       int
	money       int
}

func initialModel(actions chan<- game.Action) *model {
	m := &model{actions: actions}
	m.mode = &draftMode{}
	m.vp = viewport.New(logWidth-2, 7)
	return m
}

func (m *model) Init() tea.Cmd {
	return m.mode.Init(m)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if e, ok := msg.(game.Event); ok {
		if str := eventString(e); str != "" {
			m.events = append(m.events, str)
			m.vp.SetContent(strings.Join(m.events, "\n"))
			m.vp.GotoBottom()
		}
		switch ev := e.(type) {
		case game.PhaseEvent:
			m.turn = ev.Turn
			m.phase = ev.Phase
		case game.IngredientDraftedEvent:
			m.ingredients = append(m.ingredients, ev.Ingredient)
		case game.DishCreatedEvent:
			m.dishes = append(m.dishes, ev.Dish)
		case game.ServiceEndEvent:
			m.ingredients = nil
		}
		if pay, ok := e.(game.ServiceResultEvent); ok {
			m.money = pay.Money
		}
	}

	if wm, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = wm.Width
	}

	var vpCmd tea.Cmd
	m.vp, vpCmd = m.vp.Update(msg)
	m.vp.Width = logWidth - 2

	newMode, modeCmd := m.mode.Update(m, msg)
	if newMode != nil {
		m.mode = newMode
		initCmd := m.mode.Init(m)
		return m, tea.Batch(modeCmd, initCmd, vpCmd)
	}

	return m, tea.Batch(vpCmd, modeCmd)
}

func (m *model) View() string {
	main := m.mode.View(m)

	var infoBuilder strings.Builder
	infoBuilder.WriteString(titleStyle.Render("Game Info") + "\n")
	infoBuilder.WriteString(
		fmt.Sprintf(
			"Turn: %d\nPhase: %s\nMoney: $%d\n",
			m.turn, m.phase, m.money,
		),
	)
	infoBuilder.WriteString("Dishes:\n")
	if len(m.dishes) == 0 {
		infoBuilder.WriteString("  (none)\n")
	} else {
		for _, d := range m.dishes {
			name := d.Name
			if !m.hasIngredients(d) {
				name = missingStyle.Render(name)
			}
			infoBuilder.WriteString("- " + name + "\n")
		}
	}
	info := paneStyle.Render(infoBuilder.String())
	logView := paneStyle.Render(titleStyle.Render("Events") + "\n" + m.vp.View())

	mainWidth := m.width - logWidth
	if mainWidth < 0 {
		mainWidth = 0
	}
	main = lipgloss.NewStyle().Width(mainWidth).Render(main)

	content := lipgloss.JoinHorizontal(lipgloss.Top, main, logView)
	status := statusStyle.Render(m.mode.Status(m))
	message := messageStyle.Render(m.message)
	return lipgloss.JoinVertical(lipgloss.Left, info, content, status, message)
}

func (m *model) hasIngredients(d dish.Dish) bool {
	for _, need := range d.Ingredients {
		found := false
		for _, have := range m.ingredients {
			if have == need {
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

func eventString(e game.Event) string {
	switch e := e.(type) {
	case game.PhaseEvent:
		return fmt.Sprintf("Turn %d: %s phase", e.Turn, e.Phase)
	case game.DraftOptionsEvent:
		if len(e.Reveal) == 10 {
			return "Draft phase started"
		}
		return ""
	case game.IngredientDraftedEvent:
		return fmt.Sprintf("Ingredient drafted: %s", e.Ingredient.Name)
	case game.DesignOptionsEvent:
		return "Design phase begins"
	case game.DishCreatedEvent:
		return fmt.Sprintf("Dish created: %s", e.Dish.Name)
	case game.ServiceResultEvent:
		dishName := "no dish"
		if e.Dish != nil {
			dishName = e.Dish.Name
		}
		if e.Payment > 0 {
			return fmt.Sprintf("%s served %s for $%d", e.Customer.Name, dishName, e.Payment)
		}
		return fmt.Sprintf("%s was not served", e.Customer.Name)
	default:
		return e.EventType()
	}
}

// ---- Draft Mode ----
type draftMode struct {
	draft  []ingredient.Ingredient
	cursor int
}

func (d *draftMode) Init(m *model) tea.Cmd {
	m.message = ""
	return nil
}

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
	drafted   []ingredient.Ingredient
	cursor    int
	selected  map[int]bool
	name      textinput.Model
	dishes    []string
	selecting bool
	confirm   bool
}

func (d *designMode) Init(m *model) tea.Cmd {
	d.selected = make(map[int]bool)
	d.name = textinput.New()
	d.name.Placeholder = "Dish name"
	d.name.Blur()
	d.dishes = []string{}
	d.selecting = true
	d.confirm = false
	m.message = ""
	return nil
}

func (d *designMode) Update(m *model, msg tea.Msg) (uiMode, tea.Cmd) {
	var cmd tea.Cmd
	if !d.selecting {
		d.name, cmd = d.name.Update(msg)
	}
	switch msg := msg.(type) {
	case game.DishCreatedEvent:
		d.dishes = append(d.dishes, msg.Dish.Name)
		m.message = fmt.Sprintf("Added dish '%s'!", msg.Dish.Name)
		d.name.SetValue("")
		d.selected = make(map[int]bool)
		d.confirm = false
	case game.ServiceResultEvent:
		return &serviceMode{current: &msg}, nil
	case tea.KeyMsg:
		if msg.String() != "enter" {
			d.confirm = false
			m.message = ""
		}
		switch msg.String() {
		case "ctrl+c", "q":
			m.actions <- game.FinishDesignAction{}
			return nil, tea.Quit
		case "up", "k":
			if d.selecting && d.cursor > 0 {
				d.cursor--
			}
		case "down", "j":
			if d.selecting && d.cursor < len(d.drafted)-1 {
				d.cursor++
			}
		case "enter":
			if d.selecting {
				if d.selected[d.cursor] {
					delete(d.selected, d.cursor)
				} else {
					d.selected[d.cursor] = true
				}
			} else {
				if !d.confirm {
					d.confirm = true
					m.message = "press enter again to confirm"
				} else {
					name := strings.TrimSpace(d.name.Value())
					if len(m.dishes) >= 10 {
						m.message = "Maximum of 10 dishes reached"
						break
					}
					if len(d.dishes) >= 2 {
						m.message = "Maximum of 2 dishes this turn reached"
						break
					}
					if len(d.selected) == 0 {
						m.message = "select at least one ingredient to create a dish!"
					} else if name != "" {
						var indices []int
						for i := range d.drafted {
							if d.selected[i] {
								indices = append(indices, i)
							}
						}
						m.actions <- game.CreateDishAction{Name: name, Indices: indices}
						m.message = ""
					} else {
						m.message = ""
					}
					d.confirm = false
				}
			}
		case "tab":
			d.confirm = false
			m.message = ""
			if d.selecting {
				d.selecting = false
				d.name.Focus()
			} else {
				d.selecting = true
				d.name.Blur()
			}
		case "f", "F":
			if d.selecting {
				m.message = ""
				m.actions <- game.FinishDesignAction{}
				return nil, nil
			}
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
	return paneStyle.Render(b.String())
}

func (d *designMode) Status(m *model) string {
	return "up/down: move • enter: select • tab: name • enter x2: create dish • f: finish • q: quit"
}

// ---- Service Mode ----
type serviceMode struct {
	current  *game.ServiceResultEvent
	finished bool
}

func (s *serviceMode) Init(m *model) tea.Cmd {
	m.message = ""
	return nil
}

func (s *serviceMode) Update(m *model, msg tea.Msg) (uiMode, tea.Cmd) {
	switch msg := msg.(type) {
	case game.ServiceResultEvent:
		s.current = &msg
	case game.ServiceEndEvent:
		s.finished = true
	case game.DraftOptionsEvent:
		return &draftMode{draft: msg.Reveal}, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return nil, tea.Quit
		case "enter":
			m.actions <- game.ContinueAction{}
		}
	}
	return nil, nil
}

func (s *serviceMode) View(m *model) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Service") + "\n")
	if s.current != nil {
		var craving []string
		if len(s.current.Customer.Cravings) > 0 {
			for _, ing := range s.current.Customer.Cravings[0].Ingredients {
				craving = append(craving, ing.Name)
			}
		}
		b.WriteString(fmt.Sprintf("%s: %s -> ", s.current.Customer.Name, strings.Join(craving, ", ")))
		if s.current.Dish != nil {
			b.WriteString(servedStyle.Render(s.current.Dish.Name))
		} else {
			b.WriteString(missingStyle.Render("no dish"))
		}
		if s.current.Payment > 0 {
			b.WriteString(fmt.Sprintf(" ($%d)", s.current.Payment))
		}
		b.WriteString("\n")
	}
	return paneStyle.Render(b.String())
}

func (s *serviceMode) Status(m *model) string {
	if s.finished {
		return "enter: next turn • q: quit"
	}
	return "enter: next customer • q: quit"
}

var (
	titleStyle    = lipgloss.NewStyle().Bold(true)
	paneStyle     = lipgloss.NewStyle().Padding(0, 1)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	missingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	servedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	statusStyle   = lipgloss.NewStyle().Padding(0, 1)
	messageStyle  = lipgloss.NewStyle().Padding(0, 1)
	logStyle      = paneStyle.Copy().Width(logWidth)
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
