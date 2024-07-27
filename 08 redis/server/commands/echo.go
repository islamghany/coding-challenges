package commands

import "redis/foundation/enconder/resp"

func Echo(repsArray []resp.RESPData) resp.RESPData {
	if len(repsArray) < 2 {
		return resp.NewError(InvalidArguments)
	}
	msg, ok := repsArray[1].Data.(string)
	if !ok {
		return resp.NewError(InvalidArguments)
	}

	return resp.NewBulkString(msg)
}
