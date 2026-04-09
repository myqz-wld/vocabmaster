package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/internal/session"
	"github.com/vocabmaster/vocabmaster/internal/ui"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "复习到期单词",
	RunE: func(cmd *cobra.Command, args []string) error {
		lang, _ := cmd.Flags().GetString("lang")
		count, _ := cmd.Flags().GetInt("count")

		cfg := session.Config{
			Lang:  lang,
			Count: count,
		}

		s := session.NewReviewSession(db, lib, cfg)
		result, err := s.Run()
		if err != nil {
			return err
		}

		ui.DisplaySessionSummary(result.Reviewed, 0, result.Correct, result.Total)
		return nil
	},
}

func init() {
	reviewCmd.Flags().String("lang", "all", "语言: en, ja, all")
	reviewCmd.Flags().Int("count", 20, "最多复习词数 (0=全部)")
	rootCmd.AddCommand(reviewCmd)
}
