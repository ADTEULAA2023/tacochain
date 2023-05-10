package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/ADTEULAA2023/tacochain/internal/blockchain"
)

const Difficulty = 12
const dbPath = `E:\Code\go\src\github.com\ADTEULAA2023\tacochain\tmp\blocks`

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

//printUsage will display what options are availble to the user
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println(" add -block <BLOCK_DATA> - add a block to the chain")
	fmt.Println(" print - prints the blocks in the chain")
}

//validateArgs ensures the cli was given valid input
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		//go exit will exit the application by shutting down the goroutine
		// if you were to use os.exit you might corrupt the data
		runtime.Goexit()
	}
}

//addBlock allows users to add blocks to the chain via the cli
func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data, Difficulty)
	fmt.Println("Added Block!")
}

//printChain will display the entire contents of the blockchain
func (cli *CommandLine) printChain() {
	iterator := cli.blockchain.Iterator()

	for {
		block, err := iterator.Next()
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Printf("hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block, Difficulty)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		// This works because the Genesis block has no PrevHash to point to.
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

//run will start up the command line
func (cli *CommandLine) run() error {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			return err
		}

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			return err
		}

	default:
		cli.printUsage()
		runtime.Goexit()
	}
	// Parsed() will return true if the object it was used on has been called
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}

	return nil
}

func main() {
	chain, err := blockchain.InitBlockChain(dbPath, Difficulty)
	if err != nil {
		os.Exit(1)
	}

	defer chain.Database.Close()
	defer os.Exit(0)

	cli := CommandLine{chain}

	cli.run()
}
