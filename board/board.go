package board

import (
	"github.com/ascii-arcade/wish-template/games"
	"github.com/ascii-arcade/wish-template/messages"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen interface {
	setModel(*Model)
	update(tea.Msg) (tea.Model, tea.Cmd)
	view() string
}

type Model struct {
	Width  int
	Height int
	screen screen
	style  lipgloss.Style

	Player *games.Player
	Game   *games.Game
}

func NewModel(width, height int, style lipgloss.Style) Model {
	m := Model{
		Width:  width,
		Height: height,
		style:  style,
	}

	m.screen = m.newTableScreen()
	return m
}

func (m Model) Init() tea.Cmd {
	return waitForRefreshSignal(m.Player.UpdateChan)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height, m.Width = msg.Height, msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.Game.RemovePlayer(m.Player.Name)
			return m, tea.Quit
		}

	case messages.RefreshBoard:
		return m, waitForRefreshSignal(m.Player.UpdateChan)

	default:
		return m.activeScreen().update(msg)
	}

	return m, nil
}

func (m Model) View() string {
	return m.activeScreen().view()
}

func (m *Model) activeScreen() screen {
	if m.Game.InProgress() {
		return m.newTableScreen()
	} else {
		return m.newLobbyScreen()
	}
}

func waitForRefreshSignal(ch chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return messages.RefreshBoard(<-ch)
	}
}
