package goatpress

import (
	"testing"
)

func TestMakeGameType(t *testing.T) {
	gameType := newGameType(5, DefaultWordSet)
	if gameType.BoardSize != 5 {
		t.Errorf("gameType.BoardSize is wrong", gameType.BoardSize, 5)
	}
	if !gameType.Words.Includes("aa") {
		t.Errorf("gameType.Words isn't loaded")
	}
}

func TestGame(t *testing.T) {
	gameType := newGameType(5, DefaultWordSet)
	game := gameType.NewGame()
	if game.Board.Size != 5 {
		t.Errorf("game.Board not set up correctly")
	}
}

func TestValidWordChecking(t *testing.T) {
	gameType := newGameType(5, DefaultWordSet)
	game := gameType.NewGame()

	game.Board.Letters[0][0] = "h"
	game.Board.Letters[0][1] = "e"
	game.Board.Letters[0][2] = "l"
	game.Board.Letters[0][3] = "l"
	game.Board.Letters[0][4] = "o"

	if game.IsValidWord("h") {
		t.Errorf("'h' was a valid word when it shouldn't have been")
	}
	if !game.IsValidWord("hello") {
		t.Errorf("hello wasn't a valid word when it should have been")
	}
}

func TestValidMoveChecking(t *testing.T) {
	gameType := newGameType(5, DefaultWordSet)
	game := gameType.NewGame()

	game.Board.Letters[0][0] = "h"
	game.Board.Letters[0][1] = "e"
	game.Board.Letters[0][2] = "l"
	game.Board.Letters[0][3] = "l"
	game.Board.Letters[0][4] = "o"

	if game.IsValidMove(Move{false, []Tile{newTile(0, 0)}, "h"}) {
		t.Errorf("'h' was a valid move when it shouldn't have been")
	}
	if !game.IsValidMove(Move{false, []Tile{newTile(0, 0), newTile(0, 1), newTile(0, 2), newTile(0, 3), newTile(0, 4)}, "hello"}) {
		t.Errorf("hello wasn't a valid move when it should have been")
	}
	if !game.IsValidMove(Move{true, []Tile{}, ""}) {
		t.Errorf("'' wasn't a valid move when it should have been")
	}
}

func TestPlayGameWithWinner(t *testing.T) {
	words := testWordSet()
	words.Add("hell")
	gameType := newGameType(5, words)
	game := gameType.NewGame()
	state := game.CurrentGameState()
	SetupBoard(&game.Board)
	if game.Board.ToString() != "hello state jenga pages valid" {
		t.Errorf("Board.ToString didn't return right string", game.Board.ToString(), "")
	}

	expColorString := "00000 00000 00000 00000 00000"
	if state.ColorString != expColorString {
		t.Errorf("colorstring isn't right didn't return right string", state.ColorString, expColorString)
	}
	if state.WhoseMove != 1 {
		t.Errorf("first move is for player 1")
	}
	if state.Score1 != 0 {
		t.Errorf("Score1 is not 0")
	}
	if state.Score2 != 0 {
		t.Errorf("Score2 is not 0")
	}
	if state.Colors[0][0] != 0 {
		t.Errorf("Colors are not blank")
	}

	game.Move(game.Board.MoveFromTiles([]Tile{newTile(0, 0), newTile(0, 1), newTile(0, 2), newTile(0, 3), newTile(0, 4)}))

	state = game.CurrentGameState()
	if len(state.Moves) != 1 {
		t.Errorf("wrong number of moves", len(state.Moves), 1)
	}
	expColorString = "11111 00000 00000 00000 00000"
	if state.ColorString != expColorString {
		t.Errorf("moving didn't update colors correctly", state.ColorString, expColorString)
	}
	if state.WhoseMove != 2 {
		t.Errorf("second move is for player 2")
	}
	if state.Score1 != 5 {
		t.Errorf("Score1 is not 5")
	}
	if state.Score2 != 0 {
		t.Errorf("Score2 is not 0")
	}

	game.Move(game.Board.MoveFromTiles([]Tile{newTile(1, 0), newTile(1, 1), newTile(1, 2), newTile(1, 3), newTile(1, 4)}))

	state = game.CurrentGameState()
	if len(state.Moves) != 2 {
		t.Errorf("wrong number of moves", len(state.Moves), 2)
	}
	expColorString = "11111 22222 00000 00000 00000"
	if state.ColorString != expColorString {
		t.Errorf("moving didn't update colors correctly", state.ColorString, expColorString)
	}
	if state.WhoseMove != 1 {
		t.Errorf("third move is for player 1")
	}
	if state.Score1 != 5 {
		t.Errorf("Score1 is not 5")
	}
	if state.Score2 != 5 {
		t.Errorf("Score2 is not 5")
	}

	// invalid moves
	result := game.Move(game.Board.MoveFromTiles([]Tile{newTile(2, 0), newTile(2, 2)}))
	if result != MOVE_INVALID_WORD {
		t.Errorf("didn't identify an invalid word")
	}
	result = game.Move(game.Board.MoveFromTiles([]Tile{newTile(2, 0)}))
	if result != MOVE_TOO_SHORT {
		t.Errorf("didn't identify a too short word")
	}
	result = game.Move(game.Board.MoveFromTiles([]Tile{newTile(0, 0), newTile(0, 1), newTile(0, 2), newTile(0, 3), newTile(0, 4)}))
	if result != MOVE_ALREADY_PLAYED {
		t.Errorf("didn't identify already played word")
	}
	result = game.Move(game.Board.MoveFromTiles([]Tile{newTile(0, 0), newTile(0, 1), newTile(0, 2), newTile(0, 3)}))
	if result != MOVE_PREFIX_WORD {
		t.Errorf("didn't identify prefix word")
	}

	game.Move(game.Board.MoveFromTiles([]Tile{newTile(2, 0), newTile(2, 1), newTile(2, 2), newTile(2, 3), newTile(2, 4)}))
	game.Move(game.Board.MoveFromTiles([]Tile{newTile(3, 0), newTile(3, 1), newTile(3, 2), newTile(3, 3), newTile(3, 4)}))
	game.Move(game.Board.MoveFromTiles([]Tile{newTile(4, 0), newTile(4, 1), newTile(4, 2), newTile(4, 3), newTile(4, 4)}))

	state = game.CurrentGameState()
	if len(state.Moves) != 5 {
		t.Errorf("wrong number of moves", len(state.Moves), 5)
	}
	expColorString = "11111 22222 11111 22222 11111"
	if state.ColorString != expColorString {
		t.Errorf("moving didn't update colors correctly", state.ColorString, expColorString)
	}
	if state.WhoseMove != 0 {
		t.Errorf("should be game over")
	}
	if state.Score1 != 15 {
		t.Errorf("Score1 is not 15")
	}
	if state.Score2 != 10 {
		t.Errorf("Score2 is not 10")
	}
}

func TestPlayGameWithPassing(t *testing.T) {
	gameType := newGameType(5, testWordSet())
	game := gameType.NewGame()
	state := game.CurrentGameState()
	SetupBoard(&game.Board)

	result1 := game.Move(MakePassMove())
	result2 := game.Move(MakePassMove())
	if result1 != MOVE_OK || result2 != MOVE_OK {
		t.Errorf("passes weren't ok moves")
	}
	state = game.CurrentGameState()
	if len(state.Moves) != 2 {
		t.Errorf("wrong number of moves", len(state.Moves), 2)
	}
	if state.WhoseMove != 0 {
		t.Errorf("should be game over")
	}
	if state.Score1 != 0 {
		t.Errorf("Score1 is not 0")
	}
	if state.Score2 != 0 {
		t.Errorf("Score2 is not 0")
	}
}
