package srch

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Checks if a path is a directory
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

// Get files in directory
func GetAllFilesInDirectory(path string, recursive bool) []string {
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

var files []string

func visitFile(fp string, fi os.DirEntry, err error) error {
	if err != nil {
		return nil
	}

	if fi.IsDir() {
		return nil
	}
	files = append(files, fp)
	return nil
}

func GetFilesRecursively(directoryPath string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.WalkDir(directoryPath, visitFile)
	if err != nil {
		return files, err
	}
	return files, nil
}

func GetCompletePath(relative_path string) (string, error) {
	working_dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return working_dir + "/" + relative_path, nil
}
