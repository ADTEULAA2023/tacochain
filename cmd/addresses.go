package cmd

import (
	"fmt"

	"github.com/ADTEULAA2023/tacochain/internal/wallet"
	"github.com/spf13/cobra"
)

var (
	listAddresses = &cobra.Command{
		Use:   "list",
		Short: `Lists all of the addresses in the blockchain`,
		RunE: func(cmd *cobra.Command, args []string) error {
			wallets, _ := wallet.CreateWallets()
			addresses := wallets.GetAllAddresses()

			for _, address := range addresses {
				fmt.Println(address)
			}

			return nil
		},
	}
)

func newAddressesCMD() *cobra.Command {
	addressesCMD := &cobra.Command{
		Use:   "chains",
		Short: "Taco chain is a blockchain client used to interact in the taco network",
	}

	addressesCMD.AddCommand(listAddresses)

	return addressesCMD
}
