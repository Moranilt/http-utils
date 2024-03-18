package handler

import "testing"

var nameTests = []struct {
	name     string
	field    string
	expected bool
	index    int
}{
	{
		name:     "valid array name",
		field:    "files[]",
		index:    5,
		expected: true,
	},
	{
		name:     "ends not with close bracket",
		field:    "files[]sad",
		index:    -1,
		expected: false,
	},
	{
		name:     "not array element",
		field:    "files",
		index:    -1,
		expected: false,
	},
	{
		name:     "empty element",
		field:    "",
		index:    -1,
		expected: false,
	},
	{
		name:     "not valid index element",
		field:    "files[ashds]",
		index:    5,
		expected: false,
	},
	{
		name:     "empty index",
		field:    "files[]",
		expected: true,
		index:    5,
	},
	{
		name:     "with index",
		field:    "files[1]",
		index:    5,
		expected: false,
	},
}

func TestValidArrayName(t *testing.T) {
	for _, test := range nameTests {
		t.Run(test.name, func(t *testing.T) {
			i, valid := isValidNameArray(test.field)

			if valid != test.expected {
				t.Errorf("got %t, expected %t", valid, test.expected)
			}

			if i != test.index {
				t.Errorf("got %d, expected %d", i, test.index)
			}
		})
	}
}

func BenchmarkValidArrayName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidNameArray("sahgdjhgajhds[]")
	}
}

func TestExtractSubName(t *testing.T) {
	tests := []struct {
		name           string
		field          string
		expected_start string
		expected_sub   string
		valid          bool
	}{
		{
			name:           "valid array name",
			field:          "content[title]",
			expected_start: "content",
			expected_sub:   "title",
			valid:          true,
		},
		{
			name:           "not valid name",
			field:          "content[]",
			expected_start: "content[]",
			expected_sub:   "",
			valid:          false,
		},
		{
			name:           "not valid name",
			field:          "content",
			expected_start: "content",
			expected_sub:   "",
			valid:          false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			start, sub, valid := extractSubName(test.field)
			if valid != test.valid {
				t.Errorf("got %t, expected %t", valid, test.valid)
			}
			if start != test.expected_start {
				t.Errorf("got %s, expected %s", start, test.expected_start)
			}
			if sub != test.expected_sub {
				t.Errorf("got %s, expected %s", sub, test.expected_sub)
			}
		})
	}
}
