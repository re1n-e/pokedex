package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " Hello World ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello			WorlD			",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("The actual slice len dosen't match expected slice len for %s", c.input)
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Expected %s got %s", expectedWord, word)
			}
		}
	}
}
