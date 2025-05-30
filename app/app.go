package app

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"

	"github.com/ascii-arcade/wish-template/board"
	"github.com/ascii-arcade/wish-template/games"
	"github.com/ascii-arcade/wish-template/menu"
	"github.com/ascii-arcade/wish-template/messages"
)

type Model struct {
	active tea.Model
	menu   menu.Model
	board  board.Model
}

func (m Model) Init() tea.Cmd {
	return m.active.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.SwitchViewMsg:
		m.active = msg.Model
		initcmd := m.active.Init()
		return m, initcmd
	case messages.NewGame:
		err := m.newGame()
		if err == nil {
			m.active = m.board
			m.board.Init()
		}
		return m, func() tea.Msg {
			return messages.SwitchViewMsg{
				Model: m.board,
			}
		}
	case messages.JoinGame:
		err := m.joinGame(msg.GameCode, false)
		if err == nil {
			m.active = m.board
			m.board.Init()
		}
		return m, func() tea.Msg {
			return messages.SwitchViewMsg{
				Model: m.board,
			}
		}
	}

	var cmd tea.Cmd
	m.active, cmd = m.active.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.active.View()
}

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	renderer := bubbletea.MakeRenderer(s)

	m := Model{
		board: board.NewModel(pty.Window.Width, pty.Window.Height, renderer),
		menu:  menu.NewModel(pty.Window.Width, pty.Window.Height, renderer),
	}
	m.active = m.menu

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m *Model) newGame() error {
	newGame := games.New()
	m.board.Game = newGame
	return m.joinGame(newGame.Code, true)
}

func (m *Model) joinGame(code string, isNew bool) error {
	game, ok := games.Get(code)
	if !ok {
		return errors.New("game does not exist")
	}
	m.board.Game = game

	player := game.AddPlayer(isNew)
	m.board.Player = player

	return nil
}
