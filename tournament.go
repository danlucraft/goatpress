package goatpress

type Tournament struct {
  GameType GameType
}

func newTournament(gt GameType) *Tournament {
  return &Tournament{gt}
}

func (t *Tournament) Size() int {
  return 0
}
