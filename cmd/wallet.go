package cmd

import (
	"fmt"

	"github.com/ADTEULAA2023/tacochain/internal/wallet"
	"github.com/spf13/cobra"
)

var (
	createCMD = &cobra.Command{
		Short: `Creates a new wallet in the blockchain`,
		RunE: func(cmd *cobra.Command, args []string) error {
			wallets, err := wallet.CreateWallets()
			if err != nil {
				return err
			}
			address := wallets.AddWallet()
			wallets.SaveFile()

			fmt.Printf("New address is: %s\n", address)
			return nil
		},
	}
)

func newWalletCMD() *cobra.Command {
	walletCMD := &cobra.Command{
		Use:   "chains",
		Short: "Wallet operations",
	}

	walletCMD.AddCommand(createCMD)

	return walletCMD
}
