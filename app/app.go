package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"

	"github.com/ascii-arcade/wish-template/games"
	"github.com/ascii-arcade/wish-template/messages"
	"github.com/ascii-arcade/wish-template/screen"
)

type Model struct {
	player *games.Player
	game   *games.Game

	width  int
	height int
	screen screen.Screen
	style  lipgloss.Style

	error string
}

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	m := Model{
		width:  pty.Window.Width,
		height: pty.Window.Height,
		style:  bubbletea.MakeRenderer(s).NewStyle(),
	}

	m.screen = m.newSplashScreen()
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return doneMsg{}
		}),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height, m.width = msg.Height, msg.Width
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case messages.RefreshBoard:
		return m, func() tea.Msg {
			return messages.RefreshBoard(<-m.player.UpdateChan)
		}

	case messages.SwitchScreenMsg:
		m.screen = msg.Screen.WithModel(&m)
		return m, nil

		// case messages.NewGame:
		// 	err := m.newGame()
		// 	if err == nil {
		// 		m.active = m.board
		// 		m.board.Init()
		// 	}
		// 	return m, func() tea.Msg {
		// 		return messages.SwitchViewMsg{
		// 			Model: m.board,
		// 		}
		// 	}
		// case messages.JoinGame:
		// 	err := m.joinGame(msg.GameCode, false)
		// 	if err == nil {
		// 		m.active = m.board
		// 		m.board.Init()
		// 	}
		// 	return m, func() tea.Msg {
		// 		return messages.SwitchViewMsg{
		// 			Model: m.board,
		// 		}
		// 	}
	}

	activeScreenModel, cmd := m.screen.Update(msg)
	return activeScreenModel.(*Model), cmd
}

func (m Model) View() string {
	return m.screen.View()
}

func (m *Model) newGame() string {
	m.game = games.New()
	return m.game.Code
}

func (m *Model) joinGame(code string, isNew bool) error {
	game, err := games.GetOpenGame(code)
	if err != nil {
		return err
	}

	player, err := game.AddPlayer(isNew)
	if err != nil {
		return err
	}

	m.game = game
	m.player = player

	return nil
}

func (m *Model) setError(err string) {
	m.error = err
}

func (m *Model) clearError() {
	m.error = ""
}
