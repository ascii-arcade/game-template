package games

func (s *Game) NextTurn() {
	s.withLock(func() {
		if len(s.players) > s.currentTurnIndex+1 {
			s.currentTurnIndex++
		} else {
			s.currentTurnIndex = 0
		}
	})
}
