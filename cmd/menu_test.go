package cmd

import (
	"github.com/spf13/cobra"
	"testing"
)

func TestMenu(t *testing.T) {
	ExecuteCommand(rootCmd, "console")
}

func ExecuteCommand(root *cobra.Command, args ...string) error {
	root.SetArgs(args)
	err := root.Execute()
	return err
}
