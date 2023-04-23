package blockchain

import (
	"fmt"
	"strconv"
)

type BlockChain struct {
	blocks []*Block
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		blocks: []*Block{newGenesisBlock()},
	}
}

func newGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func (bc *BlockChain) AddBlock(data string) *BlockChain {
	block := NewBlock(data, bc.blocks[len(bc.blocks)-1].Hash)
	bc.blocks = append(bc.blocks, block)
	return bc
}

func (bc *BlockChain) Print() {
	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

	}
}
