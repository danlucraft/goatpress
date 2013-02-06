package goatpress

import (
  "testing"
)

func TestNewTournament(t *testing.T) {
  gameType := newGameType(5, DefaultWordSet)
  tourney := newTournament(*gameType)
  if tourney.Size() != 0 {
    t.Errorf("fresh tournament has more than one player!")
  }
  player1 := newInternalPlayer("Alice", newRandomFinder(DefaultWordSet))
  player2 := newInternalPlayer("Bob", newRandomFinder(DefaultWordSet))

  tourney.RegisterPlayer(player1)
  tourney.RegisterPlayer(player1)
  tourney.RegisterPlayer(player2)
  if tourney.Size() != 2 {
    t.Errorf("two players added but T size is not 2")
  }
  tourney.DeregisterPlayer(player2)
  if tourney.Size() != 1 {
    t.Errorf("deregister player didn't remove it")
  }

  tourney.PlayMatch()
  if len(tourney.Matches) != 0 {
    t.Errorf("can't play a match with one player")
  }

  tourney.RegisterPlayer(player2)
  tourney.PlayMatch()
  if len(tourney.Matches) != 1 {
    t.Errorf("T PlayMatch didn't play a match")
  }
  tourney.PlayMatch()
  tourney.PlayMatch()
  tourney.PlayMatch()
  tourney.PlayMatch()
  if len(tourney.Matches) != 5 {
    t.Errorf("T PlayMatch didn't play a match")
  }
}
