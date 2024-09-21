package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getLCSStr(dp [][]int, s1, s2 string) string {
	i, j := len(s1), len(s2)
	strLen := dp[i][j] // Length of LCS
	str := make([]rune, strLen)
	index := strLen - 1 // Index for inserting characters into the result string

	for i > 0 && j > 0 {
		if s1[i-1] == s2[j-1] {
			str[index] = rune(s1[i-1])
			index--
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}
	return string(str)
}

func LCS(s1, s2 string) string {
	n := len(s1)
	m := len(s2)
	dp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return getLCSStr(dp, s1, s2)
}

func getLCSArray(dp [][]int, arr1, arr2 []string) []string {
	i, j := len(arr1), len(arr2)
	lcs := make([]string, dp[i][j])
	index := dp[i][j] - 1 // Start from the end of the LCS

	for i > 0 && j > 0 {
		if arr1[i-1] == arr2[j-1] {
			lcs[index] = arr1[i-1] // Match found, part of LCS
			index--
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}
	return lcs
}

func LCSArrays(arr1, arr2 []string) []string {
	n := len(arr1)
	m := len(arr2)
	dp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if arr1[i-1] == arr2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return getLCSArray(dp, arr1, arr2)
}

type diff struct {
	Op    string
	Line  int
	Value string
}

type diffArray []diff

func generateDiffs(arr1, arr2, lcsArray []string) (diffArray, diffArray) {
	d1 := make(diffArray, 0)
	d2 := make(diffArray, 0)
	i, j, k := 0, 0, 0

	for i < len(arr1) && j < len(arr2) {
		if arr1[i] == lcsArray[k] && arr2[j] == lcsArray[k] {
			k++
			i++
			j++
		} else if arr1[i] == lcsArray[k] {
			d2 = append(d2, diff{Op: "ADD", Line: j, Value: arr2[j]})
			j++
		} else if arr2[j] == lcsArray[k] {
			d1 = append(d1, diff{Op: "DEL", Line: i, Value: arr1[i]})
			i++
		} else {
			d1 = append(d1, diff{Op: "MOD", Line: i, Value: arr1[i]})
			d2 = append(d2, diff{Op: "MOD", Line: j, Value: arr2[j]})
			i++
			j++
		}
	}

	for i < len(arr1) {
		d1 = append(d1, diff{Op: "DEL", Line: i, Value: arr1[i]})
		i++
	}

	for j < len(arr2) {
		d2 = append(d2, diff{Op: "ADD", Line: j, Value: arr2[j]})
		j++
	}

	return d1, d2
}

func executediff(d1, d2 diffArray) string {
	str := strings.Builder{}
	i, j := 0, 0
	for i < len(d1) && j < len(d2) {
		if d1[i].Op == "DEL" {
			str.WriteString(fmt.Sprintf("DEL %d %s\n", d1[i].Line, d1[i].Value))
			i++
		} else if d2[j].Op == "ADD" {
			str.WriteString(fmt.Sprintf("ADD %d %s\n", d2[j].Line, d2[j].Value))
			j++
		} else {
			str.WriteString(fmt.Sprintf("MOD %d %s\n", d1[i].Line, d1[i].Value))
			i++
			j++
		}
	}
	for i < len(d1) {
		str.WriteString(fmt.Sprintf("DEL %d %s\n", d1[i].Line, d1[i].Value))
		i++
	}
	for j < len(d2) {
		str.WriteString(fmt.Sprintf("ADD %d %s\n", d2[j].Line, d2[j].Value))
		j++
	}
	return str.String()
}

func readFile(file string) []string {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <file1> <file2>")
		return
	}

	arr1 := readFile(args[0])
	arr2 := readFile(args[1])

	lcsArray := LCSArrays(arr1, arr2)

	d1, d2 := generateDiffs(arr1, arr2, lcsArray)
	fmt.Print(executediff(d1, d2))

}
