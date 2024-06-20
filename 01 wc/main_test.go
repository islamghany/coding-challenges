package main

import (
	"log"
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	type wantCount struct {
		count  int
		option string
	}
	tc := []wantCount{
		{option: "-c"},
		{option: "-l"},
		{option: "-w"},
		{option: "-m"},
	}
	for _, tt := range tc {
		t.Run(tt.option, func(t *testing.T) {
			want := exec.Command("wc", tt.option, "./test.txt")
			wcOut, err := want.Output()
			if err != nil {
				log.Fatal(err)
			}
			got := exec.Command("./ccwc", tt.option, "./test.txt")
			gotOut, err := got.Output()
			if err != nil {
				log.Fatal(err)
			}
			if string(gotOut) != string(wcOut) {
				t.Errorf("Got: %v, Want: %v", string(gotOut), string(wcOut))
			}
		})
	}
}
