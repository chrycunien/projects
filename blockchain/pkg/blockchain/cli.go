package blockchain

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	CommandAddBlock   = "addblock"
	CommandPrintChain = "printchain"
)

type CLI struct {
	bc *BlockChain
}

func NewCLI(bc *BlockChain) *CLI {
	return &CLI{
		bc: bc,
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet(CommandAddBlock, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(CommandPrintChain, flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block Data")

	switch os.Args[1] {
	case CommandAddBlock:
		addBlockCmd.Parse(os.Args[2:])
	case CommandPrintChain:
		printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	switch {
	case addBlockCmd.Parsed():
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	case printChainCmd.Parsed():
		cli.printChain()
	}
}

func (cli *CLI) printUsage() {
	log.Println("Usage:")
	log.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	log.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string) {
	if err := cli.bc.AddBlock(data); err != nil {
		fmt.Println("Success!")
	}
}

func (cli *CLI) printChain() {
	cli.bc.Print()
}
