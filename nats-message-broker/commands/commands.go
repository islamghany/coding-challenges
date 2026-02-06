package commands

import (
	"fmt"
	"nats/parser"
	"net"
)

func HandleCommand(cmd *parser.Cmd, conn net.Conn) error {
	// Switch over the command name as a []byte to avoid string conversion
	switch cmd.Name {
	case parser.CONNECT:
		return connectCMD(conn)
	// case bytes.Equal(cmdName, []byte("CONNECT")):
	// 	return cmd.parseCONNECT(line[spaceIndex+1:]) // Pass JSON part only
	default:
		return fmt.Errorf("Unknown Command: %s", cmd.Name)
	}
}

func connectCMD(conn net.Conn) error {
	_, err := conn.Write([]byte("+OK\r\n"))
	return err
}
