package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/moby/patternmatcher"
)

func List(args []string, callback func(path string) error) (int, error) {
	if len(args) != 2 {
		return 2, fmt.Errorf("usage: docker-context-ls <context>")
	}

	root := args[1]

	if stat, err := os.Stat(root); err != nil {
		return 1, err
	} else {
		if !stat.IsDir() {
			return 1, fmt.Errorf("%s: context is not a directory", root)
		}
	}

	patterns := make([]string, 0)

	content, err := os.ReadFile(filepath.Join(root, ".dockerignore"))
	if err != nil {
		if !os.IsNotExist(err) {
			return 1, err
		}
	} else {
		patterns = strings.Split(string(content), "\n")
	}

	pm, err := patternmatcher.New(patterns)
	if err != nil {
		return 1, err
	}

	// walk the context directory
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		rpath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if rpath == "." {
			return nil
		}

		ignored, err := pm.MatchesOrParentMatches(filepath.ToSlash(rpath))
		if err != nil {
			return err
		}

		if ignored {
			if d.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if !d.IsDir() {
			if err := callback(rpath); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 1, err
	}

	return 0, nil
}

func main() {
	code, err := List(os.Args, func(path string) error {
		fmt.Println(path)
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	os.Exit(code)
}
