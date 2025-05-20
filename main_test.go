package main

import (
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"
)

func TestDockerBuildContextListing(t *testing.T) {
	tmp := t.TempDir()

	dockerignore := []string{
		"# comment should be ignored",
		"*.log",
		"**/ignore.txt",
		"ignore/",
		"nested/dir/to/ignore/",
		"build/",
	}

	if err := os.WriteFile(filepath.Join(tmp, ".dockerignore"), []byte(strings.Join(dockerignore, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		"1.txt":                      "keep",
		"2.log":                      "ignore; wildcard",
		"ignore.txt":                 "ignore",
		"ignore/Dockerfile":          "ignore",
		"ignore/dir/3.txt":           "ignore",
		"nested/4.txt":               "keep",
		"nested/dir/5.txt":           "keep",
		"nested/dir/6.log":           "keep",
		"nested/dir/ignore.txt":      "ignore",
		"nested/dir/to/ignore/7.txt": "ignore",
		"Dockerfile":                 "keep; special case: Dockerfile",
	}

	for rpath, content := range files {
		path := filepath.Join(tmp, rpath)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	results := make([]string, 0)

	if code, err := List([]string{"", tmp}, func(path string) error {
		results = append(results, path)
		return nil
	}); err != nil || code != 0 {
		t.Fatal(err)
	}

	expected := []string{".dockerignore"}
	for rpath, content := range files {
		if !strings.HasPrefix(content, "keep") {
			continue
		}

		expected = append(expected, rpath)
	}

	sort.Strings(expected)
	sort.Strings(results)

	if !slices.Equal(results, expected) {
		t.Fatalf("results not expected: '%s' vs. '%s'", strings.Join(results, ", "), strings.Join(expected, ", "))
	}
}
