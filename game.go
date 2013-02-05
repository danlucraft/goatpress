package goatpress

type Game struct {
  BoardSize int
  Words     WordSet
}

func newGame(size int, words WordSet) *Game {
  return &Game{size, words}
}
