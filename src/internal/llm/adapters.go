package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCodex(ctx context.Context, prompt string, options Options) (string, error) {
	outputFile, err := os.CreateTemp("", "vocabmaster-codex-*.txt")
	if err != nil {
		return "", fmt.Errorf("创建 Codex 输出文件失败: %w", err)
	}
	outputPath := outputFile.Name()
	if err := outputFile.Close(); err != nil {
		_ = os.Remove(outputPath)
		return "", fmt.Errorf("关闭 Codex 输出文件失败: %w", err)
	}
	defer os.Remove(outputPath)

	cmd := exec.CommandContext(ctx, "codex", codexExecArgs(outputPath, prompt, options)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("调用 Codex CLI 失败: %w%s", err, commandOutputPreview(output))
	}

	result, err := os.ReadFile(outputPath)
	if err != nil {
		return "", fmt.Errorf("读取 Codex 响应失败: %w", err)
	}
	raw := strings.TrimSpace(string(result))
	if raw == "" {
		raw = strings.TrimSpace(string(output))
	}
	if raw == "" {
		return "", fmt.Errorf("Codex 返回了空响应")
	}
	return raw, nil
}

func codexExecArgs(outputPath, prompt string, options Options) []string {
	args := []string{
		"exec",
		"-c", `approval_policy="never"`,
		"--sandbox", "read-only",
		"--skip-git-repo-check",
		"--ephemeral",
		"--ignore-rules",
		"--color", "never",
	}
	if options.Model != "" {
		args = append(args, "--model", options.Model)
	}
	if options.Thinking != "" {
		args = append(args, "-c", "model_reasoning_effort="+strconv.Quote(options.Thinking))
	}
	return append(args, "--output-last-message", outputPath, prompt)
}

func runClaude(ctx context.Context, prompt string, options Options) (string, error) {
	cmd := exec.CommandContext(ctx, "claude", claudeExecArgs(prompt, options)...)
	output, err := cmd.Output()
	if err != nil {
		preview := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			preview = commandOutputPreview(exitErr.Stderr)
		}
		return "", fmt.Errorf("调用 Claude CLI 失败: %w%s", err, preview)
	}

	var response struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(output, &response); err != nil {
		return "", fmt.Errorf("解析 Claude 响应失败: %w", err)
	}
	if strings.TrimSpace(response.Result) == "" {
		return "", fmt.Errorf("Claude 返回了空响应")
	}
	return response.Result, nil
}

func claudeExecArgs(prompt string, options Options) []string {
	args := []string{
		"-p", prompt,
		"--output-format", "json",
		"--permission-mode", "dontAsk",
		"--tools", "",
		"--no-session-persistence",
	}
	if options.Model != "" {
		args = append(args, "--model", options.Model)
	}
	if options.Thinking != "" {
		args = append(args, "--effort", options.Thinking)
	}
	return args
}

func runGrok(ctx context.Context, prompt string, options Options) (string, error) {
	cmd := exec.CommandContext(ctx, "grok", grokExecArgs(prompt, options)...)
	output, err := cmd.Output()
	if err != nil {
		preview := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			preview = commandOutputPreview(exitErr.Stderr)
		}
		return "", fmt.Errorf("调用 Grok CLI 失败: %w%s", err, preview)
	}
	raw := strings.TrimSpace(string(output))
	if raw == "" {
		return "", fmt.Errorf("Grok 返回了空响应")
	}
	return raw, nil
}

func grokExecArgs(prompt string, options Options) []string {
	args := []string{
		"--verbatim",
		"--no-plan",
		"--no-subagents",
		"--no-memory",
		"--disable-web-search",
		"--permission-mode", "dontAsk",
		"--tools", "",
		"--max-turns", "1",
		"--output-format", "plain",
	}
	if options.Model != "" {
		args = append(args, "--model", options.Model)
	}
	if options.Thinking != "" {
		args = append(args, "--reasoning-effort", options.Thinking)
	}
	return append(args, "-p", prompt)
}

func commandOutputPreview(output []byte) string {
	preview := strings.TrimSpace(string(output))
	if preview == "" {
		return ""
	}
	if len(preview) > 300 {
		preview = preview[:300] + "..."
	}
	return ": " + preview
}
