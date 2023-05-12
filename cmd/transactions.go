package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ADTEULAA2023/tacochain/internal/blockchain"
	"github.com/ADTEULAA2023/tacochain/pkg"
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
			studyData, err := ioutil.ReadFile(args[2])
			if err != nil {
				return fmt.Errorf("amount must be a number: %w", err)
			}

			chain, err := blockchain.ContinueBlockChain(from)
			if err != nil {
				return err
			}
			defer chain.Database.Close()

			tx, err := blockchain.NewTransaction(from, to, studyData, chain)
			if err != nil {
				return err
			}

			chain.AddBlock([]*blockchain.Transaction{tx})
			fmt.Println("Success!")

			return nil
		},
	}

	listRecordsCMD = &cobra.Command{
		Use:   "list",
		Short: `Lists all of the transactions made to the wallet`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("must specify address")
			}

			address := args[0]
			chain, err := blockchain.ContinueBlockChain(address)
			if err != nil {
				return err
			}
			defer chain.Database.Close()

			UTXOs, err := chain.FindUTXO(address)
			if err != nil {
				return err
			}

			fmt.Printf("Studies on: %s\n", address)
			for _, out := range UTXOs {
				fmt.Println("===== Public key:", string(out.Address))
				fmt.Printf("%s\n", string(out.Data))
			}

			return nil
		},
	}

	readTransactionCMD = &cobra.Command{
		Use:     "read ADDRESS TX_PUB_KEY TX_PRIV_KEY",
		Short:   `Reads a transaction from the blockchain and decodes it`,
		Aliases: []string{"r"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("must specify address, pubKey, and privKey")
			}

			address := args[0]
			txPubKey := args[1]
			txPrivKey := args[2]

			chain, err := blockchain.ContinueBlockChain(address)
			if err != nil {
				return err
			}
			defer chain.Database.Close()

			UTXOs, err := chain.FindUTXO(address)
			if err != nil {
				return err
			}

			for _, out := range UTXOs {
				if txPubKey == out.TxPubKey {
					data, err := pkg.DecodeTransactionData(txPrivKey, txPubKey, string(out.Data))
					if err != nil {
						return err
					}

					log.Println(data)
					return nil
				}
			}

			return fmt.Errorf("could not find transaction")
		},
	}
)

func newTransactionsCMD() *cobra.Command {
	transactionsCMD := &cobra.Command{
		Use:     "transactions",
		Aliases: []string{"transaction", "tx"},
		Short:   "Taco chain is a blockchain client used to interact in the taco network",
	}

	transactionsCMD.AddCommand(sendCMD)
	transactionsCMD.AddCommand(listRecordsCMD)
	transactionsCMD.AddCommand(readTransactionCMD)

	return transactionsCMD
}
