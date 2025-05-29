package root

import (
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"

	"github.com/ascii-arcade/wish-template/internal/board"
	"github.com/ascii-arcade/wish-template/internal/game"
	generateRandom "github.com/ascii-arcade/wish-template/internal/generate_random"
	"github.com/ascii-arcade/wish-template/internal/menu"
	"github.com/ascii-arcade/wish-template/internal/messages"
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
	switch msg.(type) {
	case messages.SwitchToMenu:
		m.active = m.menu
	case messages.SwitchToGame:
		m.active = m.board
		m.board.Init()
	}

	var cmd tea.Cmd
	newModel, cmd := m.active.Update(msg)
	m.active = newModel
	return m, cmd
}

func (m rootModel) View() string {
	return m.active.View()
}

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	renderer := bubbletea.MakeRenderer(s)

	name := generateRandom.Name()

	m := rootModel{
		board: board.Model{
			Player:   name,
			Term:     pty.Term,
			Width:    pty.Window.Width,
			Height:   pty.Window.Height,
			Renderer: renderer,
		},
		menu: menu.Model{
			Term:     pty.Term,
			Width:    pty.Window.Width,
			Height:   pty.Window.Height,
			Renderer: renderer,
		},
	}
	m.active = m.menu

	err := m.newGame()
	if err != nil {
		log.Fatal("Could not create new game", "error", err)
	}

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m *rootModel) newGame() error {
	code := generateRandom.Code([]string{})
	game.Games[code] = game.NewGame()
	err := m.joinGame(code)

	if err != nil {
		return err
	}

	return nil
}

func (m *rootModel) joinGame(code string) error {
	updateCh := make(chan int)
	m.board.UpdateCh = updateCh
	m.board.GameCode = code

	state, exists := game.Games[code]
	if !exists {
		return errors.New("game does not exist")
	}

	state.AddClient(updateCh)
	state.Players[m.board.Player] = &game.Player{
		Name:      m.board.Player,
		TurnOrder: len(state.Players) + 1,
	}

	state.Refresh()

	return nil
}
