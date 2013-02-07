package goatpress

import (
  "io/ioutil"
  "strings"
  "math/rand"
  "runtime"
  "path"
)

const defaultDataPath string = "data/words"

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

var DefaultWordSet = defaultWordSet()

func defaultWordSet() *HashWordSet {
  _, filename, _, _ := runtime.Caller(1)
  path := path.Join(path.Dir(filename), defaultDataPath)
  return newWordSetFromFile(path)
}

func newWordSetFromFile(path string) *HashWordSet {
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

