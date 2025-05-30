package board

import (
	tea "github.com/charmbracelet/bubbletea"
)

type lobbyScreen struct {
}

func (l lobbyScreen) Update(m *Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "s":
		m.screen = tableScreen{}
	}

	return m, nil
}

func (l lobbyScreen) View(m *Model) string {
	return m.renderer.NewStyle().Render("hello world! press s to start!")
}
