
package goatpress

import (
  "testing"
)

func testMakeGame(t *testing.T) {
  game := newGame(5, defaultWordSet())
  if game.BoardSize != 5 {
    t.Errorf("game.BoardSize is wrong", game.BoardSize, 5)
  }
  if !game.Words.Includes("aa") {
    t.Errorf("game.Words isn't loaded")
  }

}
