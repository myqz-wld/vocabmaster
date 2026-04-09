package cmd

import (
	"fmt"

	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/internal/model"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "浏览词库",
	RunE: func(cmd *cobra.Command, args []string) error {
		lang, _ := cmd.Flags().GetString("lang")
		level, _ := cmd.Flags().GetInt("level")
		statusFilter, _ := cmd.Flags().GetString("status")
		limit, _ := cmd.Flags().GetInt("limit")

		words := lib.FilterByLang(lang)
		words = lib.FilterByDifficulty(words, model.DifficultyLevel(level))

		learnedSet, err := db.GetLearnedWordSet()
		if err != nil {
			return err
		}

		var filtered []*model.Word
		for _, w := range words {
			switch statusFilter {
			case "learned":
				if !learnedSet[w.ID] {
					continue
				}
			case "unlearned":
				if learnedSet[w.ID] {
					continue
				}
			}
			filtered = append(filtered, w)
			if limit > 0 && len(filtered) >= limit {
				break
			}
		}

		if len(filtered) == 0 {
			fmt.Println("\n  没有符合条件的词汇。")
			return nil
		}

		// Simple table output
		fmt.Println()
		header := fmt.Sprintf("  %-20s %-12s %-14s %-18s %-6s %-4s",
			"ID", "单词", "发音", "中文释义", "难度", "状态")
		fmt.Println(header)
		fmt.Println("  " + repeatDash(80))

		for _, w := range filtered {
			st := "未学"
			if learnedSet[w.ID] {
				st = "已学"
			}
			fmt.Printf("  %-20s %-12s %-14s %-18s %-6s %-4s\n",
				padRight(w.ID, 20),
				padRight(w.Text, 12),
				padRight(w.Pronunciation, 14),
				padRight(truncate(w.ChineseDef, 15), 18),
				padRight(w.Difficulty.String(), 6),
				st,
			)
		}

		fmt.Printf("\n  共 %d 条结果\n", len(filtered))
		return nil
	},
}

func init() {
	listCmd.Flags().String("lang", "all", "语言: en, ja, all")
	listCmd.Flags().Int("level", 0, "难度: 1, 2, 3, 0(全部)")
	listCmd.Flags().String("status", "all", "状态: all, learned, unlearned")
	listCmd.Flags().Int("limit", 20, "最多显示条数")
	rootCmd.AddCommand(listCmd)
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

func padRight(s string, width int) string {
	sw := runewidth.StringWidth(s)
	if sw >= width {
		return s
	}
	padding := width - sw
	result := s
	for i := 0; i < padding; i++ {
		result += " "
	}
	return result
}

func repeatDash(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "-"
	}
	return s
}
