package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type buildInfo struct {
	Name        string `json:"name"`
	PackageName string `json:"package"`
	Version     string `json:"version,omitempty"`
	Commit      string `json:"commit"`
	ShortCommit string `json:"shortCommit"`
	Branch      string `json:"branch,omitempty"`
	Dirty       bool   `json:"dirty"`
	BuiltAt     string `json:"builtAt"`
}

func main() {
	root, err := repoRoot()
	if err != nil {
		fatal(err)
	}

	commit, err := gitOutput(root, "rev-parse", "HEAD")
	if err != nil {
		fatal(fmt.Errorf("read git commit: %w", err))
	}

	branch, _ := gitOutput(root, "branch", "--show-current")
	status, err := gitOutput(root, "status", "--porcelain")
	if err != nil {
		fatal(fmt.Errorf("read git status: %w", err))
	}

	info := buildInfo{
		Name:        "vocabmaster",
		PackageName: "github.com/vocabmaster/vocabmaster",
		Version:     version(root),
		Commit:      commit,
		ShortCommit: shortCommit(commit),
		Branch:      branch,
		Dirty:       strings.TrimSpace(status) != "",
		BuiltAt:     time.Now().UTC().Format(time.RFC3339),
	}

	outPath := filepath.Join(root, "build", "build-info.json")
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		fatal(err)
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fatal(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		fatal(err)
	}
}

func repoRoot() (string, error) {
	root, err := gitOutput(".", "rev-parse", "--show-toplevel")
	if err == nil {
		return root, nil
	}

	wd, wdErr := os.Getwd()
	if wdErr != nil {
		return "", wdErr
	}
	return wd, nil
}

func version(root string) string {
	if value := strings.TrimSpace(os.Getenv("VERSION")); value != "" {
		return value
	}

	tag, err := gitOutput(root, "describe", "--tags", "--abbrev=0", "--match", "v[0-9]*")
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(tag, "v")
}

func gitOutput(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := strings.TrimSpace(string(exitErr.Stderr))
			if stderr != "" {
				return "", fmt.Errorf("%s: %s", strings.Join(args, " "), stderr)
			}
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func shortCommit(commit string) string {
	if len(commit) <= 12 {
		return commit
	}
	return commit[:12]
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "write-build-info: %v\n", err)
	os.Exit(1)
}
