package main

import "net"

func main() {
	conn, err := net.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.Write([]byte("Hello, World!"))
	buf := make([]byte, 1024)
	conn.Read(buf)
	println("text", string(buf))

}
