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
  return &Game{*bg.newBoard(gt.BoardSize), *gt, make([]Move, 0)}
}

// Game: an instance of a GameType and a Board

type Game struct {
  Board     Board
  Type      GameType
  Moves     []Move
}

const (
  MOVE_INVALID_WORD = iota
  MOVE_TOO_SHORT    = iota
  MOVE_OK           = iota
)

func (game *Game) IsValidWord(word string) bool {
  return game.Type.Words.Includes(word)
}

func (game *Game) IsValidMove(move Move) bool {
  if len(move.Tiles) == 0 { return true } // pass
  word := game.Board.WordFromTiles(move.Tiles)
  return game.IsValidWord(word) // TODO implement previous move checking
}

func (game *Game) ColorMask() ColorMask {
  return newColorMask(&game.Board, game.Moves)
}

func (game *Game) UncoloredSquareCount() int {
  r := 0
  colors := game.ColorMask()
  l := game.Type.BoardSize
  for i := 0; i < l*l; i++ {
    if colors[i / 5][i % 5] == 0 {
      r++
    }
  }
  return r
}

func (game *Game) ColorString() string {
  colorMask := game.ColorMask()
  return colorMask.ToString()
}

func (game *Game) CurrentGameState() GameState {
  colorMask := game.ColorMask()
  return GameState{game.WhoseMove(), colorMask.Score(1), colorMask.Score(2),
                    game.Board, game.ColorMask(), game.ColorString(), game.Moves}
}

// 0 game over
// 1 player 1
// 2 player 2
func (game *Game) WhoseMove() int {
  if game.UncoloredSquareCount() == 0 { return 0 }
  l := len(game.Moves)
  if l > 1 && game.Moves[l-2].IsPass && game.Moves[l-1].IsPass {
    return 0
  }
  return (len(game.Moves) % 2) + 1
}

func (game *Game) Move(move Move) int {
  if !move.IsPass && len(move.Word) < 2 {
    return MOVE_TOO_SHORT
  }
  if !move.IsPass && !game.IsValidWord(move.Word) {
    return MOVE_INVALID_WORD
  }
  game.Moves = append(game.Moves, move)
  return MOVE_OK
}

// *** GameState: a representation of the current state 
// of the game.

type GameState struct {
  WhoseMove int
  Score1    int
  Score2    int
  Board     Board
  Colors    [][]int
  ColorString string
  Moves     []Move
}

