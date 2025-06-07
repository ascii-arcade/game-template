package games

func (s *Game) Count(player *Player) error {
	return s.withLock(func() error {
		player.incrementCount()
		return nil
	})
}
