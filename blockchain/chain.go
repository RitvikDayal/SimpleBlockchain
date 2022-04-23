package blockchain

/*
A Simple block chain implementation in Go.
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Block struct {
	Timestamp time.Time `json:"timestamp"`
	Data      []byte    `json:"data"`
	PrevHash  []byte    `json:"prev_hash"`
	Hash      []byte    `json:"hash"`
	Nonce     int       `json:"nonce"`
}

type Blockchain struct {
	Blocks []*Block `json:"blocks"`
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

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func InitBlockChain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func PrintBlock(block *Block) {
	fmt.Println("Prev. hash: ", block.PrevHash)
	fmt.Println("Data: ", block.Data)
	fmt.Println("Hash: ", block.Hash)
	fmt.Println()
}

func SaveBlockchain(bc *Blockchain) {
	data, err := json.MarshalIndent(bc, "", " ")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	err = ioutil.WriteFile("blockchain.json", data, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}
