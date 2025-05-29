package root

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"

	"github.com/ascii-arcade/wish-template/board"
	"github.com/ascii-arcade/wish-template/game"
	"github.com/ascii-arcade/wish-template/menu"
	"github.com/ascii-arcade/wish-template/messages"
)

type rootModel struct {
	active tea.Model
	menu   menu.Model
	board  board.Model
}

func (m rootModel) Init() tea.Cmd {
	return m.active.Init()
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		err := m.joinGame(msg.GameCode)
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

func (m rootModel) View() string {
	return m.active.View()
}

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	renderer := bubbletea.MakeRenderer(s)

	m := rootModel{
		board: board.Model{
			Term:     pty.Term,
			Width:    pty.Window.Width,
			Height:   pty.Window.Height,
			Renderer: renderer,
		},
		menu: menu.NewModel(pty.Term, pty.Window.Width, pty.Window.Height, renderer),
	}
	m.active = m.menu

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m *rootModel) newGame() error {
	newGame := game.New()
	game.Games[newGame.Code] = newGame
	m.board.Game = newGame
	return m.joinGame(newGame.Code)
}

func (m *rootModel) joinGame(code string) error {
	updateCh := make(chan int)
	m.board.UpdateCh = updateCh

	g, ok := game.Get(code)
	if !ok {
		return errors.New("game does not exist")
	}
	m.board.Game = g

	player := g.AddPlayer(updateCh)
	m.board.Player = player

	g.Refresh()

	return nil
}
