package goatpress

import (
  "testing"
)

func TestBoardGenerator(t *testing.T) {
  board := BoardGenerator{5, make([]string, 2)}
  if board.Size != 5 {
    t.Errorf("board.Size", board.Size, 5)
  }
}

func TestWordFileBoardGenerator(t *testing.T) {
  board := wordFileBoardGenerator(5, "/Users/dan/Dropbox/projects/go/src/goatpress/data/words")
  if board.Words[0] != "A" {
    t.Errorf("board.Words[0] not right word", board.Words[0], "A")
  }

  if len(board.Words) != 235886 {
    t.Errorf("board.Words not right length", len(board.Words), 235886)
  }
}

