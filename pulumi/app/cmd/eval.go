/*
Copyright © 2024 NAME HERE the0x@the0x.tech
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// evalCmd represents the eval command
var evalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Send a request to bedrock ",
	Long:  `This will send a request to AWS bedrock and return a response.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Send me something, I'm listening!")
	},
}

func init() {
	rootCmd.AddCommand(evalCmd)
}
