package blockchain

import (
	"log"

	"github.com/dgraph-io/badger"
)

//  Contains utility functions for working with blocks.

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func (bci *BlockChainIterator) Next() *Block {
	var block *Block

	err := bci.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(bci.CurrentHash)
		HandleError(err)
		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		return err
	})
	HandleError(err)

	bci.CurrentHash = block.PrevHash

	return block
}

func (bc *Blockchain) Print() {
	iter := bc.Iterator()
	for {
		block := iter.Next()
		if block == nil {
			break
		}
		PrintBlock(block)
	}
}
