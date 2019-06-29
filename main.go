package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	handle(err)

	// Find gitignore
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
	gpath := cwd
	for i := 0; i < attempt; i++ {
		gpath += "/.."
	}
	gpath += "/.gitignore"

	// Check if gitignore was found
	_, err := os.Stat(gpath)
	if !os.IsNotExist(err) {
		return filepath.Clean(gpath), true
	}

	// If not found, run again
	return findGitignore(cwd, attempt+1)
}
