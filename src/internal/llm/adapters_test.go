package llm

import (
	"reflect"
	"strings"
	"testing"
)

func TestCodexExecArgsUseCurrentNonInteractiveConfig(t *testing.T) {
	args := codexExecArgs("/tmp/out.txt", "prompt", Options{})
	want := []string{
		"exec",
		"-c", `approval_policy="never"`,
		"--sandbox", "read-only",
		"--skip-git-repo-check",
		"--ephemeral",
		"--ignore-rules",
		"--color", "never",
		"--output-last-message", "/tmp/out.txt",
		"prompt",
	}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("codexExecArgs() = %q, want %q", args, want)
	}
	if strings.Contains(strings.Join(args, "\n"), "--ask-for-approval") {
		t.Fatal("codexExecArgs() uses obsolete --ask-for-approval flag")
	}
}

func TestCodexExecArgsIncludeModelAndThinking(t *testing.T) {
	args := codexExecArgs("/tmp/out.txt", "prompt", Options{Model: "gpt-test", Thinking: "high"})
	for _, sequence := range [][]string{
		{"--model", "gpt-test"},
		{"-c", `model_reasoning_effort="high"`},
	} {
		if !containsArgs(args, sequence) {
			t.Errorf("args %q do not contain %q", args, sequence)
		}
	}
}

func TestClaudeExecArgsIncludeModelAndThinking(t *testing.T) {
	args := claudeExecArgs("prompt", Options{Model: "sonnet", Thinking: "max"})
	for _, sequence := range [][]string{
		{"--output-format", "json"},
		{"--permission-mode", "dontAsk"},
		{"--tools", ""},
		{"--model", "sonnet"},
		{"--effort", "max"},
	} {
		if !containsArgs(args, sequence) {
			t.Errorf("args %q do not contain %q", args, sequence)
		}
	}
}

func TestGrokExecArgsUseSingleTurnPlainOutput(t *testing.T) {
	args := grokExecArgs("prompt", Options{Model: "grok-4.5", Thinking: "high"})
	for _, sequence := range [][]string{
		{"--permission-mode", "dontAsk"},
		{"--tools", ""},
		{"--max-turns", "1"},
		{"--output-format", "plain"},
		{"--model", "grok-4.5"},
		{"--reasoning-effort", "high"},
		{"-p", "prompt"},
	} {
		if !containsArgs(args, sequence) {
			t.Errorf("args %q do not contain %q", args, sequence)
		}
	}
	for _, flag := range []string{"--verbatim", "--no-plan", "--no-subagents", "--no-memory", "--disable-web-search"} {
		if !containsArgs(args, []string{flag}) {
			t.Errorf("args %q do not contain %q", args, flag)
		}
	}
}

func containsArgs(args, sequence []string) bool {
	if len(sequence) == 0 || len(sequence) > len(args) {
		return false
	}
	for i := 0; i <= len(args)-len(sequence); i++ {
		if reflect.DeepEqual(args[i:i+len(sequence)], sequence) {
			return true
		}
	}
	return false
}
