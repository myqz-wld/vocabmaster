package appconfig

import (
	"os"
	"runtime"
	"testing"
)

func TestLoadMissingConfig(t *testing.T) {
	settings, err := Load(t.TempDir())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if settings.LLM != (LLMSettings{}) {
		t.Fatalf("Load() = %#v", settings)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	want := Settings{LLM: LLMSettings{
		Adapter:  "codex",
		Model:    "gpt-5.6-luna",
		Thinking: "high",
	}}

	if err := Save(dir, want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	got, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got != want {
		t.Fatalf("Load() = %#v, want %#v", got, want)
	}

	if runtime.GOOS != "windows" {
		info, err := os.Stat(Path(dir))
		if err != nil {
			t.Fatalf("Stat() error = %v", err)
		}
		if gotMode := info.Mode().Perm(); gotMode != 0600 {
			t.Fatalf("config mode = %o, want 600", gotMode)
		}
	}
}

func TestLoadRejectsInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(Path(dir), []byte("{"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if _, err := Load(dir); err == nil {
		t.Fatal("Load() error = nil")
	}
}
