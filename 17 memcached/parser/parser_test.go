package commandsparser

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNewCommandsParser(t *testing.T) {
	cp := NewParser()
	if cp == nil {
		t.Error("Expected a new CommandsParser")
	}
	reader := bytes.NewReader([]byte("set key 0 0 5\r\nvalue\r\n"))
	cmd, err := cp.Parse(reader)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	fmt.Printf("cmd: %+v, err: %v\n", cmd, err)
	if cmd.Name != SetCommand {
		t.Errorf("Expected Name to be %s, got %s", SetCommand, cmd.Name)
	}
	if cmd.Key != "key" {
		t.Error("Expected key to be key")
	}
	if cmd.Flags != 0 {
		t.Error("Expected flags to be 0")
	}
	if cmd.Expiry != 0 {
		t.Error("Expected expiry to be 0")
	}
	if cmd.Bytes != 5 {
		t.Error("Expected bytes to be 5")
	}
	if string(cmd.Value) != "value" {
		t.Error("Expected value to be value")
	}

}

func TestGetCommand(t *testing.T) {
	cp := NewParser()
	reader := bytes.NewReader([]byte("get key\r\n"))
	cmd, err := cp.Parse(reader)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cmd.Name != GetCommand {
		t.Errorf("Expected Name to be %s, got %s", GetCommand, cmd.Name)
	}
	if cmd.Key != "key" {
		t.Error("Expected key to be key")
	}
	if cmd.Flags != 0 {
		t.Error("Expected flags to be 0")
	}
	if cmd.Expiry != 0 {
		t.Error("Expected expiry to be 0")
	}
	if cmd.Bytes != 0 {
		t.Error("Expected bytes to be 0")
	}
	if string(cmd.Value) != "" {
		t.Error("Expected value to be empty")
	}
}

func TestDeleteCommand(t *testing.T) {
	cp := NewParser()
	reader := bytes.NewReader([]byte("delete key\r\n"))
	cmd, err := cp.Parse(reader)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cmd.Name != DeleteCommand {
		t.Errorf("Expected Name to be %s, got %s", DeleteCommand, cmd.Name)
	}
	if cmd.Key != "key" {
		t.Error("Expected key to be key")
	}
	if cmd.Flags != 0 {
		t.Error("Expected flags to be 0")
	}
	if cmd.Expiry != 0 {
		t.Error("Expected expiry to be 0")
	}
	if cmd.Bytes != 0 {
		t.Error("Expected bytes to be 0")
	}
	if string(cmd.Value) != "" {
		t.Error("Expected value to be empty")
	}
}
