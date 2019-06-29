package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Check args
	if len(os.Args) < 2 {
		log.Fatal("Usage: gtg filename --flags")
	}

	// Get current working directory
	cwd, err := os.Getwd()
	handle(err)

	// Find gitignore
	gitignore, found := findGitignore(cwd, 0)
	if !found {
		log.Fatalf("Gitignore not found in path: %s", cwd)
	}

	// Open gitignore file
	gtgf, err := os.OpenFile(gitignore, os.O_APPEND|os.O_RDWR, os.ModePerm)
	handle(err)
	defer gtgf.Close()

	// Check if file/pattern is already in gitignore
	reader := bufio.NewReader(gtgf)
	for {
		line, err := reader.ReadString('\n')

		if strings.TrimSpace(line) == os.Args[1] {
			return
		}

		if err == io.EOF {
			break
		}
		handle(err)
	}

	// If not, append
	_, err = gtgf.WriteString(os.Args[1] + "\n")
	handle(err)
	gtgf.Sync()
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
	handle(err)

	// If not found, run again
	return findGitignore(cwd, attempt+1)
}
