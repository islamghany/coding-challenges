package commands

import (
	"fmt"
	"os"
)

func (cmder *Commander) Pwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(dir)
}