package llm

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/vocabmaster/vocabmaster/src/internal/model"
)

type Adapter string

const (
	AdapterAuto   Adapter = "auto"
	AdapterCodex  Adapter = "codex"
	AdapterClaude Adapter = "claude"
	AdapterGrok   Adapter = "grok"

	providerTimeout = 60 * time.Second
)

type Options struct {
	Adapter  string
	Model    string
	Thinking string
}

type Client struct {
	options   Options
	providers []llmProvider
}

type llmProvider struct {
	id         Adapter
	name       string
	executable string
	run        func(context.Context, string, Options) (string, error)
}

var llmProviders = []llmProvider{
	{id: AdapterCodex, name: "Codex", executable: "codex", run: runCodex},
	{id: AdapterClaude, name: "Claude", executable: "claude", run: runClaude},
	{id: AdapterGrok, name: "Grok", executable: "grok", run: runGrok},
}

var execLookPath = exec.LookPath

func NewClient(options Options) (*Client, error) {
	normalized, err := NormalizeOptions(options)
	if err != nil {
		return nil, err
	}

	providers := llmProviders
	if Adapter(normalized.Adapter) != AdapterAuto {
		providers = nil
		for _, provider := range llmProviders {
			if provider.id == Adapter(normalized.Adapter) {
				providers = append(providers, provider)
				break
			}
		}
	}

	return &Client{options: normalized, providers: providers}, nil
}

func NormalizeOptions(options Options) (Options, error) {
	options.Adapter = strings.ToLower(strings.TrimSpace(options.Adapter))
	options.Model = strings.TrimSpace(options.Model)
	options.Thinking = strings.ToLower(strings.TrimSpace(options.Thinking))
	if options.Adapter == "" {
		options.Adapter = string(AdapterAuto)
	}

	adapter := Adapter(options.Adapter)
	switch adapter {
	case AdapterAuto, AdapterCodex, AdapterClaude, AdapterGrok:
	default:
		return Options{}, fmt.Errorf("无效的 LLM adapter %q，可选值: auto, codex, claude, grok", options.Adapter)
	}

	if adapter == AdapterAuto && (options.Model != "" || options.Thinking != "") {
		return Options{}, fmt.Errorf("--llm-model 和 --llm-thinking 需要显式设置 --llm-adapter（codex、claude 或 grok）")
	}
	if strings.ContainsAny(options.Model, "\r\n") {
		return Options{}, fmt.Errorf("LLM model 不能包含换行符")
	}
	if options.Thinking != "" && !supportsThinking(adapter, options.Thinking) {
		return Options{}, fmt.Errorf(
			"adapter %s 不支持思考程度 %q，可选值: %s",
			adapter,
			options.Thinking,
			strings.Join(thinkingLevels(adapter), ", "),
		)
	}

	return options, nil
}

func thinkingLevels(adapter Adapter) []string {
	switch adapter {
	case AdapterCodex:
		return []string{"low", "medium", "high", "xhigh", "max"}
	case AdapterClaude:
		return []string{"low", "medium", "high", "xhigh", "max"}
	case AdapterGrok:
		return []string{"low", "medium", "high", "xhigh"}
	default:
		return nil
	}
}

func supportsThinking(adapter Adapter, thinking string) bool {
	for _, supported := range thinkingLevels(adapter) {
		if thinking == supported {
			return true
		}
	}
	return false
}

func (c *Client) EnrichWord(word *model.Word) (*model.Word, error) {
	if word == nil {
		return nil, fmt.Errorf("word is nil")
	}

	prompt := buildEnrichPrompt(word)
	var failures []string
	attempted := false

	for _, provider := range c.providers {
		if !isProviderAvailable(provider) {
			failures = append(failures, fmt.Sprintf("%s: 未找到 CLI", provider.id))
			continue
		}

		attempted = true
		ctx, cancel := context.WithTimeout(context.Background(), providerTimeout)
		raw, err := provider.run(ctx, prompt, c.options)
		if ctx.Err() == context.DeadlineExceeded {
			err = fmt.Errorf("调用 %s CLI 超时（60秒）", provider.name)
		}
		cancel()
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", provider.id, err))
			continue
		}

		enriched, err := parseEnrichResponse(word, raw)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", provider.id, err))
			continue
		}
		return enriched, nil
	}

	if !attempted {
		return nil, fmt.Errorf("未找到可用 LLM CLI（尝试顺序: %s）", c.ProviderOrder())
	}
	return nil, fmt.Errorf("LLM 增强失败（尝试顺序: %s）: %s", c.ProviderOrder(), strings.Join(failures, "; "))
}

func (c *Client) IsAvailable() bool {
	for _, provider := range c.providers {
		if isProviderAvailable(provider) {
			return true
		}
	}
	return false
}

func (c *Client) ProviderOrder() string {
	names := make([]string, 0, len(c.providers))
	for _, provider := range c.providers {
		names = append(names, string(provider.id))
	}
	return strings.Join(names, " -> ")
}

func isProviderAvailable(provider llmProvider) bool {
	_, err := execLookPath(provider.executable)
	return err == nil
}
