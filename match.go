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
  players := [3]Player {nil, m.Player1, m.Player2}
  for m.Game.WhoseMove() != 0 {
    thisPlayer := players[m.Game.WhoseMove()]
    move := thisPlayer.GetMove(m.Game.CurrentGameState())
    response := MOVE_UNMADE
    for response != MOVE_OK { // should have limit on number of invalid moves?
      response = m.Game.Move(move)
    }
  }
}

func (m *Match) Winner() int {
  return 0
}
