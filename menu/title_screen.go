package menu

import (
	"fmt"
	"strings"

	"github.com/ascii-arcade/game-template/games"
	"github.com/ascii-arcade/game-template/keys"
	"github.com/ascii-arcade/game-template/language"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type titleScreen struct {
	model *Model
	style lipgloss.Style
}

func (m *Model) newTitleScreen() *titleScreen {
	return &titleScreen{
		model: m,
		style: m.style,
	}
}

func (s *titleScreen) Update(msg tea.Msg) (any, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.model.height, s.model.width = msg.Height, msg.Width
		return s.model, nil

	case tea.KeyMsg:
		switch {
		case keys.MenuEnglish.TriggeredBy(msg.String()):
			s.model.player.LanguagePreference.SetLanguage("EN")
		case keys.MenuSpanish.TriggeredBy(msg.String()):
			s.model.player.LanguagePreference.SetLanguage("ES")
		case keys.MenuStartNewGame.TriggeredBy(msg.String()):
			newGame := games.New()
			if err := s.model.joinGame(newGame.Code, true); err != nil {
				s.model.setError(err.Error())
				return s.model, nil
			}

			return s.model, func() tea.Msg { return SwitchToBoardMsg{Game: newGame} }
		case keys.MenuJoinGame.TriggeredBy(msg.String()):
			s.model.screen = s.model.newJoinScreen()
			return s.model, nil
		}
	}

	return s.model, nil
}

func (s *titleScreen) View() string {
	var content strings.Builder
	content.WriteString(s.model.lang().Get("menu", "welcome") + "\n\n")
	content.WriteString(fmt.Sprintf(s.model.lang().Get("menu", "press_to_create"), keys.MenuStartNewGame.String(s.style)) + "\n")
	content.WriteString(fmt.Sprintf(s.model.lang().Get("menu", "press_to_join"), keys.MenuJoinGame.String(s.style)) + "\n")
	content.WriteString("\n\n")

	if s.model.lang() == language.Languages["EN"] {
		content.WriteString(fmt.Sprintf(language.Languages["ES"].Get("menu", "choose_language"), keys.MenuSpanish.String(s.style)))
	} else if s.model.lang() == language.Languages["ES"] {
		content.WriteString(fmt.Sprintf(language.Languages["EN"].Get("menu", "choose_language"), keys.MenuEnglish.String(s.style)))
	}

	return content.String()
}
