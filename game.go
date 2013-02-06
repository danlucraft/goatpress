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

func (game *Game) IsValidWord(word string) bool {
  return game.Type.Words.Includes(word)
}

func (game *Game) IsValidMove(move Move) bool {
  if len(move.Tiles) == 0 { return true } // pass
  word := game.Board.WordFromTiles(move.Tiles)
  return game.IsValidWord(word) // TODO implement previous move checking
}

func (game *Game) CurrentGameState() GameState {
  colors := make([][]int, game.Type.BoardSize)
  for i := 0; i < game.Type.BoardSize; i++ {
    colors[i] = make([]int, game.Type.BoardSize)
  }
  return GameState{0, 0, game.Board, colors, make([]Move, 0)}
}

// *** GameState: a representation of the current state 
// of the game.

type GameState struct {
  Score1 int
  Score2 int
  Board  Board
  Colors [][]int
  Moves  []Move
}




