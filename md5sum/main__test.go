package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun(t *testing.T) {
	dir := t.TempDir()
	a := writeFile(t, dir, "a.txt", "File one")
	b := writeFile(t, dir, "b.txt", "File two")
	missing := filepath.Join(dir, "missing.txt")

	const (
		hEmpty = "d41d8cd98f00b204e9800998ecf8427e"
		hHello = "8b1a9953c4611296a827abf8c47804d7"
		hA     = "5d1c020d53f38ccf82ec532f35c9ca27"
		hB     = "51ae396c5595862e990e318c2176addb"
	)

	tests := []struct {
		name     string
		args     []string
		stdin    string
		wantOut  string
		wantErr  string // substring expected on stderr; "" means stderr must be empty
		wantCode int
	}{
		{"empty stdin", nil, "", hEmpty + "  -\n", "", 0},
		{"stdin no args", nil, "Hello", hHello + "  -\n", "", 0},
		{"explicit dash", []string{"-"}, "Hello", hHello + "  -\n", "", 0},
		{"binary mode", []string{"-b"}, "Hello", hHello + " *-\n", "", 0},
		{"two files", []string{a, b}, "",
			hA + "  " + a + "\n" + hB + "  " + b + "\n", "", 0},

		{"missing file continues", []string{a, missing, b}, "",
			hA + "  " + a + "\n" + hB + "  " + b + "\n",
			"missing.txt: no such file or directory\n", 1},
		{"dash mixed with files", []string{a, "-", b}, "Hello",
			hA + "  " + a + "\n" + hHello + "  -\n" + hB + "  " + b + "\n", "", 0},
		{"unknown flag", []string{"-x"}, "", "", "-x", 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out, errOut bytes.Buffer
			code := run(tc.args, strings.NewReader(tc.stdin), &out, &errOut)

			if got := out.String(); got != tc.wantOut {
				t.Errorf("stdout:\n got %q\nwant %q", got, tc.wantOut)
			}
			if tc.wantErr == "" && errOut.Len() > 0 {
				t.Errorf("unexpected stderr: %q", errOut.String())
			}
			if tc.wantErr != "" && !strings.Contains(errOut.String(), tc.wantErr) {
				t.Errorf("stderr %q should contain %q", errOut.String(), tc.wantErr)
			}
			if code != tc.wantCode {
				t.Errorf("exit code = %d, want %d", code, tc.wantCode)
			}
		})
	}
}
