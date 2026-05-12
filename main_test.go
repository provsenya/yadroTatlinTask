package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestCountNames(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]int
	}{
		{
			name: "counts repeated names",
			input: `Алёна
Миша
Алёна
Дима`,
			expected: map[string]int{
				"Алёна": 2,
				"Миша":  1,
				"Дима":  1,
			},
		},
		{
			name: "ignores empty lines",
			input: `Алёна

Миша

Алёна`,
			expected: map[string]int{
				"Алёна": 2,
				"Миша":  1,
			},
		},
		{
			name: "trims spaces",
			input: `  Алёна  
Миша
   Алёна`,
			expected: map[string]int{
				"Алёна": 2,
				"Миша":  1,
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: map[string]int{},
		},
		{
			name: "only empty lines",
			input: `


`,
			expected: map[string]int{},
		},
		{
			name:     "only spaces",
			input:    "     \n   \n\t\n",
			expected: map[string]int{},
		},
		{
			name:  "single name",
			input: `Алёна`,
			expected: map[string]int{
				"Алёна": 1,
			},
		},
		{
			name: "names with different letter case are different",
			input: `алёна
Алёна
АЛЁНА`,
			expected: map[string]int{
				"алёна": 1,
				"Алёна": 1,
				"АЛЁНА": 1,
			},
		},
		{
			name: "names with hyphen are counted",
			input: `Анна-Мария
Анна-Мария
Миша`,
			expected: map[string]int{
				"Анна-Мария": 2,
				"Миша":       1,
			},
		},
		{
			name: "names with internal spaces are counted as full name",
			input: `Анна Мария
Анна Мария
Анна`,
			expected: map[string]int{
				"Анна Мария": 2,
				"Анна":       1,
			},
		},
		{
			name:  "windows line endings",
			input: "Алёна\r\nМиша\r\nАлёна\r\n",
			expected: map[string]int{
				"Алёна": 2,
				"Миша":  1,
			},
		},
		{
			name: "trailing newline",
			input: `Алёна
Миша
`,
			expected: map[string]int{
				"Алёна": 1,
				"Миша":  1,
			},
		},
		{
			name: "many repetitions of one name",
			input: `Дима
Дима
Дима
Дима
Дима`,
			expected: map[string]int{
				"Дима": 5,
			},
		},
		{
			name: "latin and cyrillic names",
			input: `John
Алёна
John
Миша`,
			expected: map[string]int{
				"John":  2,
				"Алёна": 1,
				"Миша":  1,
			},
		},
		{
			name:  "name with tab spaces around it",
			input: "\tАлёна\t\nМиша\n\tАлёна",
			expected: map[string]int{
				"Алёна": 2,
				"Миша":  1,
			},
		},
		{
			name: "numbers are counted as text names",
			input: `123
123
Алёна`,
			expected: map[string]int{
				"123":   2,
				"Алёна": 1,
			},
		},
		{
			name: "symbols are counted as text names",
			input: `!!!
!!!
Алёна`,
			expected: map[string]int{
				"!!!":   2,
				"Алёна": 1,
			},
		},
		{
			name: "mixed empty lines spaces and names",
			input: `

 Алёна 

Миша

 Алёна
`,
			expected: map[string]int{
				"Алёна": 2,
				"Миша":  1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)

			actual, err := CountNames(reader)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestSortNameCounts(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []NameCount
	}{
		{
			name: "sorts by count descending",
			input: map[string]int{
				"Дима":  1,
				"Миша":  3,
				"Алёна": 2,
			},
			expected: []NameCount{
				{Name: "Миша", Count: 3},
				{Name: "Алёна", Count: 2},
				{Name: "Дима", Count: 1},
			},
		},
		{
			name: "sorts by name when counts are equal",
			input: map[string]int{
				"Миша":  1,
				"Алёна": 1,
				"Дима":  1,
			},
			expected: []NameCount{
				{Name: "Алёна", Count: 1},
				{Name: "Дима", Count: 1},
				{Name: "Миша", Count: 1},
			},
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: []NameCount{},
		},
		{
			name: "one element",
			input: map[string]int{
				"Алёна": 1,
			},
			expected: []NameCount{
				{Name: "Алёна", Count: 1},
			},
		},
		{
			name: "same counts with latin and cyrillic names",
			input: map[string]int{
				"John":  1,
				"Алёна": 1,
				"Bob":   1,
			},
			expected: []NameCount{
				{Name: "Bob", Count: 1},
				{Name: "John", Count: 1},
				{Name: "Алёна", Count: 1},
			},
		},
		{
			name: "large count goes first",
			input: map[string]int{
				"Алёна": 1000,
				"Миша":  2,
				"Дима":  500,
			},
			expected: []NameCount{
				{Name: "Алёна", Count: 1000},
				{Name: "Дима", Count: 500},
				{Name: "Миша", Count: 2},
			},
		},
		{
			name: "zero count is sorted after positive counts",
			input: map[string]int{
				"Алёна": 0,
				"Миша":  2,
				"Дима":  1,
			},
			expected: []NameCount{
				{Name: "Миша", Count: 2},
				{Name: "Дима", Count: 1},
				{Name: "Алёна", Count: 0},
			},
		},
		{
			name: "negative count is sorted after zero and positive counts",
			input: map[string]int{
				"Алёна": -1,
				"Миша":  2,
				"Дима":  0,
			},
			expected: []NameCount{
				{Name: "Миша", Count: 2},
				{Name: "Дима", Count: 0},
				{Name: "Алёна", Count: -1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := SortNameCounts(tt.input)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}
