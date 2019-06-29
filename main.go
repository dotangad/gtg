package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	cwd, err := os.Getwd()
	handle(err)

	gitignore, found := findGitignore(cwd, 0)
	if !found {
		log.Fatalf("Gitignore not found in path: %s", cwd)
	}

	fmt.Println(gitignore)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func findGitignore(cwd string, attempt int) (string, bool) {
	if attempt > 3 {
		return "", false
	}

	// Construct prefix to traverse up directories
	path := cwd
	for i := 0; i < attempt; i++ {
		path += "/.."
	}
	path += "/.gitignore"

	// Check if gitignore was found
	found, err := fileExists(path)
	handle(err)
	if found {
		return path, true
	}

	// If not found, run again
	return findGitignore(cwd, attempt+1)
}

func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}
