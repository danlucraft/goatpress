package goatpress

import (
  "math/rand"
  "fmt"
  "encoding/json"
  "os"
  "strconv"
)

type Tournament struct {
  GameType GameType
  Players  map[string]Player
  Matches  []Match
  MatchLog *os.File
  Scores   Scores
  AllPlayerNames map[string]bool
}

type Scores struct {
  Games         map[string]int
  Moves         map[string]int
  Wins          map[string]int
  Losses        map[string]int
}

func newTournament(gt GameType) *Tournament {
  g := make(map[string]int)
  m := make(map[string]int)
  w := make(map[string]int)
  l := make(map[string]int)
  tm := Scores{g,m,w,l}
  return newTournamentWithScores(gt, tm)
}

func newTournamentWithScores(gt GameType, scores Scores) *Tournament {
  return &Tournament{gt, make(map[string]Player), make([]Match, 0), nil, scores, make(map[string]bool)}
}

func (t *Tournament) RegisterPlayer(p Player) {
  if _, present := t.Players[p.Name()]; present {
    return
  }
  t.AllPlayerNames[p.Name()] = true
  t.Players[p.Name()] = p
}

func (t *Tournament) DeregisterPlayer(name string) {
  delete(t.Players, name)
}

func (t *Tournament) PlayerList() string {
  s := ""
  for name, _ := range t.Players {
    s += name + " ("
    s += strconv.Itoa(t.ScoreFor(name))
    s += "), "
  }
  return s
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

    name1 := player1.Name()
    name2 := player2.Name()
    t.Scores.Games[name1] += 1
    t.Scores.Games[name2] += 1
    t.Scores.Moves[name1] += len(match.Game.Moves)
    t.Scores.Moves[name2] += len(match.Game.Moves)
    winnerIx := match.Winner()
    if winnerIx == 0 {
      fmt.Printf("DRAW\n")
    } else if winnerIx == 1 {
      t.Scores.Wins[name1] += 1
      t.Scores.Losses[name2] += 1
      fmt.Printf("Winner was %s\n", player1.Name())
    } else if winnerIx == 2 {
      t.Scores.Wins[name2] += 1
      t.Scores.Losses[name1] += 1
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
  panic("your programming is bad")
}

func (t *Tournament) ScoreFor(name string) int {
  return t.Scores.Wins[name]
}

func (t *Tournament) Size() int {
  return len(t.Players)
}

func unmarshalTournament(gt GameType, bs []byte) *Tournament {
  var s Scores
  json.Unmarshal(bs, &s)
  return newTournamentWithScores(gt, s)
}

func (t *Tournament) Marshal() []byte {
  b, _ := json.Marshal(t.Scores)
  return b
}


