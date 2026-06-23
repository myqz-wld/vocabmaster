package buildinfo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	appName           = "vocabmaster"
	moduleName        = "github.com/vocabmaster/vocabmaster"
	installedInfo     = "vocabmaster.build-info.json"
	localBuildInfo    = "build-info.json"
	statusOK          = 0
	statusMismatch    = 1
	statusUnavailable = 2
)

var (
	ErrMetadataMissing = errors.New("build metadata missing")
	ErrSourceNotFound  = errors.New("source checkout not found")
)

type Info struct {
	Name        string `json:"name"`
	PackageName string `json:"package"`
	Version     string `json:"version,omitempty"`
	Commit      string `json:"commit"`
	ShortCommit string `json:"shortCommit"`
	Branch      string `json:"branch,omitempty"`
	Dirty       bool   `json:"dirty"`
	BuiltAt     string `json:"builtAt"`
}

type Source struct {
	Repo        string
	Commit      string
	ShortCommit string
	Branch      string
	Dirty       bool
	OriginMain  string
}

func PrintStatus(out io.Writer, checkOnly bool) int {
	info, metadataPath, err := LoadInstalled()
	if err != nil {
		if errors.Is(err, ErrMetadataMissing) {
			fmt.Fprintf(out, "安装元数据: 缺失 (%s)\n", metadataPath)
			fmt.Fprintln(out, "状态: 无法判断已安装版本，安装元数据缺失")
		} else {
			fmt.Fprintf(out, "安装元数据: 无法读取 (%s): %v\n", metadataPath, err)
			fmt.Fprintln(out, "状态: 无法判断已安装版本，安装元数据无效")
		}
		if checkOnly {
			return statusUnavailable
		}
		return statusOK
	}

	printInstalled(out, info, metadataPath)

	source, err := FindSourceCheckout(".")
	if err != nil {
		if errors.Is(err, ErrSourceNotFound) {
			fmt.Fprintln(out, "本地源码: 从当前目录向上未找到 VocabMaster checkout")
		} else {
			fmt.Fprintf(out, "本地源码: 无法读取 git 状态: %v\n", err)
		}
		fmt.Fprintln(out, "状态: 无法和本地 checkout 比较")
		if checkOnly {
			return statusUnavailable
		}
		return statusOK
	}

	printSource(out, source)
	switch {
	case sameCommit(info.Commit, source.Commit):
		if source.Dirty {
			fmt.Fprintln(out, "状态: 已安装版本匹配当前 checkout commit；当前源码还有未提交改动")
		} else {
			fmt.Fprintln(out, "状态: 已安装版本匹配当前 checkout commit")
		}
		return statusOK
	case source.OriginMain != "" && sameCommit(info.Commit, source.OriginMain):
		fmt.Fprintf(out, "本地 origin/main: %s\n", formatCommit(source.OriginMain, ""))
		fmt.Fprintln(out, "状态: 已安装版本匹配本地 origin/main，但不匹配当前 checkout commit")
		if checkOnly {
			return statusMismatch
		}
		return statusOK
	default:
		if source.OriginMain != "" {
			fmt.Fprintf(out, "本地 origin/main: %s\n", formatCommit(source.OriginMain, ""))
		}
		fmt.Fprintln(out, "状态: 已安装版本不匹配当前 checkout commit")
		if checkOnly {
			return statusMismatch
		}
		return statusOK
	}
}

func LoadInstalled() (Info, string, error) {
	candidates, err := metadataCandidates()
	if err != nil {
		return Info{}, "", err
	}
	for _, path := range candidates {
		info, err := readInfo(path)
		if err == nil {
			return info, path, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return Info{}, path, err
		}
	}
	if len(candidates) == 0 {
		return Info{}, "", ErrMetadataMissing
	}
	return Info{}, candidates[0], ErrMetadataMissing
}

func FindSourceCheckout(start string) (Source, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return Source{}, err
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		data, err := os.ReadFile(goModPath)
		if err == nil && strings.Contains(string(data), "module "+moduleName) {
			return readSource(dir)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return Source{}, ErrSourceNotFound
}

func readInfo(path string) (Info, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Info{}, err
	}

	var info Info
	if err := json.Unmarshal(data, &info); err != nil {
		return Info{}, err
	}
	if strings.TrimSpace(info.Commit) == "" {
		return Info{}, errors.New("metadata commit is empty")
	}
	if info.Name == "" {
		info.Name = appName
	}
	if info.ShortCommit == "" {
		info.ShortCommit = shortCommit(info.Commit)
	}
	return info, nil
}

func metadataCandidates() ([]string, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}

	var paths []string
	addForExecutable := func(path string) {
		dir := filepath.Dir(path)
		paths = appendUnique(paths, filepath.Join(dir, installedInfo))
		paths = appendUnique(paths, filepath.Join(dir, localBuildInfo))
	}

	addForExecutable(exe)
	if resolved, err := filepath.EvalSymlinks(exe); err == nil && resolved != exe {
		addForExecutable(resolved)
	}

	return paths, nil
}

func readSource(repo string) (Source, error) {
	commit, err := gitOutput(repo, "rev-parse", "HEAD")
	if err != nil {
		return Source{}, err
	}
	branch, _ := gitOutput(repo, "branch", "--show-current")
	status, err := gitOutput(repo, "status", "--porcelain")
	if err != nil {
		return Source{}, err
	}
	originMain, _ := gitOutput(repo, "rev-parse", "--verify", "origin/main")

	return Source{
		Repo:        repo,
		Commit:      commit,
		ShortCommit: shortCommit(commit),
		Branch:      branch,
		Dirty:       strings.TrimSpace(status) != "",
		OriginMain:  originMain,
	}, nil
}

func gitOutput(repo string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := strings.TrimSpace(string(exitErr.Stderr))
			if stderr != "" {
				return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), stderr)
			}
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func printInstalled(out io.Writer, info Info, metadataPath string) {
	version := info.Version
	if version == "" {
		version = "unknown"
	}
	fmt.Fprintf(out, "%s %s\n", displayName(info), version)
	fmt.Fprintf(out, "安装元数据: %s\n", metadataPath)
	fmt.Fprintf(out, "安装 commit: %s\n", formatCommit(info.Commit, info.ShortCommit))
	if info.Branch != "" {
		fmt.Fprintf(out, "构建分支: %s\n", info.Branch)
	}
	fmt.Fprintf(out, "构建 dirty: %t\n", info.Dirty)
	if info.BuiltAt != "" {
		fmt.Fprintf(out, "构建时间: %s\n", info.BuiltAt)
	}
}

func printSource(out io.Writer, source Source) {
	fmt.Fprintf(out, "本地源码: %s\n", source.Repo)
	fmt.Fprintf(out, "本地 commit: %s\n", formatCommit(source.Commit, source.ShortCommit))
	if source.Branch != "" {
		fmt.Fprintf(out, "本地分支: %s\n", source.Branch)
	}
	fmt.Fprintf(out, "本地 dirty: %t\n", source.Dirty)
}

func displayName(info Info) string {
	if info.Name != "" {
		return info.Name
	}
	return appName
}

func sameCommit(left string, right string) bool {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)
	if left == "" || right == "" {
		return false
	}
	return left == right || strings.HasPrefix(left, right) || strings.HasPrefix(right, left)
}

func formatCommit(commit string, short string) string {
	if short == "" {
		short = shortCommit(commit)
	}
	if short == "" || short == commit {
		return commit
	}
	return fmt.Sprintf("%s (%s)", short, commit)
}

func shortCommit(commit string) string {
	if len(commit) <= 12 {
		return commit
	}
	return commit[:12]
}

func appendUnique(paths []string, next string) []string {
	for _, path := range paths {
		if path == next {
			return paths
		}
	}
	return append(paths, next)
}
