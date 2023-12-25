package main

import "fmt"

func main() {
  fmt.Println(makeSkipList("hello"))
}

func makeSkipList(pattern string) map[string]int {
  skip_list := make(map[string]int)
  pattern_length := len(pattern)
  for index, char := range pattern {
    skip_list[string(char)] = max(1, pattern_length - index - 1)
  }
  return skip_list
}

