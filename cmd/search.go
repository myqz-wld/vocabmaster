package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/internal/ui"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "搜索单词",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		results := lib.Search(query)

		if len(results) == 0 {
			fmt.Printf("\n  未找到匹配 \"%s\" 的词汇。\n", query)
			return nil
		}

		fmt.Printf("\n  找到 %d 条结果：\n", len(results))
		for _, w := range results {
			linkedWords := lib.GetLinkedWords(w.ID)
			ui.DisplayWordCard(w, linkedWords)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
