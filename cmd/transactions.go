package cmd

import (
	"fmt"
	"strconv"

	"github.com/ADTEULAA2023/tacochain/internal/blockchain"
	"github.com/spf13/cobra"
)

var (
	sendCMD = &cobra.Command{
		Use:   "send",
		Short: `Sends information to address`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("must specify from, to, and amount")
			}

			from, to := args[0], args[1]
			amount, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("amount must be a number: %w", err)
			}

			chain := blockchain.ContinueBlockChain(from)
			defer chain.Database.Close()

			tx := blockchain.NewTransaction(from, to, amount, chain)

			chain.AddBlock([]*blockchain.Transaction{tx})
			fmt.Println("Success!")

			return nil
		},
	}
)

func newTransactionsCMD() *cobra.Command {
	addressesCMD := &cobra.Command{
		Use:   "transactions",
		Short: "Taco chain is a blockchain client used to interact in the taco network",
	}

	addressesCMD.AddCommand(listAddresses)

	return addressesCMD
}
