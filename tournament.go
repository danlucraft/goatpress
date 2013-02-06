package goatpress

import (
  "math/rand"
  "fmt"
)

type Tournament struct {
  GameType GameType
  Players  map[string]Player
  Matches  []Match
}

func newTournament(gt GameType) *Tournament {
  return &Tournament{gt, make(map[string]Player), make([]Match, 0)}
}

func (t *Tournament) RegisterPlayer(p Player) {
  if _, present := t.Players[p.Name()]; present {
    return
  }
  t.Players[p.Name()] = p
}

func (t *Tournament) DeregisterPlayer(p Player) {
  delete(t.Players, p.Name())
}

func (t *Tournament) PlayMatch() {
  if len(t.Players) > 1 {
    player1 := t.RandomPlayer()
    player2 := t.RandomPlayer()
    for player1.Name() == player2.Name() {
      player2 = t.RandomPlayer()
    }
    fmt.Printf("Match between %s and %s... \n", player1.Name(), player2.Name())
    match := newMatch(&t.GameType, player1, player2)
    t.Matches = append(t.Matches, *match)
    match.Play()
    winnerIx := match.Winner()
    if winnerIx == 0 {
      fmt.Printf("DRAW\n")
    } else if winnerIx == 1 {
      fmt.Printf("Winner was %s\n", player1.Name())
    } else if winnerIx == 2 {
      fmt.Printf("Winner was %s\n", player2.Name())
    }
  }
}

func (t *Tournament) RandomPlayer() Player {
  i := rand.Intn(len(t.Players))
  j := 0
  for _, player := range t.Players {
    if j == i {
      return player
    }
    j++
  }
  panic("your indexing is bad")
}

func (t *Tournament) Size() int {
  return len(t.Players)
}

