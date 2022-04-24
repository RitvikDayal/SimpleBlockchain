package blockchain

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/dgraph-io/badger"
)

//  Contains utility functions for working with blocks.

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
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

func (bc *Blockchain) UnspentTransactions(address string) []Transaction {
	var UTXOs []Transaction
	iter := bc.Iterator()
	STXOs := make(map[string][]int)

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if STXOs[txID] != nil {
					for _, out := range STXOs[txID] {
						if out == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlockedWith(address) {
					UTXOs = append(UTXOs, *tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					if in.CanUnlockIn(address) {
						inTxId := hex.EncodeToString(in.Txid)
						STXOs[inTxId] = append(STXOs[inTxId], in.Vout)
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return UTXOs
}
