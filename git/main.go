package main

import (
	"flag"
	"fmt"
	"mygit/cmd"
	"os"
)

func main() {

	flag.Parse()

	args := flag.Args()

	cmder := cmd.NewCommand()

	err := cmder.Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
