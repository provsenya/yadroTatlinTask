package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func CountNames(r io.Reader) (map[string]int, error) {
	counts := make(map[string]int)

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())

		if name == "" {
			continue
		}

		counts[name]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return counts, nil
}

func PrintNameCounts(counts map[string]int) {
	for name, count := range counts {
		fmt.Printf("%s:%d\n", name, count)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file_path>")
		return
	}

	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	counts, err := CountNames(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	PrintNameCounts(counts)
}
