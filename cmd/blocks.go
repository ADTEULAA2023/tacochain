package cmd

import (
	"github.com/spf13/cobra"
)

var addBlock = &cobra.Command{
	Use: "blocks",

	Short: "Taco chain is a blockchain client used to interact in the taco network",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func newBlocksCMD() *cobra.Command {
	blocksCMD := &cobra.Command{
		Use:     "blocks",
		Aliases: []string{"block", "b"},
		Short:   "Addresses all chains operations",
	}

	blocksCMD.AddCommand(addBlock)
	return blocksCMD
}
