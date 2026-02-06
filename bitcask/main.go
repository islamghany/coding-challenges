package main

import (
	"bitcask/bitcask"
	"flag"
	"fmt"
	"os"
)

func main() {

	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	var dbPath string

	flag.StringVar(&dbPath, "db", "bitcask.db", "Path to the database file")
	flag.Parse()

	if dbPath == "" {
		fmt.Println("No database path provided")
		return fmt.Errorf("No database path provided")
	}

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage: bitcask -db <path> <command> <key> [value]")
		return fmt.Errorf("Invalid number of arguments")
	}

	cmd := args[0]

	db, err := bitcask.NewBitcask(dbPath)
	if err != nil {
		fmt.Println("Error creating database:", err)
		return err
	}

	switch cmd {
	case "set":
		if len(args) < 3 {
			fmt.Println("Usage: bitcask -db <path> set <key> <value>")
			return fmt.Errorf("Invalid number of arguments")
		}
		key := args[1]
		value := args[2]
		err := db.Set(key, value)
		if err != nil {
			fmt.Println("Error setting key:", err)
			return err
		}
	case "get":

		if len(args) < 2 {
			fmt.Println("Usage: bitcask -db <path> get <key>")
			return fmt.Errorf("Invalid number of arguments")
		}
		key := args[1]
		value, err := db.Get(key)
		if err != nil {
			fmt.Println("Error getting key:", err)
			return err
		}
		fmt.Println(value)
	case "del":
		if len(args) < 2 {
			fmt.Println("Usage: bitcask -db <path> del <key>")
			return fmt.Errorf("Invalid number of arguments")
		}
		key := args[1]
		err := db.Delete(key)
		if err != nil {
			fmt.Println("Error deleting key:", err)
			return err
		}
		fmt.Println("Key deleted successfully")
	case "merge":
		err := db.Merge()
		if err != nil {
			fmt.Println("Error merging database:", err)
			return err
		}
		fmt.Println("Database merged successfully")
	default:
		fmt.Println("Invalid command")
		return fmt.Errorf("Invalid command")
	}
	return nil
}
