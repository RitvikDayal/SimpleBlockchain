/*
Simple block chain.
*/

package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Block struct {
	Timestamp time.Time `json:"timestamp"`
	Data      string    `json:"data"`
	PrevHash  []byte    `json:"prev_hash"`
	Hash      []byte    `json:"hash"`
}

func (block *Block) SetHash() {
	timestamp := []byte(block.Timestamp.String())
	headers := append(block.PrevHash, block.Data...)
	headers = append(headers, timestamp...)
	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}

func NewBlock(data string, prevHash []byte) *Block {
	block := &Block{time.Now(), data, prevHash, []byte{}}
	block.SetHash()
	return block
}

type Blockchain struct {
	Blocks []*Block `json:"blocks"`
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func PrintBlock(block *Block) {
	fmt.Println("Prev. hash: ", block.PrevHash)
	fmt.Println("Data: ", block.Data)
	fmt.Println("Hash: ", block.Hash)
	fmt.Println()
}

func saveBlockchain(bc *Blockchain) {
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

func main() {
	/*
		Infinite loop to ask user for input.
	*/
	bc := NewBlockchain()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// Ask user for Data to store in the block.
		fmt.Printf("Enter data to store in the block: ")
		var data string
		scanner.Scan()
		data = scanner.Text()

		// Add the block to the blockchain.
		bc.AddBlock(data)

		// Print number of Blocks in the blockchain.
		fmt.Println("Number of Blocks in the blockchain: ", len(bc.Blocks))

		// Ask if user what to print the blockchain.
		fmt.Printf("Do you want to print the blockchain? (y/n)")
		var print string
		scanner.Scan()
		print = scanner.Text()

		if print == "y" {
			for _, block := range bc.Blocks {
				PrintBlock(block)
			}
		}

		// Ask if user wants to exit.
		fmt.Printf("Do you want to exit? (y/n)")
		var exit string
		fmt.Scanln(&exit)
		if exit == "y" {
			// save the blockchain to a file.
			fmt.Println("Saving blockchain to a file...")
			saveBlockchain(bc)
			fmt.Println("Blockchain saved to a file.")
			break
		}
	}
}
