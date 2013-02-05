
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

  if game.IsValidWord([][]int {[]int {0,0}}) {
    t.Errorf("'h' was a valid word when it shouldn't have been")
  }
  if !game.IsValidWord([][]int {[]int {0,0}, []int {0,1}, []int {0,2}, []int {0,3}, []int {0,4}}) {
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

  if game.IsValidMove([][]int {[]int {0,0}}) {
    t.Errorf("'h' was a valid move when it shouldn't have been")
  }
  if !game.IsValidMove([][]int {[]int {0,0}, []int {0,1}, []int {0,2}, []int {0,3}, []int {0,4}}) {
    t.Errorf("hello wasn't a valid move when it should have been")
  }
  if !game.IsValidMove([][]int {}) {
    t.Errorf("'' wasn't a valid move when it should have been")
  }
}

