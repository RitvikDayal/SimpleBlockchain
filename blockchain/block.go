package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

type Block struct {
	Timestamp    time.Time      `json:"timestamp"`
	Transactions []*Transaction `json:"transactions"`
	PrevHash     []byte         `json:"prev_hash"`
	Hash         []byte         `json:"hash"`
	Nonce        int            `json:"nonce"`
}

func (block *Block) HashTxs() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{time.Now(), txs, prevHash, []byte{}, 0}

	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func PrintBlock(block *Block) {
	fmt.Println("Prev. hash: ", block.PrevHash)
	fmt.Println("Data: ", block.Transactions)
	fmt.Println("Hash: ", block.Hash)
	fmt.Println()
}

func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	HandleError(err)

	return result.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	HandleError(err)

	return &block
}
