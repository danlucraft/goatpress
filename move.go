package goatpress

type Move struct {
  IsPass      bool
  Tiles       []Tile
  Word   string
}

func MakePassMove() Move {
  return Move{true, make([]Tile, 0), ""}
}

func (move *Move) HasTile(tile Tile) bool {
  for _, tile2 := range move.Tiles {
    if tile2[0] == tile[0] && tile2[1] == tile[1] {
      return true
    }
  }
  return false
}

