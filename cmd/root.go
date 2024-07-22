package cmd

import (
	"log"

	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "swadloon",
	Run: func(cmd *cobra.Command, args []string) {
		pretty.Println(cmd)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
