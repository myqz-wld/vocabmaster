package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/internal/ui"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "重置学习进度（恢复为未学习状态）",
	RunE: func(cmd *cobra.Command, args []string) error {
		word, _ := cmd.Flags().GetString("word")
		all, _ := cmd.Flags().GetBool("all")

		if word != "" {
			// 支持 word ID 或 text 搜索，仿照 info 命令
			wordID := word
			if lib.GetWord(word) == nil {
				results := lib.Search(word)
				if len(results) == 1 {
					wordID = results[0].ID
				} else if len(results) > 1 {
					fmt.Printf("  找到多个匹配词汇，请指定精确 ID：\n")
					for _, r := range results {
						fmt.Printf("    %s (%s)\n", r.ID, r.Text)
					}
					return nil
				}
			}

			if !ui.ReadConfirmation(fmt.Sprintf("确定要重置 %s 的学习进度吗？该词将恢复为未学习状态。", wordID)) {
				fmt.Println("  已取消。")
				return nil
			}
			record, err := db.GetReviewRecord(wordID)
			if err != nil || record == nil {
				fmt.Printf("  未找到 %s 的学习记录。请使用词汇 ID（如 en_apple）。\n", word)
				return nil
			}
			if err := db.ResetWord(wordID); err != nil {
				return err
			}
			fmt.Printf("  已重置 %s 的学习进度。\n", wordID)
			return nil
		}

		if all {
			if !ui.ReadConfirmation("确定要重置所有学习进度吗？所有词将恢复为未学习状态，此操作不可撤销！") {
				fmt.Println("  已取消。")
				return nil
			}
			records, err := db.GetAllReviewRecords()
			if err != nil {
				return err
			}
			if err := db.ResetAll(); err != nil {
				return err
			}
			fmt.Printf("  已重置所有 %d 个词的学习进度。\n", len(records))
			return nil
		}

		fmt.Println("  请指定 --word <id> 或 --all")
		return nil
	},
}

func init() {
	resetCmd.Flags().String("word", "", "重置特定词汇（支持 ID 或单词文本）")
	resetCmd.Flags().Bool("all", false, "重置所有进度")
	rootCmd.AddCommand(resetCmd)
}
