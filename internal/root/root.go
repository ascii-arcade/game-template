package root

import (
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"

	"github.com/ascii-arcade/wish/internal/game"
	gameView "github.com/ascii-arcade/wish/internal/game_view"
	generateRandom "github.com/ascii-arcade/wish/internal/generate_random"
	"github.com/ascii-arcade/wish/internal/menu"
	"github.com/ascii-arcade/wish/internal/messages"
)

type rootModel struct {
	active tea.Model
	menu   menu.Model
	game   gameView.Model
}

func (m rootModel) Init() tea.Cmd {
	return m.active.Init()
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case messages.SwitchToMenu:
		m.active = m.menu
	case messages.SwitchToGame:
		m.active = m.game
		m.game.Init()
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
		game: gameView.Model{
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
	m.game.UpdateCh = updateCh
	m.game.GameCode = code

	state, exists := game.Games[code]
	if !exists {
		return errors.New("game does not exist")
	}

	state.AddClient(updateCh)
	state.Players[m.game.Player] = &game.Player{
		Name:      m.game.Player,
		TurnOrder: len(state.Players) + 1,
	}

	state.Refresh()

	return nil
}
