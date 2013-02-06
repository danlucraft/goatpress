package goatpress

import (
  "testing"
)

func TestNewTournament(t *testing.T) {
  gameType := newGameType(5, testWordSet())
  tourney := newTournament(*gameType)
  if tourney.Size() != 0 {
    t.Errorf("fresh tournament has more than one player!")
  }
}
