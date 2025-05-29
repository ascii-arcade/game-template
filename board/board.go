package board

import (
	"fmt"
	"log"

	"github.com/ascii-arcade/wish-template/game"
	"github.com/ascii-arcade/wish-template/messages"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Player   string
	Term     string
	Width    int
	Height   int
	Renderer *lipgloss.Renderer

	GameCode string
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
	state := m.gameState()
	state.LockState()
	defer func() {
		state.UnlockState()
		state.Refresh()
	}()

	switch msg.String() {
	case "a":
		state.Players[m.Player].Count++
	case "q", "ctrl+c":
		state.RemoveClient(m.UpdateCh, m.Player)
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) View() string {
	counts := ""
	for _, p := range m.gameState().OrderedPlayers() {
		counts += fmt.Sprintf("%s: %d\n", p.Name, p.Count)
	}

	return m.Renderer.NewStyle().Render(fmt.Sprintf("You are %s", m.Player)) +
		"\n\n" + counts +
		"\n\n'" + m.GameCode + "'" +
		"\n\n" + m.Renderer.NewStyle().Render("Press 'q' to quit")
}

func (m *Model) gameState() *game.Game {
	state, exists := game.Games[m.GameCode]
	if !exists {
		log.Fatal("Game does not exist", "code", m.GameCode)
	}
	return state
}

func waitForRefreshSignal(ch chan int) tea.Cmd {
	return func() tea.Msg {
		v := <-ch
		return messages.RefreshGame(v)
	}
}
