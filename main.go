package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type NameCount struct {
	Name  string
	Count int
}

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

func SortNameCounts(counts map[string]int) []NameCount {
	result := make([]NameCount, 0, len(counts))

	for name, count := range counts {
		result = append(result, NameCount{
			Name:  name,
			Count: count,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Name < result[j].Name
		}

		return result[i].Count > result[j].Count
	})

	return result
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

	result := SortNameCounts(counts)

	for _, item := range result {
		fmt.Printf("%s:%d\n", item.Name, item.Count)
	}
}
