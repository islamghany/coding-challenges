package main

import (
	"os/exec"
	"testing"
)

func TestWithCutCommand(t *testing.T) {

	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "Test with -f option",
			args: []string{"-f1", "sample.tsv"},
		},
		{
			name: "Test with -f multiple fields",
			args: []string{"-f1,2,3", "sample.tsv"},
		},
		{
			name: "Test with -f  and -d option",
			args: []string{"-f1", "-d,", "fourchords.csv"},
		},
		{
			name: "Test with -f multiple fields and -d option",
			args: []string{"-f1,2,3", "-d,", "fourchords.csv"},
		},
		{
			name: "Test with -f option and -d option and pipe",
			args: []string{"-f1", "-d,", "fourchords.csv"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("cut", tc.args...)
			cutOutput, err := cmd.Output()
			if err != nil {
				t.Fatalf("Error executing cut command: %v", err)
			}
			cmdCut := exec.Command("./cccut", tc.args...)
			cccutOutput, err := cmdCut.Output()
			if err != nil {
				t.Fatalf("Error executing cccut command: %v", err)
			}
			if string(cutOutput) != string(cccutOutput) {
				t.Fatalf("Expected output: %s, got: %s", string(cutOutput), string(cccutOutput))
			}

		})
	}
}
