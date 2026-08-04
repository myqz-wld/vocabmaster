package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/src/internal/appconfig"
	"github.com/vocabmaster/vocabmaster/src/internal/llm"
)

var (
	configAdapter  string
	configModel    string
	configThinking string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "查看或修改本地配置",
	Args:  cobra.NoArgs,
	RunE:  showConfig,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "显示本地配置",
	Args:  cobra.NoArgs,
	RunE:  showConfig,
}

var configSetLLMCmd = &cobra.Command{
	Use:   "set-llm",
	Short: "保存默认 LLM adapter、模型和思考程度",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := llm.NormalizeOptions(llm.Options{
			Adapter:  configAdapter,
			Model:    configModel,
			Thinking: configThinking,
		})
		if err != nil {
			return err
		}

		settings.LLM = appconfig.LLMSettings{
			Adapter:  options.Adapter,
			Model:    options.Model,
			Thinking: options.Thinking,
		}
		if err := appconfig.Save(dataDir, settings); err != nil {
			return fmt.Errorf("保存本地配置失败: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "已保存 LLM 配置: adapter=%s, model=%s, thinking=%s\n", options.Adapter, displayLLMDefault(options.Model), displayLLMDefault(options.Thinking))
		fmt.Fprintf(cmd.OutOrStdout(), "配置文件: %s\n", appconfig.Path(dataDir))
		fmt.Fprintln(cmd.OutOrStdout(), "命令行 --llm-* 参数仍可临时覆盖这些默认值。")
		return nil
	},
}

func showConfig(cmd *cobra.Command, args []string) error {
	options, err := llm.NormalizeOptions(llm.Options{
		Adapter:  settings.LLM.Adapter,
		Model:    settings.LLM.Model,
		Thinking: settings.LLM.Thinking,
	})
	if err != nil {
		return fmt.Errorf("本地 LLM 配置无效: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "本地 LLM 配置:")
	fmt.Fprintf(cmd.OutOrStdout(), "  adapter: %s\n", options.Adapter)
	fmt.Fprintf(cmd.OutOrStdout(), "  model: %s\n", displayLLMDefault(options.Model))
	fmt.Fprintf(cmd.OutOrStdout(), "  thinking: %s\n", displayLLMDefault(options.Thinking))
	fmt.Fprintf(cmd.OutOrStdout(), "  文件: %s\n", appconfig.Path(dataDir))
	return nil
}

func displayLLMDefault(value string) string {
	if value == "" {
		return "(使用 adapter 默认值)"
	}
	return value
}

func isConfigCommand(cmd *cobra.Command) bool {
	return cmd == configCmd || cmd.Parent() == configCmd
}

func init() {
	configSetLLMCmd.Flags().StringVar(&configAdapter, "adapter", "auto", "LLM adapter: auto, codex, claude, grok")
	configSetLLMCmd.Flags().StringVar(&configModel, "model", "", "LLM 模型 ID（留空使用 adapter 默认值）")
	configSetLLMCmd.Flags().StringVar(&configThinking, "thinking", "", "LLM 思考程度（留空使用 adapter 默认值）")
	configCmd.AddCommand(configShowCmd, configSetLLMCmd)
	rootCmd.AddCommand(configCmd)
}
