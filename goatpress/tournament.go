package goatpress

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
)

type Tournament struct {
	GameType       GameType
	Players        map[string]Player
	Matches        []Match
	DataPath       string
	Scores         Scores
	AllPlayerNames map[string]bool
	InProgress     map[string]bool
}

type Scores struct {
	Games      map[string]int
	Moves      map[string]int
	Wins       map[string]int
	Losses     map[string]int
	WinProduct map[string]int
}

func newTournament(gt GameType, dataPath string) *Tournament {
	if dataPath == "" {
		return blankTournament(gt, "")
	}
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		return blankTournament(gt, dataPath)
	} else {
		b, _ := ioutil.ReadFile(dataPath)
		return unmarshalTournament(gt, b, dataPath)
	}
	return blankTournament(gt, dataPath)
}

func blankTournament(gt GameType, dataPath string) *Tournament {
	g := make(map[string]int)
	m := make(map[string]int)
	w := make(map[string]int)
	l := make(map[string]int)
	tb := make(map[string]int)
	tm := Scores{g, m, w, l, tb}
	return newTournamentWithScores(gt, tm, dataPath)
}

func newTournamentWithScores(gt GameType, scores Scores, dataPath string) *Tournament {
	names := make(map[string]bool)
	for name, _ := range scores.Games {
		names[name] = true
	}
	return &Tournament{gt, make(map[string]Player), make([]Match, 0), dataPath, scores, names, make(map[string]bool)}
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

type MatchResult struct {
	Player1   string
	Player2   string
	Winner    int
	MoveCount int
}

func (t *Tournament) PlayMatch() {
	if t.NonPlayingCount() > 2 {
		player1 := t.RandomNonplayingPlayer()
		player2 := t.RandomNonplayingPlayer()
		for player1.Name() == player2.Name() || (player1.Name() == "Random" && player2.Name() == "Random2") || (player2.Name() == "Random" && player1.Name() == "Random2") {
			player2 = t.RandomNonplayingPlayer()
		}
		fmt.Printf("Match between %s and %s... \n", player1.Name(), player2.Name())
		match := newMatch(&t.GameType, player1, player2)
		t.InProgress[player1.Name()] = true
		t.InProgress[player2.Name()] = true
		go AsyncPlayMatch(match)
	}
}

func (t *Tournament) RecordResult(result MatchResult) {
	name1 := result.Player1
	name2 := result.Player2
	t.Scores.Games[name1] += 1
	t.Scores.Games[name2] += 1
	t.Scores.Moves[name1] += result.MoveCount
	t.Scores.Moves[name2] += result.MoveCount
	winnerIx := result.Winner
	if winnerIx == 0 {
		fmt.Printf("DRAW\n")
	} else if winnerIx == 1 {
		t.Scores.Wins[name1] += 1
		t.Scores.Losses[name2] += 1
		t.Scores.WinProduct[name1+">"+name2] += 1
		fmt.Printf("Winner was %s\n", name1)
	} else if winnerIx == 2 {
		t.Scores.Wins[name2] += 1
		t.Scores.Losses[name1] += 1
		t.Scores.WinProduct[name2+">"+name1] += 1
		fmt.Printf("Winner was %s\n", name2)
	}
	t.InProgress[name1] = false
	t.InProgress[name2] = false
}

func AsyncPlayMatch(match *Match) {
	match.Play()
	result := MatchResult{match.Player1.Name(),
		match.Player2.Name(),
		match.Winner(),
		len(match.Game.Moves)}
	matchResults <- result
}

func (t *Tournament) RandomNonplayingPlayer() Player {
	nonPlaying := make([]string, 0)
	for name, _ := range t.Players {
		if !t.InProgress[name] {
			nonPlaying = append(nonPlaying, name)
		}
	}
	i := rand.Intn(len(nonPlaying))
	playerName := nonPlaying[i]
	return t.Players[playerName]
}

func (t *Tournament) NonPlayingCount() int {
	count := 0
	for name, _ := range t.Players {
		if !t.InProgress[name] {
			count++
		}
	}
	return count
}

func (t *Tournament) ScoreFor(name string) int {
	return t.Scores.Wins[name]
}

func (t *Tournament) Size() int {
	return len(t.Players)
}

func unmarshalTournament(gt GameType, bs []byte, dataPath string) *Tournament {
	var s Scores
	json.Unmarshal(bs, &s)
	return newTournamentWithScores(gt, s, dataPath)
}

func (t *Tournament) Marshal() []byte {
	b, err := json.Marshal(t.Scores)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return b
}

func (t *Tournament) Save() {
	if t.DataPath != "" {
		ioutil.WriteFile(t.DataPath, t.Marshal(), 0644)
	}
}
