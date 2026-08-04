package cmd

import (
	"testing"

	"github.com/vocabmaster/vocabmaster/src/internal/appconfig"
	"github.com/vocabmaster/vocabmaster/src/internal/llm"
)

func TestMergeLLMOptionsUsesSavedDefaults(t *testing.T) {
	saved := appconfig.LLMSettings{Adapter: "codex", Model: "gpt-5.6-luna", Thinking: "high"}
	got := mergeLLMOptions(saved, llm.Options{Adapter: "auto"}, false, false, false)
	want := llm.Options{Adapter: "codex", Model: "gpt-5.6-luna", Thinking: "high"}
	if got != want {
		t.Fatalf("mergeLLMOptions() = %#v, want %#v", got, want)
	}
}

func TestMergeLLMOptionsFlagsOverrideSavedDefaults(t *testing.T) {
	saved := appconfig.LLMSettings{Adapter: "codex", Model: "gpt-5.6-luna", Thinking: "high"}
	flags := llm.Options{Adapter: "codex", Model: "gpt-test", Thinking: "low"}
	got := mergeLLMOptions(saved, flags, false, true, true)
	want := llm.Options{Adapter: "codex", Model: "gpt-test", Thinking: "low"}
	if got != want {
		t.Fatalf("mergeLLMOptions() = %#v, want %#v", got, want)
	}
}

func TestMergeLLMOptionsChangingAdapterClearsProviderSpecificDefaults(t *testing.T) {
	saved := appconfig.LLMSettings{Adapter: "codex", Model: "gpt-5.6-luna", Thinking: "high"}
	flags := llm.Options{Adapter: "claude"}
	got := mergeLLMOptions(saved, flags, true, false, false)
	want := llm.Options{Adapter: "claude"}
	if got != want {
		t.Fatalf("mergeLLMOptions() = %#v, want %#v", got, want)
	}
}

func TestMergeLLMOptionsDefaultsToAuto(t *testing.T) {
	got := mergeLLMOptions(appconfig.LLMSettings{}, llm.Options{Adapter: "auto"}, false, false, false)
	want := llm.Options{Adapter: "auto"}
	if got != want {
		t.Fatalf("mergeLLMOptions() = %#v, want %#v", got, want)
	}
}
