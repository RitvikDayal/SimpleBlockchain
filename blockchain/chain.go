package blockchain

/*
A Simple block chain implementation in Go.
*/

import (
	"fmt"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath              = "./tmp/blocks"
	dbFile              = "./tmp/blocks/MANIFEST"
	genesisCoinbaseData = "Genesis coinbase"
)

type Blockchain struct {
	// Blocks []*Block `json:"blocks"`
	LastHash []byte `json:"lastHash"`
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (bc *Blockchain) AddBlock(txs []*Transaction) {
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

	newBlock := NewBlock(txs, lastHash)
	err = bc.Database.Update(func(txn *badger.Txn) error {

		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleError(err)

		err = txn.Set([]byte("lh"), newBlock.Hash)
		bc.LastHash = newBlock.Hash

		return err
	})
	HandleError(err)
}

func InitBlockChain(address string) *Blockchain {
	var lastHash []byte

	if DBexists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesisCoinbaseData)

		genesis := NewGenesisBlock(cbtx)
		err := txn.Set(genesis.Hash, genesis.Serialize())
		HandleError(err)

		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash

		return err
	})

	HandleError(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

func ContinueBlockChain(address string) *Blockchain {
	var lastHash []byte

	if !DBexists() {
		fmt.Println("No existing blockchain found")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	HandleError(err)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleError(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	HandleError(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

func (bc *Blockchain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.LastHash, bc.Database}
}
