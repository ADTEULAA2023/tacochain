package cmd

import (
	"fmt"
	"os"

	"github.com/ADTEULAA2023/tacochain/internal/wallet"
	"github.com/spf13/cobra"
)

var (
	createCMD = &cobra.Command{
		Use:   "create",
		Short: `Creates a new wallet in the blockchain`,
		RunE: func(cmd *cobra.Command, args []string) error {
			wallets, err := wallet.CreateWallets()
			if err != nil && !os.IsNotExist(err) {
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
		Use:     "wallets",
		Short:   "Wallet operations",
		Aliases: []string{"wallet", "w"},
	}

	walletCMD.AddCommand(createCMD)

	return walletCMD
}
