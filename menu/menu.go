package menu

import (
	"github.com/ascii-arcade/wish-template/messages"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Term     string
	Width    int
	Height   int
	Renderer *lipgloss.Renderer

	isJoining     bool
	gameCodeInput textinput.Model
}

func NewModel(term string, width, height int, renderer *lipgloss.Renderer) Model {
	ti := textinput.New()
	ti.Width = 9
	ti.CharLimit = 7

	return Model{
		Term:     term,
		Width:    width,
		Height:   height,
		Renderer: renderer,

		isJoining:     false,
		gameCodeInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.isJoining && m.gameCodeInput.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyEsc:
				m.isJoining = false
				return m, nil
			case tea.KeyEnter:
				if len(m.gameCodeInput.Value()) == 7 {
					return m, func() tea.Msg { return messages.JoinGame{GameCode: m.gameCodeInput.Value()} }
				}
			}

			m.gameCodeInput, cmd = m.gameCodeInput.Update(msg)
			return m, cmd
		}
	} else {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.Height, m.Width = msg.Height, msg.Width

		case tea.KeyMsg:
			switch msg.String() {
			case "n":
				return m, func() tea.Msg { return messages.NewGame{} }
			case "j":
				m.isJoining = true
				m.gameCodeInput.Focus()
				m.gameCodeInput.SetValue("")
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var content string
	if m.isJoining {
		content = "Enter the game code to join:\n\n" + m.gameCodeInput.View()
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
