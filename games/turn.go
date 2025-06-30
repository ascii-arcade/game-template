package games

import "github.com/ascii-arcade/game-template/messages"

func (s *Game) NextTurn() {
	s.withLock(func() {
		if len(s.players) > s.CurrentTurnIndex+1 {
			s.CurrentTurnIndex++
		} else {
			s.CurrentTurnIndex = 0
		}

		winner := s.GetWinner()
		if winner != nil {
			for _, player := range s.players {
				player.update(messages.WinnerScreen)
			}
			return
		}
	})
}
