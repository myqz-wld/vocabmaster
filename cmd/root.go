package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/internal/library"
	"github.com/vocabmaster/vocabmaster/internal/store"
)

var (
	dataDir string
	db      *store.SQLiteStore
	lib     *library.Library
)

var rootCmd = &cobra.Command{
	Use:   "vocabmaster",
	Short: "命令行背单词工具 - 支持英文和日文",
	Long:  "VocabMaster 是一个命令行背单词工具，支持英文和日文单词学习，\n基于 SM-2 间隔重复算法，带有中文释义、发音标注和例句。",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" || cmd.Name() == "version" {
			return nil
		}

		if dataDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("获取用户目录失败: %w", err)
			}
			dataDir = filepath.Join(home, ".vocabmaster")
		}

		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return fmt.Errorf("创建数据目录失败: %w", err)
		}

		var err error
		db, err = store.NewSQLiteStore(filepath.Join(dataDir, "vocabmaster.db"))
		if err != nil {
			return fmt.Errorf("打开数据库失败: %w", err)
		}

		lib, err = library.NewLibrary()
		if err != nil {
			return fmt.Errorf("加载词库失败: %w", err)
		}

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if db != nil {
			db.Close()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "数据目录 (默认 ~/.vocabmaster)")
}
