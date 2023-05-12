package cmd

import (
	"fmt"
	"strconv"

	"github.com/ADTEULAA2023/tacochain/internal/blockchain"
	"github.com/spf13/cobra"
)

var (
	listChainsCMD = &cobra.Command{
		Use:   "list",
		Short: `Lists all of the chains in the blockchain`,
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := blockchain.ContinueBlockChain("")
			if err != nil {
				return err
			}
			defer chain.Database.Close()
			iterator := chain.Iterator()

			for {
				block, err := iterator.Next()
				if err != nil {
					return err
				}

				fmt.Printf("Previous hash: %x\n", block.PrevHash)
				fmt.Printf("hash: %x\n", block.Hash)
				pow := blockchain.NewProofOfWork(block)
				fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
				fmt.Println()
				// This works because the Genesis block has no PrevHash to point to.
				if len(block.PrevHash) == 0 {
					break
				}
			}

			return nil
		},
	}

	createChainCMD = &cobra.Command{
		Use:   "create",
		Short: `Create new chain in the blockchain`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("must specify address")
			}

			newChain, err := blockchain.InitBlockChain(args[0])
			if err != nil {
				return err
			}
			newChain.Database.Close()
			fmt.Println("Finished creating chain")

			return nil
		},
	}
)

func newChainsCMD() *cobra.Command {
	chainCMD := &cobra.Command{
		Use:     "chains",
		Aliases: []string{"chain", "c"},
		Short:   "Addresses all chains operations",
	}

	chainCMD.AddCommand(createChainCMD)
	chainCMD.AddCommand(listChainsCMD)

	return chainCMD
}
