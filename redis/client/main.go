package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"redis/foundation/enconder/resp"
)

func parseCommand(args []string) []string {
	cmds := make([]string, 0)
	for _, v := range args {
		if v != "" && v != " " {
			cmds = append(cmds, v)
		}
	}
	return cmds
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		println("Invalid number of arguments")
		os.Exit(1)
	}
	cmds := parseCommand(args)

	conn, err := net.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	respArr := resp.RESPData{
		Data: []resp.RESPData{},
		Type: resp.Array,
	}
	respArr.Type = resp.Array
	data := []resp.RESPData{}
	for _, v := range cmds {
		data = append(data,
			resp.RESPData{
				Data: v,
				Type: resp.BulkString,
			},
		)

	}
	respArr.Data = data
	darr, err := resp.Serialize(&respArr)

	conn.Write(darr)
	var buf bytes.Buffer
	_, err = io.Copy(&buf, conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(buf.String())
}
