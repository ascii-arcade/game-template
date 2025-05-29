package board

import (
	"fmt"

	"github.com/ascii-arcade/wish-template/internal/game"
	"github.com/ascii-arcade/wish-template/internal/messages"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Term     string
	Renderer *lipgloss.Renderer
	Player   *game.Player
	Game     *game.Game

	width  int
	height int
}

func New(term string, renderer *lipgloss.Renderer, game *game.Game, player *game.Player) Model {
	return Model{
		Term:     term,
		Renderer: renderer,
		Game:     game,
		Player:   player,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		waitForRefreshSignal(m.Player.RefreshCh),
		tea.WindowSize(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height, m.width = msg.Height, msg.Width

	case tea.KeyMsg:
		return m.handleKey(msg)

	case messages.RefreshGame:
		return m, waitForRefreshSignal(m.Player.RefreshCh)
	}

	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.Game.LockState()
	defer func() {
		m.Game.UnlockState()
		m.Game.Refresh()
	}()

	switch msg.String() {
	case "a":
		m.Player.Count++
	case "q", "ctrl+c":
		m.Game.RemovePlayer(m.Player)
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) View() string {
	counts := ""
	for _, p := range m.Game.OrderedPlayers() {
		counts += fmt.Sprintf("%s: %d\n", p.Name, p.Count)
	}

	return m.Renderer.NewStyle().Render(fmt.Sprintf("You are %s", m.Player)) +
		"\n\n" + counts +
		"\n\n'" + m.Game.Code + "'" +
		"\n\n" + m.Renderer.NewStyle().Render("Press 'q' to quit")
}

func waitForRefreshSignal(ch chan messages.RefreshGame) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
