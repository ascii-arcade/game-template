package board

import (
	"github.com/ascii-arcade/wish-template/games"
	"github.com/ascii-arcade/wish-template/messages"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen interface {
	setModel(*Model)
	update(tea.KeyMsg) (tea.Model, tea.Cmd)
	view() string
}

type Model struct {
	Width    int
	Height   int
	renderer *lipgloss.Renderer
	screen   screen

	Player *games.Player
	Game   *games.Game
}

func NewModel(width, height int, renderer *lipgloss.Renderer) Model {
	return Model{
		Width:    width,
		Height:   height,
		renderer: renderer,
		screen:   &tableScreen{},
	}
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
		default:
			return m.activeScreen().update(msg)
		}

	case messages.RefreshGame:
		return m, waitForRefreshSignal(m.Player.UpdateChan)
	}

	return m, nil
}

func (m Model) View() string {
	return m.activeScreen().view()
}

func (m *Model) activeScreen() screen {
	if m.Game.InProgress() {
		m.screen.setModel(m)
		return m.screen
	} else {
		return &lobbyScreen{model: m}
	}
}

func waitForRefreshSignal(ch chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return messages.RefreshGame(<-ch)
	}
}
