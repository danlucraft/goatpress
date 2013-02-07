package goatpress

import (
  "fmt"
  "strconv"
)

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

func (move *Move) ToMessage() string {
  tiles := ""
  for i, tile := range move.Tiles {
    tiles += strconv.FormatInt(int64(tile.X()), 16)
    tiles += strconv.FormatInt(int64(tile.Y()), 16)
    if i < len(move.Tiles) - 1 {
      tiles += ","
    }
  }
  return fmt.Sprintf("move:%s", tiles)
}

func (move *Move) ToString() string {
  return fmt.Sprintf("<Move: %s>", move.Word)
}
