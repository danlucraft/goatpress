package goatpress

import (
	"testing"
)

func TestMatchBetweenPassesPlayers(t *testing.T) {
	gameType := newGameType(5, DefaultWordSet)
	player1 := newInternalPlayer("Alice", newPassFinder())
	player2 := newInternalPlayer("Bob", newPassFinder())
	match := newMatch(gameType, player1, player2)
	match.Play()
	winner := match.Winner()
	if winner != 0 {
		t.Errorf("a game of passing wasn't a draw")
	}
}

func TestMatchBetweenRandomPlayers(t *testing.T) {
	gameType := newGameType(5, testWordSet())
	player1 := newInternalPlayer("Alice", newRandomFinder(testWordSet()))
	player2 := newInternalPlayer("Bob", newRandomFinder(testWordSet()))
	match := newMatch(gameType, player1, player2)
	SetupBoard(&match.Game.Board)
	state := match.Game.CurrentGameState()
	if state.Board.ToString() != "hello state jenga pages valid" {
		t.Errorf("Board.ToString didn't return right string", state.Board.ToString(), "")
	}
	match.Play()
	state = match.Game.CurrentGameState()
}
