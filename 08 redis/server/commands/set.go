package commands

import "redis/foundation/enconder/resp"

func (cmdr *Commander) Set(repsArray []resp.RESPData) resp.RESPData {
	if len(repsArray) < 3 {
		return resp.NewError(InvalidArguments)
	}
	key, ok1 := repsArray[1].Data.(string)
	value, ok2 := repsArray[2].Data.(string)
	if !ok1 || !ok2 {
		return resp.NewError(InvalidArguments)
	}
	cmdr.Store.Set.Add(key, value)
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
