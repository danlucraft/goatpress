package goatpress

type Player interface {
  Name() string
  MakeMove(GameState) Move
}

type InternalPlayer struct {
  MoveFinder MoveFinder
}

type ClientPlayer struct {
}
