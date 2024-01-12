package srch

import (
  "os"
  "strings"
)

// Color codes for terminal color printing
const (
  srch_reset   = "\033[0m"
  srch_red     = "\033[31m"
  srch_magenta = "\033[35m"
)

func colorizeSubstring(input string, startIndex int, length int, colorCode string) string {
  before := input[:startIndex]
  after := input[startIndex+length:]
  substring := input[startIndex : startIndex+length]

  return before + colorCode + substring + srch_reset + after
}

func colorizeAllOccurrences(input string, target string, colorCode string) string {
  var result string
  startIndex := 0

  for {
    index := strings.Index(input[startIndex:], target)
    if index == -1 {
      // No more occurrences found
      result += input[startIndex:]
      break
    }

    index += startIndex
    before := input[startIndex:index]
    word := input[index : index+len(target)]

    result += before + colorCode + word + srch_reset
    startIndex = index + len(target)
  }

  return result
}

func ColorizeOutput(path string, input string, engine *StringFinder) string {
  var output string

  short_path, _ := os.Getwd()

  path = path[len(short_path):] + ":"
  if path[0] == '/' {
    path = path[1:]
  }
  colorized_path := colorizeSubstring(path, 0, len(path), srch_magenta)

  output += colorized_path + " "
  output += colorizeAllOccurrences(input, engine.Pattern, srch_red)

  return output
}
