package goatpress

type Match struct {
  Game    *Game
  Player1 Player
  Player2 Player
  played  bool
}

func newMatch(gt *GameType, p1 Player, p2 Player) *Match {
  return &Match{gt.NewGame(), p1, p2, false}
}

func (m *Match) Play() {
  if m.played {
    return
  }
  m.played = true
}

func (m *Match) Winner() int {
  return 0
}
