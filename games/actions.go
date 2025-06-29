package games

func (s *Game) EarnPoint(player *Player) {
	s.withLock(func() {
		player.incrementPoints()
	})
}
