package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/internal/ui"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "查看学习统计",
	RunE: func(cmd *cobra.Command, args []string) error {
		lang, _ := cmd.Flags().GetString("lang")

		stats, err := db.GetStats(lang)
		if err != nil {
			return fmt.Errorf("获取统计信息失败: %w", err)
		}

		var accuracy float64
		if stats.TotalReviews > 0 {
			accuracy = float64(stats.CorrectCount) / float64(stats.TotalReviews) * 100
		}

		totalLib := lib.CountByLang(lang)

		ui.DisplayStatsTable(totalLib, &ui.StatsData{
			Learned:      stats.LearnedWords,
			Mastered:     stats.MasteredWords,
			Due:          stats.DueWords,
			TotalReviews: stats.TotalReviews,
			Accuracy:     accuracy,
			Streak:       stats.Streak,
		})
		return nil
	},
}

func init() {
	statsCmd.Flags().String("lang", "all", "语言: en, ja, all")
	rootCmd.AddCommand(statsCmd)
}
