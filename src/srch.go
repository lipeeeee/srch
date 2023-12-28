package main

import (
	// Internal
	//"bufio"
	"errors"
	"fmt"
	//"io/fs"
	"log"
	"os"

	// External
	"github.com/urfave/cli/v2"
)

/*
  main
  main processing / arg parsing
  srch(path, recursive_flag)
    makeSkipList
    search files
      for each file in lines
        search & print found
*/

func main() {
	files := getAllFilesInDirectory("./", true)
	for idx, file := range files {
		fmt.Println(fmt.Sprintf("%d:%s", idx, file))
	}

	app := &cli.App{
		Name:  "srch",
		Usage: "Performant recursive file content search tool",
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

			fmt.Println("path:", path, is_directory)
			return nil
		},
	}

	// Print any erorrs thrown
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
	}
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

// Make boyer-moore skip dictionary
func makeSkipList(pattern string) map[string]int {
	skip_list := make(map[string]int)
	pattern_length := len(pattern)
	for index, char := range pattern {
		skip_list[string(char)] = max(1, pattern_length-index-1)
	}
	return skip_list
}
