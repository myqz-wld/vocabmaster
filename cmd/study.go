package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/internal/model"
	"github.com/vocabmaster/vocabmaster/internal/session"
)

var studyCmd = &cobra.Command{
	Use:   "study",
	Short: "一键学习（自动平衡复习和新词）",
	RunE: func(cmd *cobra.Command, args []string) error {
		lang, _ := cmd.Flags().GetString("lang")
		level, _ := cmd.Flags().GetInt("level")
		newWords, _ := cmd.Flags().GetInt("new-words")

		cfg := session.Config{
			Lang:     lang,
			Level:    model.DifficultyLevel(level),
			NewWords: newWords,
		}

		s := session.NewStudySession(db, lib, cfg)
		return s.Run()
	},
}

func init() {
	studyCmd.Flags().String("lang", "all", "语言: en, ja, all")
	studyCmd.Flags().Int("level", 0, "难度: 1(初级), 2(中级), 3(高级), 0(全部)")
	studyCmd.Flags().Int("new-words", 0, "新词数量上限 (0=自动)")
	rootCmd.AddCommand(studyCmd)
}
