package games

func (s *Game) Count(pName string) {
	s.withLock(func() {
		if player, exists := s.getPlayer(pName); exists {
			player.incrementCount()
		}
	})
}
