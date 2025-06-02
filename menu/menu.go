package menu

import (
	"time"

	"github.com/ascii-arcade/wish-template/config"
	"github.com/ascii-arcade/wish-template/keys"
	"github.com/ascii-arcade/wish-template/language"
	"github.com/ascii-arcade/wish-template/messages"
	"github.com/ascii-arcade/wish-template/screen"
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

type doneMsg struct{}

type Model struct {
	Width              int
	Height             int
	screen             screen.Screen
	style              lipgloss.Style
	languagePreference *language.LanguagePreference

	errorCode     string
	gameCodeInput textinput.Model
}

func NewModel(width, height int, style lipgloss.Style, languagePreference *language.LanguagePreference) Model {
	ti := textinput.New()
	ti.Width = 9
	ti.CharLimit = 7

	m := Model{
		Width:              width,
		Height:             height,
		style:              style,
		languagePreference: languagePreference,

		gameCodeInput: ti,
	}

	m.screen = m.newSplashScreen()
	return m
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

func (m *Model) lang() *language.Language {
	return m.languagePreference.Lang
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height, m.Width = msg.Height, msg.Width
		return m, nil

	case messages.SwitchScreenMsg:
		m.screen = msg.Screen.WithModel(&m)
		return m, nil

	case tea.KeyMsg:
		if keys.ExitApplication.TriggeredBy(msg.String()) {
			return m, tea.Quit
		}
	}

	screenModel, cmd := m.screen.Update(msg)
	return screenModel.(*Model), cmd
}

func (m Model) View() string {
	if m.Width < config.MinimumWidth {
		return m.lang().Get("error.window_too_narrow")
	}
	if m.Height < config.MinimumHeight {
		return m.lang().Get("error.window_too_short")
	}

	return m.screen.View()
}

func (m *Model) setError(err string) {
	m.errorCode = err
}

func (m *Model) clearError() {
	m.errorCode = ""
}
