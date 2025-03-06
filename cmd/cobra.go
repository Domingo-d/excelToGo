package cmd

import (
	"errors"
	"excelToGo/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "Version",
		Short: "Excel to Go",
		Long:  `Excel to Go`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				tip()
				return errors.New(color.RedString("参数错误!"))
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) { tip() },

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func tip() {
	color.Green(`欢迎使用 Excel to Go v` + common.Version)
}

func init() {
	rootCmd.AddCommand(InteractiveCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("执行错误", err)
		os.Exit(-1)
	}
}
