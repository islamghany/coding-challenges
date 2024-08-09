package commands

import (
	"fmt"
	"redis/foundation/enconder/resp"
	"redis/foundation/store"
	"strconv"
	"strings"
	"time"
)

func getExpiryOptions(repsArray []resp.RESPData) ([]store.Option, error) {
	opts := make([]store.Option, 0)
	for i := 3; i < len(repsArray); i++ {
		opt := strings.ToLower(repsArray[i].Data.(string))
		switch opt {
		case "ex", "px", "exat", "pxat":
			var res store.Option = nil
			if i+1 >= len(repsArray) {
				return nil, fmt.Errorf("Invalid arguments")
			}
			exp, ok := repsArray[i+1].Data.(string)
			if !ok {
				return nil, fmt.Errorf("Invalid %s arguments", opt)
			}
			if opt == "ex" {
				expNum, err := strconv.Atoi(exp)
				if err != nil || exp[0] == '-' {
					return nil, fmt.Errorf("Invalid %s arguments", opt)
				}
				res = store.WithEX(time.Duration(expNum) * time.Second)
			} else if opt == "px" {
				expNum, err := strconv.Atoi(exp)
				if err != nil || exp[0] == '-' {
					return nil, fmt.Errorf("Invalid %s arguments", opt)
				}
				res = store.WithPX(time.Duration(expNum) * time.Millisecond)
			} else if opt == "exat" {
				expTime, err := time.Parse(time.RFC3339, exp)
				if err != nil || exp[0] == '-' {
					return nil, fmt.Errorf("Invalid %s arguments", opt)
				}
				res = store.WithEXAT(expTime)
			} else if opt == "pxat" {
				expTime, err := time.Parse(time.RFC3339, exp)
				if err != nil || exp[0] == '-' {
					return nil, fmt.Errorf("Invalid %s arguments", opt)
				}
				res = store.WithPXAT(expTime)
			}
			if res != nil {
				opts = append(opts, res)
			}
		}
	}
	return opts, nil
}

func (cmdr *Commander) Set(repsArray []resp.RESPData) resp.RESPData {
	if len(repsArray) < 3 {
		return resp.NewError(InvalidArguments)
	}
	key, ok1 := repsArray[1].Data.(string)
	value, ok2 := repsArray[2].Data.(string)
	if !ok1 || !ok2 {
		return resp.NewError(InvalidArguments)
	}
	opts, err := getExpiryOptions(repsArray)
	if err != nil {
		return resp.NewError(err.Error())
	}

	cmdr.Store.Set.Add(key, value, opts...)
	return resp.NewSimpleString("OK")
}

func (cmdr *Commander) Get(repsArray []resp.RESPData) resp.RESPData {
	if len(repsArray) < 2 {
		return resp.NewError(InvalidArguments)
	}
	key, ok := repsArray[1].Data.(string)
	if !ok {
		return resp.NewError(InvalidArguments)
	}
	value, ok := cmdr.Store.Set.Get(key)
	if !ok {
		return resp.NewNil()
	}

	return resp.NewBulkString(value)
}
