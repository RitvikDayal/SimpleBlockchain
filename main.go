package main

import (
	"os"

	"github.com/ritvikdayal/SimpleBlockchain/cli"
)


func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()

}
