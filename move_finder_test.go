package goatpress

import (
  "testing"
)

func TestPassFinder(t *testing.T) {
  finder := newPassFinder()
  gameType := newGameType(5, DefaultWordSet)
  game := gameType.NewGame()
  move := finder.GetMove(game.CurrentGameState())
  if !move.IsPass {
    t.Errorf("PassFinder didn't return a pass")
  }
}

func TestRandomMoveFinder(t *testing.T) {
  testWords := newWordSet()
  testWords.Add("he")
  testWords.Add("jojo")

  gameType := newGameType(5, testWords)
  game := gameType.NewGame()

  for i := 0; i < 25; i++ {
    game.Board.Letters[i / 5][i % 5] = "q"
  }
  game.Board.Letters[0][0] = "h"
  game.Board.Letters[0][1] = "e"

  finder := newRandomFinder(testWords)
  move := finder.GetMove(game.CurrentGameState())

  if move.IsPass {
    t.Errorf("RandomFinder returned a pass")
  }
  if move.Word != "he" {
    t.Errorf("RandomFinder didn't return the only word 'he'", move.Word, "he")
  }
}
