package menu

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const logo = `++------------------------------------------------------------------------------++
++------------------------------------------------------------------------------++
||                                                                              ||
||                                                                              ||
||      _    ____   ____ ___ ___        _    ____   ____    _    ____  _____    ||
||     / \  / ___| / ___|_ _|_ _|      / \  |  _ \ / ___|  / \  |  _ \| ____|   ||
||    / _ \ \___ \| |    | | | |_____ / _ \ | |_) | |     / _ \ | | | |  _|     ||
||   / ___ \ ___) | |___ | | | |_____/ ___ \|  _ <| |___ / ___ \| |_| | |___    ||
||  /_/   \_\____/ \____|___|___|   /_/   \_\_| \_\\____/_/   \_\____/|_____|   ||
||                                                                              ||
||                                                                              ||
||                                                                              ||
++------------------------------------------------------------------------------++
++------------------------------------------------------------------------------++`

type screen interface {
	setModel(*Model)
	update(tea.KeyMsg) (tea.Model, tea.Cmd)
	view() string
}

type doneMsg struct{}

type Model struct {
	Width    int
	Height   int
	renderer *lipgloss.Renderer
	screen   screen

	isSplashing   bool
	gameCodeInput textinput.Model
}

func NewModel(width, height int, renderer *lipgloss.Renderer) Model {
	ti := textinput.New()
	ti.Width = 9
	ti.CharLimit = 7

	return Model{
		Width:    width,
		Height:   height,
		renderer: renderer,
		screen:   &splashScreen{},

		gameCodeInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return doneMsg{}
		}),
		tea.WindowSize(),
		textinput.Blink,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height, m.Width = msg.Height, msg.Width

	case doneMsg:
		m.screen = &optionScreen{}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			return m.activeScreen().update(msg)
		}
	}

	return m, nil
}

func (m Model) View() string {
	return m.activeScreen().view()
}

func (m *Model) activeScreen() screen {
	m.screen.setModel(m)
	return m.screen
}
