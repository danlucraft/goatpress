package goatpress

import (
  "io/ioutil"
  "strings"
  "math/rand"
)

const defaultDataPath string = "/Users/dan/Dropbox/projects/go/src/goatpress/data/words"

type Board struct {
  Size    int
  Letters [][]string
}

// *** WordSet

type WordSet interface {
  Add(string)
  Includes(string) bool
  ChooseRandom() string
  Length() int
}

type HashWordSet struct {
  words map[string]bool
  words2 []string
}

func newWordSet() *HashWordSet {
  return &HashWordSet{make(map[string]bool), make([]string, 0)}
}

func (set *HashWordSet) Add(word string) {
  set.words[word] = true
  set.words2 = append(set.words2, word)
}

func (set *HashWordSet) Includes(word string) bool {
  _, found := set.words[word]
  return found
}

func (set *HashWordSet) ChooseRandom() string {
  return set.words2[rand.Intn(len(set.words2))]
}

func (set *HashWordSet) Length() int {
  return len(set.words2)
}

func newWordSetFromFile(path string) WordSet {
  b, err := ioutil.ReadFile(path)
  if err != nil {
    panic(err)
  }
  allWords := strings.Split(string(b), "\n")
  allWords = allWords[:len(allWords)-1]
  // strip words with uppercase in them and shorter than 2 characters
  words := newWordSet()
  for _, w := range allWords {
    if w == strings.ToLower(w) && len(w) > 1 {
      words.Add(w)
    }
  }
  return words
}

// *** BoardGenerator

type BoardGenerator struct {
  Words WordSet
}

func defaultBoardGenerator() *BoardGenerator {
  return newBoardGenerator(newWordSetFromFile(defaultDataPath))
}

func newBoardGenerator(words WordSet) *BoardGenerator {
  return &BoardGenerator{words}
}

func (bg *BoardGenerator) newBoard(size int) *Board {
  letters := make([][]string, size)
  for i := 0; i < size; i++ {
    letters[i] = make([]string, size)
  }
  current := 0
  for current < size*size {
    word := bg.Words.ChooseRandom()
    for _, char := range word {
      if current < size*size {
        letters[current / 5][current % 5] = string(char)
      }
      current += 1
    }
  }

  return &Board{size, letters}
}
