package root

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"

	"github.com/ascii-arcade/wish-template/internal/board"
	"github.com/ascii-arcade/wish-template/internal/game"
	"github.com/ascii-arcade/wish-template/internal/menu"
	"github.com/ascii-arcade/wish-template/internal/messages"
)

type rootModel struct {
	active tea.Model
	sess   ssh.Session
}

func (m rootModel) Init() tea.Cmd {
	return m.active.Init()
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.SwitchViewMsg:
		m.active = msg.NewModel
		initcmd := msg.NewModel.Init()
		return m, initcmd
	case messages.NewGame:
		model, _ := m.joinGame(game.NewGame().Code)
		return m, func() tea.Msg {
			return messages.SwitchViewMsg{
				NewModel: model,
			}
		}
	case messages.JoinGame:
		model, _ := m.joinGame(msg.GameCode)
		return m, func() tea.Msg {
			return messages.SwitchViewMsg{
				NewModel: model,
			}
		}
	}

	updateModel, cmd := m.active.Update(msg)
	m.active = updateModel
	return m, cmd
}

func (m rootModel) View() string {
	return m.active.View()
}

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	renderer := bubbletea.MakeRenderer(s)

	m := rootModel{
		active: menu.NewModel(pty.Term, pty.Window.Width, pty.Window.Height, renderer),
		sess:   s,
	}

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m *rootModel) joinGame(code string) (tea.Model, error) {
	state, exists := game.GetGame(code)
	if !exists {
		return nil, errors.New("game does not exist")
	}
	player := state.AddPlayer()

	pty, _, _ := m.sess.Pty()
	renderer := bubbletea.MakeRenderer(m.sess)

	state.Refresh()

	return board.New(pty.Term, renderer, state, player), nil
}
