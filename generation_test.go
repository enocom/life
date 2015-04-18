package main

import "testing"

// Rule 1: Any live cell with fewer than two live neighbours dies.
// Rule 2: Any live cell with two or three live neighbours lives.
// Rule 3: Any live cell with more than three live neighbours dies.
// Rule 4: Any dead cell with exactly three live neighbours becomes a live cell.

func TestWillSurvive(t *testing.T) {
	cases := []struct {
		cell      string
		neighbors []string
		want      bool
	}{
		{"o", []string{}, false},                                       // Rule 1
		{"o", []string{"o"}, false},                                    // Rule 1
		{"o", []string{"o", "o"}, true},                                // Rule 2
		{"o", []string{"o", "o", "o"}, true},                           // Rule 2
		{"o", []string{"o", "o", "o", "o"}, false},                     // Rule 3
		{"o", []string{"o", "o", "o", "o", "o"}, false},                // Rule 3
		{"o", []string{"o", "o", "o", "o", "o", "o"}, false},           // Rule 3
		{"o", []string{"o", "o", "o", "o", "o", "o", "o"}, false},      // Rule 3
		{"o", []string{"o", "o", "o", "o", "o", "o", "o", "o"}, false}, // Rule 3

		{" ", []string{"o"}, false},                                    // Rule 4
		{" ", []string{"o", "o"}, false},                               // Rule 4
		{" ", []string{"o", "o", "o"}, true},                           // Rule 4
		{" ", []string{"o", "o", "o", "o"}, false},                     // Rule 4
		{" ", []string{"o", "o", "o", "o", "o"}, false},                // Rule 4
		{" ", []string{"o", "o", "o", "o", "o", "o"}, false},           // Rule 4
		{" ", []string{"o", "o", "o", "o", "o", "o", "o"}, false},      // Rule 4
		{" ", []string{"o", "o", "o", "o", "o", "o", "o", "o"}, false}, // Rule 4
	}

	for _, test := range cases {
		ok := willSurvive(test.cell, test.neighbors)

		if ok != test.want {
			t.Errorf("willSurvive(%q, %q) = %v want %v", test.cell, test.neighbors, ok, test.want)
		}
	}
}
