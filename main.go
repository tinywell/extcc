package main

import (
	"extcc/cmds"

	"github.com/spf13/cobra"
)

// RootCmd ...
var RootCmd = &cobra.Command{
	Use:   "extcc",
	Short: "external chaincode client，用于调用外部部署链码，方便测试链码逻辑",
}

func main() {
	RootCmd.AddCommand(cmds.InvokeCmd)
	RootCmd.AddCommand(cmds.RegCmd)
	RootCmd.AddCommand(cmds.ServerCmd)
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}
