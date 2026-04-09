package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/internal/library"
	"github.com/vocabmaster/vocabmaster/internal/llm"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "批量生成/矫正词库数据（调用 LLM）",
	RunE: func(cmd *cobra.Command, args []string) error {
		lang, _ := cmd.Flags().GetString("lang")
		count, _ := cmd.Flags().GetInt("count")
		force, _ := cmd.Flags().GetBool("force")

		if !llm.IsAvailable() {
			fmt.Println("  错误：未找到 claude CLI。请确保已安装 Claude Code。")
			return nil
		}

		words := lib.FilterByLang(lang)
		if count > 0 && len(words) > count {
			words = words[:count]
		}

		allWordIDs := getAllLibWordIDs(lib)

		processed := 0
		skipped := 0
		failed := 0

		for i, w := range words {
			// Skip if already enriched (unless force)
			if !force {
				existing, err := db.GetEnrichedWord(w.ID)
				if err == nil && existing != nil {
					fmt.Printf("  [%d/%d] 跳过已增强: %s\n", i+1, len(words), w.Text)
					skipped++
					continue
				}
			}

			fmt.Printf("  [%d/%d] 正在处理: %s ...", i+1, len(words), w.Text)

			enriched, err := llm.EnrichWord(w, allWordIDs)
			if err != nil {
				fmt.Printf(" 失败: %v\n", err)
				failed++
				continue
			}

			if err := db.SaveEnrichedWord(enriched); err != nil {
				fmt.Printf(" 保存失败: %v\n", err)
				failed++
				continue
			}

			fmt.Println(" 完成")
			processed++
		}

		fmt.Printf("\n  处理完成：成功 %d，跳过 %d，失败 %d\n", processed, skipped, failed)
		return nil
	},
}

func getAllLibWordIDs(l *library.Library) []string {
	words := l.AllWords()
	ids := make([]string, len(words))
	for i, w := range words {
		ids[i] = w.ID
	}
	return ids
}

func init() {
	generateCmd.Flags().String("lang", "all", "语言: en, ja, all")
	generateCmd.Flags().Int("count", 0, "处理数量 (0=全部)")
	generateCmd.Flags().Bool("force", false, "强制重新生成已有数据")
	rootCmd.AddCommand(generateCmd)
}
