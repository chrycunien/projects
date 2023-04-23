package main

import (
	"blockchain/pkg/blockchain"
)

func main() {
	bc := blockchain.NewBlockChain()
	bc.AddBlock("Send 1 BTC to Eric").AddBlock("Send 2 more BTC to Eric")
	bc.Print()
}
