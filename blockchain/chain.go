package blockchain

/*
A Simple block chain implementation in Go.
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Blockchain struct {
	Blocks []*Block `json:"blocks"`
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func InitBlockChain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
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
