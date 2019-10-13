package main

import "testing"

func TestNormalize(t *testing.T) {
	testCases := []struct{
		input string
		want string
	}{
		{"1234567890", "1234567890" },
		{"123 456 7891", "1234567891" },
		{"(123) 456 7892", "1234567892" },
		{"123-456-7894", "1234567894" },
		{"(123)456-7892", "1234567892" },
	}

	for _, test := range testCases {
		got := normalize(test.input)
		if got != test.want {
			t.Errorf("Got %s want %s\n", got, test.want)
		}
	}
}
