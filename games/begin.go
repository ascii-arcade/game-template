package games

func (s *Game) Begin() {
	s.withLock(func() {
		s.inProgress = true
	})
}
