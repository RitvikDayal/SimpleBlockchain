package main

/*
Simple block chain.
*/

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	// Import local block chain package.
	"github.com/ritvikdayal/SimpleBlockchain/blockchain"
)

func main() {
	/*
		Infinite loop to ask user for input.
	*/
	bc := blockchain.InitBlockChain()
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
				blockchain.PrintBlock(block)

				pow := blockchain.NewProofOfWork(block)
				fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
			}
		}

		// Ask if user wants to exit.
		fmt.Printf("Do you want to exit? (y/n)")
		var exit string
		fmt.Scanln(&exit)
		if exit == "y" {
			// save the blockchain to a file.
			fmt.Println("Saving blockchain to a file...")
			blockchain.SaveBlockchain(bc)
			fmt.Println("Blockchain saved to a file.")
			break
		}
	}
}
