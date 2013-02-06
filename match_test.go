package goatpress

import (
  "testing"
)

func TestMatchBetweenPassesPlayers(t *testing.T) {
  gameType := newGameType(5, DefaultWordSet)
  player1 := newInternalPlayer("test1", newPassFinder())
  player2 := newInternalPlayer("test2", newPassFinder())
  match := newMatch(gameType, player1, player2)
  match.Play()
  winner := match.Winner()
  if winner != 0 {
    t.Errorf("a game of passing wasn't a draw")
  }
}
