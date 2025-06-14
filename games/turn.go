package games

func (s *Game) NextTurn() {
	s.withLock(func() {
		if len(s.players) > s.CurrentTurnIndex+1 {
			s.CurrentTurnIndex++
		} else {
			s.CurrentTurnIndex = 0
		}
	})
}
