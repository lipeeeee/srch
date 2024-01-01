package main

import (
  // Internal
  "bufio"
  "errors"
  "fmt"
  "log"
  "os"
  "strings"

  // External
  "github.com/urfave/cli/v2"
)

// stringFinder efficiently finds strings in a source text. It's implemented
// using the Boyer-Moore string search algorithm:
// https://en.wikipedia.org/wiki/Boyer-Moore_string_search_algorithm
// https://www.cs.utexas.edu/~moore/publications/fstrpos.pdf (note: this aged
// document uses 1-based indexing)
type stringFinder struct {
  // pattern is the string that we are searching for in the text.
  pattern string
  length int

  // badCharSkip[b] contains the distance between the last byte of pattern
  // and the rightmost occurrence of b in pattern. If b is not in pattern,
  // badCharSkip[b] is len(pattern).
  //
  // Whenever a mismatch is found with byte b in the text, we can safely
  // shift the matching frame at least badCharSkip[b] until the next time
  // the matching char could be in alignment.
  badCharSkip [256]int

  // goodSuffixSkip[i] defines how far we can shift the matching frame given
  // that the suffix pattern[i+1:] matches, but the byte pattern[i] does
  // not. There are two cases to consider:
  //
  // 1. The matched suffix occurs elsewhere in pattern (with a different
  // byte preceding it that we might possibly match). In this case, we can
  // shift the matching frame to align with the next suffix chunk. For
  // example, the pattern "mississi" has the suffix "issi" next occurring
  // (in right-to-left order) at index 1, so goodSuffixSkip[3] ==
  // shift+len(suffix) == 3+4 == 7.
  //
  // 2. If the matched suffix does not occur elsewhere in pattern, then the
  // matching frame may share part of its prefix with the end of the
  // matching suffix. In this case, goodSuffixSkip[i] will contain how far
  // to shift the frame to align this portion of the prefix to the
  // suffix. For example, in the pattern "abcxxxabc", when the first
  // mismatch from the back is found to be in position 3, the matching
  // suffix "xxabc" is not found elsewhere in the pattern. However, its
  // rightmost "abc" (at position 6) is a prefix of the whole pattern, so
  // goodSuffixSkip[3] == shift+len(suffix) == 6+5 == 11.
  goodSuffixSkip []int
}

// srch entry point
func main() {
  var recursive bool

  app := &cli.App{
    Name:  "srch",
    Usage: "Performant recursive file content search tool",
    Flags: []cli.Flag {
      &cli.BoolFlag {
        Name:         "recursive",
        Aliases:      []string{"r"},
        Value:        false,
        Usage:        "Flag for recursive directory searching",
        Destination:  &recursive,
      },
    },
    Action: func(cCtx *cli.Context) error {
      // Check number of args passed in before doing anything
      // 0      -> provide breaf usage explanation
      // not 2  -> return an error
      switch cCtx.Args().Len() {
      case 0:
        fmt.Println(`USAGE: srch [OPTION] PATTERN [DIRECTORY/FILE]
Try 'srch --help' for more information`)
        return nil
      case 2:
        break
      default:
        return errors.New(fmt.Sprintf("%s%d",
          "Expected 2 arguments... Got ", cCtx.Args().Len()))
      }

      // Get complete path to search
      path, err := getCompletePath(cCtx.Args().Get(1))
      if err != nil {
        return err
      }
      is_directory, err := isDirectory(path)
      if err != nil {
        return err
      }

      // Initialize srch engine
      engine := makeStringFinder(cCtx.Args().Get(0))

      // If it is a file, simply call srch with the path
      if !is_directory {
        err := srch(engine, path)
        if err != nil {
          return err
        }
      } else {

      }

      return nil
    },
  }

  // Print any erorrs thrown
  if err := app.Run(os.Args); err != nil {
    fmt.Println(err.Error())
  }
}

// srch should be called with:
//  - engine:     pointer of an instance of srch engine
//  - path:       path to a specific file to analyze
func srch(engine *stringFinder, path string) error {
  // 1. Open file
  file, err := os.Open(path)
  if err != nil {
    return err
  }
  defer file.Close()

  // 2. Create file reader
  scanner := bufio.NewScanner(file)
  idx := 0

  // 3. For each line in file, we will iterate on that same line until engine
  // does not find anything, whilst printing everything we find
  for scanner.Scan() {
    idx++
    text_to_search := scanner.Text()

    // Keep iterating on each found pattern 
    found_index := engine.next(text_to_search)
    for found_index != -1 {
      printFind(path, idx, found_index, engine.pattern, text_to_search)
      text_to_search = text_to_search[found_index + len(engine.pattern):]
      found_index = engine.next(text_to_search)
    }
  }
  
  if err := scanner.Err(); err != nil {
    return err
  }
  return nil
}

// Prints to STDOUT result of a single find
func printFind(path string, line_num int, found_index int, pat string, txt string) {
  fmt.Println(fmt.Sprintf("Found: %d-%s-%d", line_num, txt, found_index + 1))
}

// Gets complete path given a relative in cli execution
func getCompletePath(relative_path string) (string, error) {
  working_dir, err := os.Getwd()
  if err != nil {
    return "", err
  }
  return working_dir + "/" + relative_path, nil
}

// Checks if a path is a directory
func isDirectory(path string) (bool, error) {
  fileInfo, err := os.Stat(path)
  if err != nil {
    return false, err
  }

  return fileInfo.IsDir(), err
}

// Get files in directory
func getAllFilesInDirectory(path string, recursive bool) []string {
  var c []string
  entries, err := os.ReadDir(path)
  if err != nil {
    log.Fatal(err)
  }

  for _, e := range entries {
    fmt.Println(e.Name())
    /*if recursive && e.IsDir() {
    	recursive_dir_call := getAllFilesInDirectory(path+"/"+e.Name(), recursive)
    	for _, entry := range recursive_dir_call {
    		entries = append(entries, entry)
    	}
    }*/
  }

  return c
}

func makeStringFinder(pattern string) *stringFinder {
  f := &stringFinder{
    pattern:        pattern,
    length:         len(pattern),
    goodSuffixSkip: make([]int, len(pattern)),
  }
  // last is the index of the last character in the pattern.
  last := len(pattern) - 1

  // Build bad character table.
  // Bytes not in the pattern can skip one pattern's length.
  for i := range f.badCharSkip {
  	f.badCharSkip[i] = len(pattern)
  }
  // The loop condition is < instead of <= so that the last byte does not
  // have a zero distance to itself. Finding this byte out of place implies
  // that it is not in the last position.
  for i := 0; i < last; i++ {
  	f.badCharSkip[pattern[i]] = last - i
  }

  // Build good suffix table.
  // First pass: set each value to the next index which starts a prefix of
  // pattern.
  lastPrefix := last
  for i := last; i >= 0; i-- {
  	if strings.HasPrefix(pattern, pattern[i+1:]) {
  		lastPrefix = i + 1
  	}
  	// lastPrefix is the shift, and (last-i) is len(suffix).
  	f.goodSuffixSkip[i] = lastPrefix + last - i
  }
  // Second pass: find repeats of pattern's suffix starting from the front.
  for i := 0; i < last; i++ {
  	lenSuffix := longestCommonSuffix(pattern, pattern[1:i+1])
  	if pattern[i-lenSuffix] != pattern[last-lenSuffix] {
  		// (last-i) is the shift, and lenSuffix is len(suffix).
  		f.goodSuffixSkip[last-lenSuffix] = lenSuffix + last - i
  	}
  }

  return f
}

func longestCommonSuffix(a, b string) (i int) {
  for ; i < len(a) && i < len(b); i++ {
  	if a[len(a)-1-i] != b[len(b)-1-i] {
  		break
  	}
  }
  return
}

// next returns the index in text of the first occurrence of the pattern. If
// the pattern is not found, it returns -1.
func (f *stringFinder) next(text string) int {
  i := len(f.pattern) - 1
  for i < len(text) {
  	// Compare backwards from the end until the first unmatching character.
  	j := len(f.pattern) - 1
  	for j >= 0 && text[i] == f.pattern[j] {
  		i--
  		j--
  	}
  	if j < 0 {
  		return i + 1 // match
  	}
  	i += max(f.badCharSkip[text[i]], f.goodSuffixSkip[j])
  }
  return -1
}

