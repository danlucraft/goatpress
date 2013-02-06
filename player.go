package goatpress

type Player interface {
  Name() string
  GetMove(GameState) Move
}

type InternalPlayer struct {
  name       string
  MoveFinder MoveFinder
}

func newInternalPlayer(name string, moveFinder MoveFinder) InternalPlayer {
  return InternalPlayer{name, moveFinder}
}

func (p InternalPlayer) Name() string {
  return p.name
}

func (p InternalPlayer) GetMove(state GameState) Move {
  moveFinder := p.MoveFinder
  return moveFinder.GetMove(state)
}

type ClientPlayer struct {

}
