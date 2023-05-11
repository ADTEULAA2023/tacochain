package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tacochain",
	Short: "Taco chain is a blockchain client used to interact in the taco network",
}

func Execute() {
	rootCmd.AddCommand(newChainsCMD())
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
