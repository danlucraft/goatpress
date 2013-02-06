package goatpress

import (
  "testing"
)

func TestHashWordSet(t *testing.T) {
  set := newWordSet()
  set.Add("hello")
  set.Add("hi")
  a1 := set.Includes("hello")
  if !a1 { t.Errorf("include hello failed", a1, true) }
  a2 := set.Includes("hippie")
  if a2 { t.Errorf("include hippie failed", a2, false) }
  a3 := set.ChooseRandom()
  if a3 != "hi" && a3 != "hello" {
    t.Errorf("ChooseRandom didn't choose one")
  }
}

func TestNewWordSetFromFile(t *testing.T) {
  set := DefaultWordSet
  if !set.Includes("aa") { t.Errorf("wordSet doesn't include aa") }
  if set.Includes("a") { t.Errorf("wordSet includes a") }

  if set.Length() != 210661 {
    t.Errorf("wordSet not right length", set.Length(), 210661)
  }
}

