package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/internal/model"
	"github.com/vocabmaster/vocabmaster/internal/ui"
)

var infoCmd = &cobra.Command{
	Use:   "info <word-id or text>",
	Short: "查看单词详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		w := lib.GetWord(query)
		if w == nil {
			results := lib.Search(query)
			if len(results) > 0 {
				w = results[0]
			}
		}

		if w == nil {
			fmt.Printf("\n  未找到词汇: %s\n", query)
			return nil
		}

		// Check for enriched version
		enriched, err := db.GetEnrichedWord(w.ID)
		if err == nil && enriched != nil {
			w = enriched
		}

		baseLinked := lib.GetLinkedWordsFor(w)
		linkedWords := make([]*model.Word, len(baseLinked))
		copy(linkedWords, baseLinked)
		for i, lw := range linkedWords {
			if e, err := db.GetEnrichedWord(lw.ID); err == nil && e != nil {
				linkedWords[i] = e
			}
		}
		ui.DisplayWordCard(w, linkedWords)

		// Show review progress
		record, err := db.GetReviewRecord(w.ID)
		if err == nil && record != nil {
			fmt.Printf("  %s学习进度：%s\n", "\033[1m", "\033[0m")
			fmt.Printf("  EF: %.2f  间隔: %d天  重复: %d次\n",
				record.EaseFactor, record.Interval, record.Repetitions)
			fmt.Printf("  总复习: %d次  正确: %d次\n",
				record.TotalReviews, record.CorrectCount)
			if !record.NextReviewAt.IsZero() {
				fmt.Printf("  下次复习: %s\n", record.NextReviewAt.Format("2006-01-02"))
			}
			fmt.Println()
		} else {
			fmt.Println("  尚未学习此词。")
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
