package resp

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestSerialize(t *testing.T) {}

func TestDeserialize(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		expected interface{}
		t        RESPType
		err      error
	}{
		// Simple String
		{
			name:     "Simple String 1",
			data:     []byte("+OK\r\n"),
			t:        SimpleString,
			expected: "OK",
		},
		{
			name:     "Simple String 2",
			data:     []byte("+Hello, World!\r\n"),
			expected: "Hello, World!",
			t:        SimpleString,
		},
		// Simple Error
		{
			name:     "Simple Error 1",
			data:     []byte("-ERR\r\n"),
			t:        Error,
			expected: "ERR",
		},
		{
			name:     "Simple Error 2",
			data:     []byte("-Error message\r\n"),
			t:        Error,
			expected: "Error message",
		},
		{
			name:     "Simple Error Error",
			data:     []byte("-Error message"),
			t:        Error,
			expected: "",
			err:      ErrMalformedType,
		},
		// Integer
		{
			name:     "Integer 1",
			data:     []byte(":123\r\n"),
			t:        Integer,
			expected: 123,
		},
		{
			name:     "Integer 2",
			data:     []byte(":0\r\n"),
			t:        Integer,
			expected: 0,
		},
		{
			name:     "Integer Error",
			data:     []byte(":asas\r\n"),
			t:        Integer,
			expected: 0,
			err:      ErrInvalidInteger,
		},
		// Bulk String
		{
			name:     "Bulk String 1",
			data:     []byte("$6\r\nfoobar\r\n"),
			t:        BulkString,
			expected: "foobar",
		},
		{
			name:     "Bulk String 2",
			data:     []byte("$0\r\n\r\n"),
			t:        BulkString,
			expected: "",
		},
		{
			name:     "Bulk String null",
			data:     []byte("$-1\r\n"),
			t:        BulkString,
			expected: nil,
		},
		{
			name:     "Bulk String Error",
			data:     []byte("$6\r\nfoobar"),
			t:        BulkString,
			expected: "",
			err:      ErrMalformedType,
		},
		// Array
		{
			name:     "Array 1",
			data:     []byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"),
			t:        Array,
			expected: []RESPData{{Data: "foo", Type: BulkString}, {Data: "bar", Type: BulkString}},
		},
		{
			name:     "Array 2",
			data:     []byte("*3\r\n:1\r\n:2\r\n:3\r\n"),
			t:        Array,
			expected: []RESPData{{Data: 1, Type: Integer}, {Data: 2, Type: Integer}, {Data: 3, Type: Integer}},
		},
		{
			name:     "Array Error",
			data:     []byte("*3\r\n:1\r\n:2\r\n:3"),
			t:        Array,
			expected: nil,
			err:      ErrMalformedType,
		},
		{
			name:     "Array nil",
			data:     []byte("*-1\r\n"),
			t:        Array,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewReader(tc.data))
			result, err := Deserialize(reader)
			if err != nil && err != tc.err {
				t.Errorf("unexpected error: %v", err)
			}
			if result != nil && result.Type != tc.t {
				t.Errorf("expected %v, got %v", tc.t, result.Type)
			}
			if result != nil && result.Type == Array {
				if result.Data == nil && tc.expected == nil {
					return
				}
				if len(result.Data.([]RESPData)) != len(tc.expected.([]RESPData)) {
					t.Errorf("expected %v, got %v", tc.expected, result.Data)
				}
				for _, v := range result.Data.([]RESPData) {
					fmt.Println(v.Data)
				}
			} else if result != nil && result.Data != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result.Data)
			}
		})
	}
}
