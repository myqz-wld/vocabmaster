package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/src/internal/appconfig"
	"github.com/vocabmaster/vocabmaster/src/internal/buildinfo"
	"github.com/vocabmaster/vocabmaster/src/internal/library"
	"github.com/vocabmaster/vocabmaster/src/internal/llm"
	"github.com/vocabmaster/vocabmaster/src/internal/store"
)

var (
	dataDir string
	db      *store.SQLiteStore
	lib     *library.Library

	llmAdapter  string
	llmModel    string
	llmThinking string
	llmClient   *llm.Client
	settings    appconfig.Settings
)

var rootCmd = &cobra.Command{
	Use:   "vocabmaster",
	Short: "命令行背单词工具 - 支持英文和日文",
	Long:  "VocabMaster 是一个命令行背单词工具，支持英文和日文单词学习，\n基于 SM-2 间隔重复算法，带有中文释义、发音标注和例句。",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" || cmd.Name() == "version" || cmd.Name() == "check-installed" {
			return nil
		}

		if err := prepareDataDir(); err != nil {
			return err
		}

		loaded, err := appconfig.Load(dataDir)
		if err != nil {
			return fmt.Errorf("读取本地配置失败: %w", err)
		}
		settings = loaded

		if isConfigCommand(cmd) {
			return nil
		}

		configuredLLM, err := llm.NewClient(mergeLLMOptions(settings.LLM, llm.Options{
			Adapter:  llmAdapter,
			Model:    llmModel,
			Thinking: llmThinking,
		}, llmFlagChanged(cmd, "llm-adapter"), llmFlagChanged(cmd, "llm-model"), llmFlagChanged(cmd, "llm-thinking")))
		if err != nil {
			return err
		}
		llmClient = configuredLLM

		db, err = store.NewSQLiteStore(filepath.Join(dataDir, "vocabmaster.db"))
		if err != nil {
			return fmt.Errorf("打开数据库失败: %w", err)
		}

		lib, err = library.NewLibrary()
		if err != nil {
			return fmt.Errorf("加载词库失败: %w", err)
		}

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if db != nil {
			db.Close()
		}
	},
}

func prepareDataDir() error {
	if dataDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户目录失败: %w", err)
		}
		dataDir = filepath.Join(home, ".vocabmaster")
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("创建数据目录失败: %w", err)
	}
	return nil
}

func mergeLLMOptions(saved appconfig.LLMSettings, flags llm.Options, adapterChanged, modelChanged, thinkingChanged bool) llm.Options {
	if strings.TrimSpace(saved.Adapter) == "" {
		saved.Adapter = string(llm.AdapterAuto)
	}
	effective := llm.Options{
		Adapter:  saved.Adapter,
		Model:    saved.Model,
		Thinking: saved.Thinking,
	}

	if adapterChanged {
		if !strings.EqualFold(strings.TrimSpace(flags.Adapter), strings.TrimSpace(saved.Adapter)) {
			if !modelChanged {
				effective.Model = ""
			}
			if !thinkingChanged {
				effective.Thinking = ""
			}
		}
		effective.Adapter = flags.Adapter
	}
	if modelChanged {
		effective.Model = flags.Model
	}
	if thinkingChanged {
		effective.Thinking = flags.Thinking
	}
	return effective
}

func llmFlagChanged(cmd *cobra.Command, name string) bool {
	if flag := cmd.Flags().Lookup(name); flag != nil && flag.Changed {
		return true
	}
	flag := cmd.InheritedFlags().Lookup(name)
	return flag != nil && flag.Changed
}

func Execute() {
	if exitCode, handled := handleRootStatusFlags(os.Args[1:]); handled {
		os.Exit(exitCode)
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func handleRootStatusFlags(args []string) (int, bool) {
	for _, arg := range args {
		switch arg {
		case "--version":
			return buildinfo.PrintStatus(os.Stdout, false), true
		case "--check-installed":
			return buildinfo.PrintStatus(os.Stdout, true), true
		case "--":
			return 0, false
		}

		if !strings.HasPrefix(arg, "-") {
			return 0, false
		}
	}
	return 0, false
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "数据目录 (默认 ~/.vocabmaster)")
	rootCmd.PersistentFlags().StringVar(&llmAdapter, "llm-adapter", "auto", "LLM adapter: auto, codex, claude, grok")
	rootCmd.PersistentFlags().StringVar(&llmModel, "llm-model", "", "LLM 模型 ID（需显式选择 adapter；留空使用默认值）")
	rootCmd.PersistentFlags().StringVar(&llmThinking, "llm-thinking", "", "LLM 思考程度（需显式选择 adapter；留空使用默认值）")
}
