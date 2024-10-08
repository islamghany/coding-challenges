package commands

import (
	"fmt"
	"os"
)

func (cmder *Commander) Ls(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	files, err := file.Readdir(0)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
	
}