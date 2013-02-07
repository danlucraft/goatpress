package goatpress

import (
  "testing"
)

func TestNameValid(t *testing.T) {
  if !ValidateName("asdf") { t.Errorf("asdf should be valid name") }
  if ValidateName("!£") { t.Errorf("!£ should not be valid name") }
  if ValidateName("") { t.Errorf("'' should not be valid name") }
}
