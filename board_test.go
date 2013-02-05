package goatpress

import (
  "testing"
)

func TestHashWordSet(t *testing.T) {
  set := newWordSet()
  set.Add("hello")
  set.Add("hi")
  a1 := set.Includes("hello")
  if !a1 { t.Errorf("include hello failed", a1, true) }
  a2 := set.Includes("hippie")
  if a2 { t.Errorf("include hippie failed", a2, false) }
  a3 := set.ChooseRandom()
  if a3 != "hi" && a3 != "hello" {
    t.Errorf("ChooseRandom didn't choose one")
  }
}

func TestNewWordSetFromFile(t *testing.T) {
  set := newWordSetFromFile(defaultDataPath)
  if !set.Includes("aa") { t.Errorf("wordSet doesn't include aa") }
  if set.Includes("a") { t.Errorf("wordSet includes a") }

  if set.Length() != 210661 {
    t.Errorf("wordSet not right length", set.Length(), 210661)
  }
}

func TestMakingNewBoards(t *testing.T) {
  bg := defaultBoardGenerator()
  board := bg.newBoard(5)
  if board.Size != 5 {
    t.Errorf("board.Size is not right", board.Size, 5)
  }
  if len(board.Letters) != 5 {
    t.Errorf("board.Letters is not right size", len(board.Letters), 5)
  }
  if len(board.Letters[0]) != 5 {
    t.Errorf("board.Letters[0] is not right size", len(board.Letters[0]), 5)
  }
  if (board.Letters[0][0] == "") {
    t.Errorf("board.Letters hasn't been filled in")
  }
}

func TestBoard(t *testing.T) {
  bg := defaultBoardGenerator()
  board := bg.newBoard(5)
  board.Letters[0][0] = "h"
  board.Letters[0][1] = "e"
  board.Letters[0][2] = "l"
  board.Letters[0][3] = "l"
  board.Letters[0][4] = "o"
}




