package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vocabmaster/vocabmaster/src/internal/buildinfo"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示安装版本和构建信息",
	Run: func(cmd *cobra.Command, args []string) {
		buildinfo.PrintStatus(cmd.OutOrStdout(), false)
	},
}

var checkInstalledCmd = &cobra.Command{
	Use:   "check-installed",
	Short: "检查已安装命令是否匹配当前源码 checkout",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(buildinfo.PrintStatus(cmd.OutOrStdout(), true))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(checkInstalledCmd)
}
