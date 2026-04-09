package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/internal/model"
	"github.com/vocabmaster/vocabmaster/internal/session"
	"github.com/vocabmaster/vocabmaster/internal/ui"
)

var learnCmd = &cobra.Command{
	Use:   "learn",
	Short: "学习新词",
	RunE: func(cmd *cobra.Command, args []string) error {
		lang, _ := cmd.Flags().GetString("lang")
		level, _ := cmd.Flags().GetInt("level")
		count, _ := cmd.Flags().GetInt("count")

		cfg := session.Config{
			Lang:  lang,
			Level: model.DifficultyLevel(level),
			Count: count,
		}

		s := session.NewLearnSession(db, lib, cfg)
		result, err := s.Run()
		if err != nil {
			return err
		}

		ui.DisplaySessionSummary(0, result.Learned, result.Correct, result.Total)
		return nil
	},
}

func init() {
	learnCmd.Flags().String("lang", "all", "语言: en, ja, all")
	learnCmd.Flags().Int("level", 0, "难度: 1(初级), 2(中级), 3(高级), 0(全部)")
	learnCmd.Flags().Int("count", 5, "学习新词数量")
	rootCmd.AddCommand(learnCmd)
}
