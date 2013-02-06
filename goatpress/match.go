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
  players  := [2]Player {m.Player1, m.Player2}
  messages := [2]int {MSG_NONE, MSG_NONE}
  m.Player1.NewGame(m.Game.CurrentGameState())
  m.Player2.NewGame(m.Game.CurrentGameState())

  for m.Game.WhoseMove() != 0 {
    playerIx := m.Game.WhoseMove() - 1
    thisPlayer := players[playerIx]

    response := MOVE_UNMADE
    i := 0
    move := MakePassMove()
    for i < 100 && response != MOVE_OK { // should have limit on number of invalid moves?
      move = thisPlayer.GetMove(messages[playerIx], "", m.Game.CurrentGameState())
      response = m.Game.Move(move)
      switch response {
      case MOVE_OK:
        messages[playerIx] = MSG_OPPONENT_MOVE
      case MOVE_ALREADY_PLAYED:
        messages[playerIx] = MSG_BAD_MOVE_ALREADY
      case MOVE_PREFIX_WORD:
        messages[playerIx] = MSG_BAD_MOVE_PREFIX
      case MOVE_INVALID_WORD:
        messages[playerIx] = MSG_BAD_MOVE_UNKNOWN
      case MOVE_TOO_SHORT:
        messages[playerIx] = MSG_BAD_MOVE_TOO_SHORT
      }
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
  colorMask := m.Game.ColorMask()
  score1 := colorMask.Score(1)
  score2 := colorMask.Score(2)
  if score1 > score2 {
    return 1
  } else if score2 > score1 {
    return 2
  }
  return 0
}

type MatchMarshaller struct {
  Player1 string
  Player2 string
  Game    GameMarshaller
}

func (m *Match) Marshaller() MatchMarshaller {
  mm := MatchMarshaller{m.Player1.Name(), m.Player2.Name(), m.Game.Marshaller()}
  return mm
}

func dummyForFmt() {
  fmt.Printf("adsf")
}

