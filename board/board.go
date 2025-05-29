package board

import (
	"fmt"
	"log"

	"github.com/ascii-arcade/wish-template/games"
	"github.com/ascii-arcade/wish-template/messages"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Term     string
	Width    int
	Height   int
	Renderer *lipgloss.Renderer

	Player   *games.Player
	Game     *games.Game
	UpdateCh chan int
}

func (m Model) Init() tea.Cmd {
	return waitForRefreshSignal(m.UpdateCh)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height, m.Width = msg.Height, msg.Width

	case tea.KeyMsg:
		return m.handleKey(msg)

	case messages.RefreshGame:
		return m, waitForRefreshSignal(m.UpdateCh)
	}

	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	game := m.gameState()
	defer func() {
		game.Refresh()
	}()

	switch msg.String() {
	case "a":
		game.Count(m.Player.Name)
	case "q", "ctrl+c":
		game.RemovePlayer(m.Player.Name)
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) View() string {
	counts := ""
	for _, p := range m.gameState().OrderedPlayers() {
		counts += fmt.Sprintf("%s: %d\n", p.Name, p.Count)
	}

	return m.Renderer.NewStyle().Render(fmt.Sprintf("You are %s", m.Player.Name)) +
		"\n\n" + counts +
		"\n\n'" + m.Game.Code + "'" +
		"\n\n" + m.Renderer.NewStyle().Render("Press 'q' to quit")
}

func (m *Model) gameState() *games.Game {
	game, exists := games.Games[m.Game.Code]
	if !exists {
		log.Fatal("Game does not exist", "code", m.Game.Code)
	}
	return game
}

func waitForRefreshSignal(ch chan int) tea.Cmd {
	return func() tea.Msg {
		v := <-ch
		return messages.RefreshGame(v)
	}
}
