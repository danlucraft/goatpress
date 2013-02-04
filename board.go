package goatpress

import (
  "io/ioutil"
  "strings"
)

type Move struct {
  Squares []int
}

type Board struct {
  Letters [][]int
  Moves []Move
}

type BoardGenerator struct {
  Size int
  Words []string
}

func wordFileBoardGenerator(size int, path string) *BoardGenerator {
  b, err := ioutil.ReadFile(path)
  if err != nil {
    panic(err)
  }
  words := strings.Split(string(b), "\n")
  words = words[:len(words)-1]
  return &BoardGenerator{size, words}
}
