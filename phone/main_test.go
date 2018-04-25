package main

import "testing"

func TestNormalize(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"(123)456-7892", "1234567892"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			got := normalize(testCase.input)

			if got != testCase.want {
				t.Errorf("got %s; want %s", got, testCase.want)
				return
			}
		})
	}
}
