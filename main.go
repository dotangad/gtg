package main

import (
	"bufio"
	"io"
	"log"
	"os"
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
	found, err = checkInFile(gtgf, os.Args[1])
	if found {
		return
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

func checkInFile(f *os.File, p string) (bool, error) {
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')

		if strings.TrimSpace(line) == p {
			return true, nil
		}

		if err == io.EOF {
			break
		}
		handle(err)
	}

	return false, nil
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
		return gpath, true
	}
	if os.IsNotExist(err) {
		// If not found, run again
		return findGitignore(cwd, attempt+1)
	}
	handle(err)

	return findGitignore(cwd, attempt+1)
}
