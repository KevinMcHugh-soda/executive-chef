package ui

import tea "github.com/charmbracelet/bubbletea"

type model struct{}

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
	return "Welcome to Executive Chef\n"
}

func Run() error {
	p := tea.NewProgram(model{})
	_, err := p.Run()
	return err
}
