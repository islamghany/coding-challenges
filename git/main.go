package main

import (
	"fmt"
	"mygit/cmd"
	"os"
)

func main() {
	cmder := cmd.NewCommand()

	if err := cmder.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
