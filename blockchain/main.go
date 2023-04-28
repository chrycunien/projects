package main

import "blockchain/pkg/blockchain"

func main() {
	bc := blockchain.NewBlockChain()
	defer bc.Close()

	cli := blockchain.NewCLI(bc)
	cli.Run()
}
