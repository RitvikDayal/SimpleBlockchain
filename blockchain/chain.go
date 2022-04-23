package blockchain

/*
A Simple block chain implementation in Go.
*/

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const dbFile = "./tmp/blocks"

type Blockchain struct {
	// Blocks []*Block `json:"blocks"`
	LastHash []byte `json:"lastHash"`
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleError(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	HandleError(err)

	newBlock := NewBlock(data, lastHash)
	err = bc.Database.Update(func(txn *badger.Txn) error {

		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleError(err)

		err = txn.Set([]byte("lh"), newBlock.Hash)
		bc.LastHash = newBlock.Hash

		return err
	})
	HandleError(err)
}

func InitBlockChain() *Blockchain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbFile)

	db, err := badger.Open(opts)
	HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := NewGenesisBlock()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			HandleError(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			HandleError(err)
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			return err
		}
	})

	HandleError(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}


func (bc *Blockchain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.LastHash, bc.Database}
}