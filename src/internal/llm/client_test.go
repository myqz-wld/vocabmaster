package llm

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestNewClientDefaultsToAllAdapters(t *testing.T) {
	client := mustClient(t, Options{})
	if got := client.ProviderOrder(); got != "codex -> claude -> grok" {
		t.Fatalf("ProviderOrder() = %q", got)
	}
}

func TestNewClientRequiresExplicitAdapterForModelAndThinking(t *testing.T) {
	for _, options := range []Options{
		{Model: "gpt-test"},
		{Thinking: "high"},
	} {
		if _, err := NewClient(options); err == nil || !strings.Contains(err.Error(), "显式设置 --llm-adapter") {
			t.Fatalf("NewClient(%+v) error = %v", options, err)
		}
	}
}

func TestNewClientValidatesAdapterThinkingLevels(t *testing.T) {
	tests := []struct {
		adapter  string
		thinking string
		wantErr  bool
	}{
		{adapter: "codex", thinking: "minimal"},
		{adapter: "codex", thinking: "max", wantErr: true},
		{adapter: "claude", thinking: "max"},
		{adapter: "claude", thinking: "minimal", wantErr: true},
		{adapter: "grok", thinking: "high"},
		{adapter: "grok", thinking: "xhigh", wantErr: true},
		{adapter: "grok", thinking: "unknown", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.adapter+"/"+tt.thinking, func(t *testing.T) {
			_, err := NewClient(Options{Adapter: tt.adapter, Thinking: tt.thinking})
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClientRejectsUnknownAdapter(t *testing.T) {
	_, err := NewClient(Options{Adapter: "other"})
	if err == nil || !strings.Contains(err.Error(), "auto, codex, claude, grok") {
		t.Fatalf("NewClient() error = %v", err)
	}
}

func TestClientUsesCodexFirst(t *testing.T) {
	withAvailableProviders(t)
	calls := []string{}
	withLLMProviders(t, []llmProvider{
		stubProvider(AdapterCodex, func(_ context.Context, _ string, options Options) (string, error) {
			calls = append(calls, "codex")
			if options.Adapter != "auto" {
				t.Errorf("options.Adapter = %q", options.Adapter)
			}
			return `{"chinese_def":"Codex 释义","examples":[{"sentence":"A test.","translation":"一个测试。"}]}`, nil
		}),
		stubProvider(AdapterClaude, func(context.Context, string, Options) (string, error) {
			calls = append(calls, "claude")
			return `{"chinese_def":"Claude 释义"}`, nil
		}),
		stubProvider(AdapterGrok, func(context.Context, string, Options) (string, error) {
			calls = append(calls, "grok")
			return `{"chinese_def":"Grok 释义"}`, nil
		}),
	})

	enriched, err := mustClient(t, Options{}).EnrichWord(testWord())
	if err != nil {
		t.Fatalf("EnrichWord() error = %v", err)
	}
	if got := strings.Join(calls, ","); got != "codex" {
		t.Fatalf("provider calls = %q", got)
	}
	if enriched.ChineseDef != "Codex 释义" {
		t.Errorf("ChineseDef = %q", enriched.ChineseDef)
	}
}

func TestClientFallsBackThroughGrok(t *testing.T) {
	withAvailableProviders(t)
	calls := []string{}
	withLLMProviders(t, []llmProvider{
		stubProvider(AdapterCodex, failingRun(&calls, "codex")),
		stubProvider(AdapterClaude, failingRun(&calls, "claude")),
		stubProvider(AdapterGrok, func(context.Context, string, Options) (string, error) {
			calls = append(calls, "grok")
			return `{"chinese_def":"Grok 释义"}`, nil
		}),
	})

	enriched, err := mustClient(t, Options{}).EnrichWord(testWord())
	if err != nil {
		t.Fatalf("EnrichWord() error = %v", err)
	}
	if got := strings.Join(calls, ","); got != "codex,claude,grok" {
		t.Fatalf("provider calls = %q", got)
	}
	if enriched.ChineseDef != "Grok 释义" {
		t.Errorf("ChineseDef = %q", enriched.ChineseDef)
	}
}

func TestExplicitAdapterDoesNotFallback(t *testing.T) {
	withAvailableProviders(t)
	calls := []string{}
	withLLMProviders(t, []llmProvider{
		stubProvider(AdapterCodex, func(context.Context, string, Options) (string, error) {
			calls = append(calls, "codex")
			return `{"chinese_def":"Codex 释义"}`, nil
		}),
		stubProvider(AdapterGrok, func(_ context.Context, _ string, options Options) (string, error) {
			calls = append(calls, "grok")
			if options.Model != "grok-4.5" || options.Thinking != "high" {
				t.Errorf("options = %+v", options)
			}
			return `{"chinese_def":"Grok 释义"}`, nil
		}),
	})

	client := mustClient(t, Options{Adapter: "grok", Model: "grok-4.5", Thinking: "HIGH"})
	enriched, err := client.EnrichWord(testWord())
	if err != nil {
		t.Fatalf("EnrichWord() error = %v", err)
	}
	if got := strings.Join(calls, ","); got != "grok" {
		t.Fatalf("provider calls = %q", got)
	}
	if enriched.ChineseDef != "Grok 释义" || client.ProviderOrder() != "grok" {
		t.Fatalf("unexpected client result: %#v, order=%q", enriched, client.ProviderOrder())
	}
}

func TestClientReportsAllProviderFailures(t *testing.T) {
	withAvailableProviders(t)
	withLLMProviders(t, []llmProvider{
		stubProvider(AdapterCodex, failingRun(nil, "codex")),
		stubProvider(AdapterClaude, failingRun(nil, "claude")),
		stubProvider(AdapterGrok, failingRun(nil, "grok")),
	})

	_, err := mustClient(t, Options{}).EnrichWord(testWord())
	if err == nil {
		t.Fatal("EnrichWord() error = nil")
	}
	for _, want := range []string{"codex -> claude -> grok", "codex failed", "claude failed", "grok failed"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not contain %q", err, want)
		}
	}
}

func TestIsAvailableChecksSelectedProviders(t *testing.T) {
	withLLMProviders(t, []llmProvider{
		stubProvider(AdapterCodex, nil),
		stubProvider(AdapterClaude, nil),
		stubProvider(AdapterGrok, nil),
	})
	oldLookPath := execLookPath
	execLookPath = func(file string) (string, error) {
		if file == "claude" {
			return "/tmp/claude", nil
		}
		return "", errors.New("not found")
	}
	t.Cleanup(func() { execLookPath = oldLookPath })

	if !mustClient(t, Options{}).IsAvailable() {
		t.Fatal("auto client should find Claude")
	}
	if mustClient(t, Options{Adapter: "grok"}).IsAvailable() {
		t.Fatal("selected Grok client should be unavailable")
	}
}

func mustClient(t *testing.T, options Options) *Client {
	t.Helper()
	client, err := NewClient(options)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	return client
}

func stubProvider(id Adapter, run func(context.Context, string, Options) (string, error)) llmProvider {
	return llmProvider{id: id, name: strings.ToUpper(string(id)), executable: string(id), run: run}
}

func failingRun(calls *[]string, name string) func(context.Context, string, Options) (string, error) {
	return func(context.Context, string, Options) (string, error) {
		if calls != nil {
			*calls = append(*calls, name)
		}
		return "", errors.New(name + " failed")
	}
}

func withLLMProviders(t *testing.T, providers []llmProvider) {
	t.Helper()
	oldProviders := llmProviders
	llmProviders = providers
	t.Cleanup(func() { llmProviders = oldProviders })
}

func withAvailableProviders(t *testing.T) {
	t.Helper()
	oldLookPath := execLookPath
	execLookPath = func(file string) (string, error) {
		return "/tmp/" + file, nil
	}
	t.Cleanup(func() { execLookPath = oldLookPath })
}
