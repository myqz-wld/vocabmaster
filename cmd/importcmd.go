package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/internal/model"
)

var importCmd = &cobra.Command{
	Use:   "import <file.json>",
	Short: "导入外部词库",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("读取文件失败: %w", err)
		}

		var file model.WordLibraryFile
		if err := json.Unmarshal(data, &file); err != nil {
			return fmt.Errorf("解析 JSON 失败: %w", err)
		}

		if len(file.Words) == 0 {
			fmt.Println("  文件中没有词汇数据。")
			return nil
		}

		imported := 0
		for _, w := range file.Words {
			if err := db.SaveEnrichedWord(&w); err != nil {
				fmt.Printf("  导入 %s 失败: %v\n", w.ID, err)
				continue
			}
			imported++
		}

		fmt.Printf("\n  成功导入 %d 个词汇。\n", imported)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
