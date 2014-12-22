package util

import (
  "os"
)

func FileExists(filename string) bool {
  if _, err := os.Stat(filename); err != nil {
    return false
  } else {
    return true
  }
}
