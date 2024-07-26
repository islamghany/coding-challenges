package resp

import (
	"bytes"
	"fmt"
)

func Serialize(respData *RESPData) ([]byte, error) {
	switch respData.Type {
	case SimpleString:
		return serializeSimpleString(respData)
	case Error:
		return serializeError(respData)
	case Integer:
		return serializeInteger(respData)
	case BulkString:
		return serializeBulkString(respData)
	case Array:
		return serializeArray(respData)
	default:
		return nil, ErrUnknownType
	}
}

func serializeSimpleString(respData *RESPData) ([]byte, error) {
	return []byte(fmt.Sprintf("+%s\r\n", respData.Data)), nil
}

func serializeError(respData *RESPData) ([]byte, error) {
	return []byte(fmt.Sprintf("-%s\r\n", respData.Data)), nil
}

func serializeInteger(respData *RESPData) ([]byte, error) {
	return []byte(fmt.Sprintf(":%d\r\n", respData.Data)), nil
}

func serializeBulkString(respData *RESPData) ([]byte, error) {
	if respData.Data == nil {
		return []byte("$-1\r\n"), nil
	}
	data := respData.Data.(string)

	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(data), data)), nil
}

func serializeArray(respData *RESPData) ([]byte, error) {
	data := respData.Data.([]RESPData)
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("*%d\r\n", len(data)))
	for _, v := range data {
		b, err := Serialize(&v)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	return buf.Bytes(), nil
}
