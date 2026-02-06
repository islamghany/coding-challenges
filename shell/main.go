package main

import (
	"fmt"
	"bufio"
	"os"
	"ccshell/commands"
	
)

func logangle(){
	fmt.Print(">> ")
}

func main() {
	cmder := commands.NewCommander()
	scanner := bufio.NewScanner(os.Stdin)
	logangle()
	for scanner.Scan(){
	cmd := scanner.Text()
	switch cmd {
		case "ls":
			cmder.Ls(".")
		case "cd":
			fmt.Println("cd")
		case "pwd":
			cmder.Pwd()
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Command not found")
		}
	logangle()
	}
}