package games

func (s *Game) Count(pName string) {
	s.withLock(func() {
		if player, exists := s.Players[pName]; exists {
			player.Count++
		}
	})
}
