package main

/*
Simple block chain.
*/

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	// Import local block chain package.
	"github.com/ritvikdayal/SimpleBlockchain/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.Blockchain
}

func (cli *CommandLine) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -data DATA -- will add a new block to the chain with the data provided.")
	fmt.Println(" print -- will print all the blocks in the chain.")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Block added to the blockchain.")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()
	for {
		block := iter.Next()
		if block == nil {
			break
		}
		fmt.Printf("\n----------------------------\n")
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("\nPoW: %s\n", strconv.FormatBool(pow.Validate()))

		fmt.Printf("----------------------------\n")

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	default:
		cli.PrintUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		} else {
			cli.addBlock(*addBlockData)
		}
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

func main() {
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.Run()

}
