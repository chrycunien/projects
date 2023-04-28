package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().UnixNano(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Nonce = nonce
	block.Hash = hash

	return block
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	_ = encoder.Encode(b)

	return result.Bytes()
}

func (b *Block) Print() {
	fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
	fmt.Printf("Data: %s\n", b.Data)
	fmt.Printf("Hash: %x\n", b.Hash)
	pow := NewProofOfWork(b)
	fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	fmt.Println()
}

func DeserializeBlock(b []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(b))
	_ = decoder.Decode(&block)

	return &block
}
