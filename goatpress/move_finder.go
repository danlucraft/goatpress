package goatpress

import (
  "strings"
)

type MoveFinder interface {
  GetMove(GameState) Move
}

// PassFinder: always passes

type PassFinder struct {
}

func (f *PassFinder) GetMove(_ GameState) Move {
  return MakePassMove()
}

func newPassFinder() *PassFinder {
  return &PassFinder{}
}

// RandomFinder: searches for any valid word and returns it

type RandomFinder struct {
  words WordSet
}

func (f *RandomFinder) GetMove(state GameState) Move {
  //alreadyMovedWords := make(map[string]bool)
  wordSet := f.words
  for i := 0; i < 10000; i++ {
    // possible word
    word := wordSet.ChooseRandom()
    alreadyPlayed := false
    for _, move := range state.Moves {
      if move.Word == word || strings.HasPrefix(move.Word, word) {
        alreadyPlayed = true
      }
    }
    if !alreadyPlayed {
      move := state.Board.RandomMoveFromWord(word)
      if !move.IsPass {
        return move
      }
    }
  }
  return MakePassMove()
}

func newRandomFinder(words WordSet) *RandomFinder {
  return &RandomFinder{words}
}

