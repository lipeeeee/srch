package srch

import (
  "os"
  "bufio"
  "fmt"
)

func Srch(engine *StringFinder, path string) error {
  // Open file
  file, err := os.Open(path)
  if err != nil {
    return err
  }
  defer file.Close()

  // Create file reader
  scanner := bufio.NewScanner(file)
  idx := 0

  for scanner.Scan() {
    idx++
    text_to_search := scanner.Text()

    // Create buffer to stored this line's indicies
    var indicies []int = make([]int, 0)
    current_found := false 

    // Keep iterating on each found pattern 
    found_index := engine.Next(text_to_search)
    for found_index != -1 {
      indicies = append(indicies, found_index)
      text_to_search = text_to_search[found_index + engine.Length:]
      found_index = engine.Next(text_to_search)
      current_found = true
    }

    if current_found {
      printFind(path, idx, indicies, engine, scanner.Text())
    }
  }
  
  if err := scanner.Err(); err != nil {
    return err
  }
  return nil
}

// Prints to STDOUT result of a single find
func printFind(path string, line_num int, indicies []int, engine *StringFinder, txt string) {
    coloredString := ColorizeOutput(path, txt, engine)
    fmt.Println(coloredString)
}

