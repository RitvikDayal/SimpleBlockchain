package blockchain

import (
	"fmt"
	"time"
)

type Block struct {
	Timestamp time.Time `json:"timestamp"`
	Data      []byte    `json:"data"`
	PrevHash  []byte    `json:"prev_hash"`
	Hash      []byte    `json:"hash"`
	Nonce     int       `json:"nonce"`
}

func NewBlock(data string, prevHash []byte) *Block {
	b_data := []byte(data)
	block := &Block{time.Now(), b_data, prevHash, []byte{}, 0}

	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func PrintBlock(block *Block) {
	fmt.Println("Prev. hash: ", block.PrevHash)
	fmt.Println("Data: ", block.Data)
	fmt.Println("Hash: ", block.Hash)
	fmt.Println()
}