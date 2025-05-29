package menu

import (
	"github.com/ascii-arcade/wish-template/internal/messages"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Term     string
	Width    int
	Height   int
	Renderer *lipgloss.Renderer

	isJoining bool
	gameCode  string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height, m.Width = msg.Height, msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return m, func() tea.Msg { return messages.SwitchToGame{} }
		case "j":
			m.isJoining = true
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	var content string
	if m.isJoining {
		content = "Enter the game code to join:\n\n"
	} else {
		content = "Welcome to the Game!\n\n"
		content += "Press 'n' to create a new game.\n"
		content += "Press 'j' to join an existing game.\n"
	}

	style := lipgloss.NewStyle().
		Width(m.Width).
		Height(m.Height)

	return style.Render(
		lipgloss.Place(
			m.Width,
			m.Height,
			lipgloss.Center,
			lipgloss.Center,
			content,
		),
	)
}
