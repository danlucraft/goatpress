package goatpress

type Player interface {
  Name() string
  MakeMove(GameState) Move
}

type InternalPlayer struct {
  MoveFinder MoveFinder
}

func newInternalPlayer(moveFinder MoveFinder) InternalPlayer {
  return InternalPlayer{moveFinder}
}

type ClientPlayer struct {

}
