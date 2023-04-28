package blockchain

import "github.com/boltdb/bolt"

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bci *BlockchainIterator) Next() (*Block, error) {
	var block *Block

	if err := bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(bci.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	}); err != nil {
		return nil, err
	}

	bci.currentHash = block.PrevBlockHash

	return block, nil
}
