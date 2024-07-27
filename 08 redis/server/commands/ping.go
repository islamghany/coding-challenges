package commands

import "redis/foundation/enconder/resp"

func Ping(repsArray []resp.RESPData) resp.RESPData {
	respData := resp.NewSimpleString("PONG")
	if len(repsArray) > 1 {
		arg, ok := repsArray[1].Data.(string)
		if !ok {
			return resp.NewError(InvalidArguments)
		}
		respData = resp.NewBulkString(arg)
	}

	return respData
}
