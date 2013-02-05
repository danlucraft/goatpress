package goatpress

// GameType: encapsulates settings for a game (size, which words, etc)

type GameType struct {
  BoardSize  int
  Words      WordSet
}

func newGameType(size int, words WordSet) *GameType {
  return &GameType{size, words}
}

func (gt *GameType) NewGame() *Game {
  bg := &BoardGenerator{gt.Words}
  return &Game{*bg.newBoard(gt.BoardSize), *gt}
}

// Game: an instance of a GameType and a Board

type Game struct {
  Board     Board
  Type      GameType
}

func (game *Game) IsValidWord(move [][]int) bool {
  word := game.Board.WordFromMove(move)
  return game.Type.Words.Includes(word)
}

func (game *Game) IsValidMove(move [][]int) bool {
  return game.IsValidWord(move) // TODO implement previous move checking
}
