package blockchain

import (
	"errors"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

func NewBlockChain() *BlockChain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic("cannot open bolt db")
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := newGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return errors.New("cannot create bolt bucket")
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				return errors.New("cannot put genesis bucket")
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				return errors.New("cannot update the last block hash")
			}

			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return err
	}); err != nil {
		panic(err)
	}

	return &BlockChain{
		tip: tip,
		db:  db,
	}
}

func newGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func (bc *BlockChain) AddBlock(data string) error {
	var lastHash []byte

	if err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	}); err != nil {
		return err
	}

	newBlock := NewBlock(data, lastHash)

	if err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if err := b.Put(newBlock.Hash, newBlock.Serialize()); err != nil {
			return err
		}
		if err := b.Put([]byte("l"), newBlock.Hash); err != nil {
			return err
		}
		bc.tip = newBlock.Hash
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (bc *BlockChain) Iterator() *BlockchainIterator {
	bci := BlockchainIterator{
		currentHash: bc.tip,
		db:          bc.db,
	}

	return &bci
}

func (bc *BlockChain) Print() {
	bci := bc.Iterator()

	var block *Block
	var err error

	for {
		if block, err = bci.Next(); err != nil {
			break
		}
		block.Print()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (bc *BlockChain) Close() {
	bc.db.Close()
}
