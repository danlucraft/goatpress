package goatpress

import (
  "fmt"
)

type Match struct {
  Game      *Game
  Player1   Player
  Player2   Player
  played    bool
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
    response := MOVE_UNMADE
    i := 0
    move := MakePassMove()
    for i < 100 && response != MOVE_OK { // should have limit on number of invalid moves?
      move = thisPlayer.GetMove(m.Game.CurrentGameState())
      fmt.Printf("cand player: %s, move: %s\n", thisPlayer.Name(), move.ToString())
      response = m.Game.Move(move)
      i++
    }

    if response != MOVE_OK {
      move = MakePassMove()
      m.Game.Move(move)
    }
    colorMask := m.Game.ColorMask()
    colorString := colorMask.ToString()
    fmt.Printf("MOVE player: %s, move: %s, colors:%s\n", thisPlayer.Name(), move.ToString(), colorString)
  }
}

func (m *Match) Winner() int {
  return 0
}
