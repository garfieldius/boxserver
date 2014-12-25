package util

import (
  "os"
  "strings"
)

func FileExists(filename string) bool {
  if _, err := os.Stat(filename); err != nil {
    return false
  } else {
    return true
  }
}

func Join(parts ...string) string {
  return strings.Join(([]string)(parts), (string)(os.PathSeparator));
}
