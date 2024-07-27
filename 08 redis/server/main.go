package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"redis/foundation/enconder/resp"
	"redis/foundation/store"
	"redis/server/commands"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	// creating store
	store := store.NewStore()
	// creating commander
	commander := commands.NewCommander(store)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn, commander)
	}
}

func handleConnection(conn net.Conn, commander *commands.Commander) {
	defer func() {
		fmt.Println("closing the connection")
		conn.Close()
	}()
	// extract the data from the connection as array of bytes
	cmds := readFromConnection(conn)
	// trying to deserialize the cmds to resp format
	respData, err := resp.Deserialize(bufio.NewReader(bytes.NewReader(cmds)))
	// if an error occur when parsing return and close the connection
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}
	// check if the cmds is an array of respData objects
	dataArr, ok := respData.Data.([]resp.RESPData)
	errResp := resp.NewError(commands.InvalidArguments)
	if !ok || len(dataArr) == 0 {
		ret, _ := resp.Serialize(&errResp)
		conn.Write(ret)
		return
	}
	cmd := strings.ToLower(dataArr[0].Data.(string))
	var res resp.RESPData

	switch cmd {
	case "ping":
		res = commander.Ping(dataArr)
	case "echo":
		res = commander.Echo(dataArr)
	case "set":
		res = commander.Set(dataArr)
	case "get":
		res = commander.Get(dataArr)
	default:
		res = errResp
	}
	ret, err := resp.Serialize(&res)
	conn.Write(ret)
}

func readFromConnection(conn net.Conn) []byte {
	var buf bytes.Buffer
	for {
		temp := make([]byte, 1024)

		n, err := conn.Read(temp)
		if err != nil {
			break
		}
		buf.Write(temp[0:n])
		if n < 1024 {
			break
		}
	}
	return buf.Bytes()
}
