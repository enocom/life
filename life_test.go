package life_test

import (
	"testing"

	"github.com/enocom/life"
)

type test struct {
	dimen  life.Dimension
	before []life.Cell
	after  []life.Cell
}

var testCases = map[string]test{
	"Rule 1": test{
		dimen:  life.Dimension{X: 1, Y: 1},
		before: []life.Cell{life.NewLiveCell()},
		after:  []life.Cell{life.NewDeadCell()},
	},
	"Rule 2": test{
		dimen: life.Dimension{X: 3, Y: 1},
		before: []life.Cell{
			life.NewLiveCell(),
			life.NewLiveCell(),
			life.NewLiveCell(),
		},
		after: []life.Cell{
			life.NewDeadCell(),
			life.NewLiveCell(),
			life.NewDeadCell(),
		},
	},
	"Rule 2 and Rule 3 and Rule 4": test{
		dimen: life.Dimension{X: 3, Y: 3},
		// - - -
		// - O O
		// O O O
		before: []life.Cell{
			life.NewDeadCell(),
			life.NewDeadCell(),
			life.NewDeadCell(),

			life.NewDeadCell(),
			life.NewLiveCell(),
			life.NewLiveCell(),

			life.NewLiveCell(),
			life.NewLiveCell(),
			life.NewLiveCell(),
		},
		// - - -
		// O - O
		// O - O
		after: []life.Cell{
			life.NewDeadCell(),
			life.NewDeadCell(),
			life.NewDeadCell(),

			life.NewLiveCell(),
			life.NewDeadCell(),
			life.NewLiveCell(),

			life.NewLiveCell(),
			life.NewDeadCell(),
			life.NewLiveCell(),
		},
	},
}

func TestNext(t *testing.T) {
	for description, tc := range testCases {
		g1 := life.NewGeneration(
			life.WithDimension(tc.dimen),
			life.WithCells(tc.before),
		)

		g2 := life.Next(g1)

		got := g2.Cells()
		want := tc.after
		if !equal(got, want) {
			t.Errorf("(%s): want %v, got %v", description, want, got)
		}
	}
}

func equal(a, b []life.Cell) bool {
	if len(a) != len(b) {
		return false
	}

	for i, cell := range a {
		if cell != b[i] {
			return false
		}
	}

	return true
}

func TestGenerationString(t *testing.T) {
	g := life.NewGeneration(
		life.WithDimension(life.Dimension{X: 2, Y: 2}),
		life.WithCells([]life.Cell{
			life.NewLiveCell(),
			life.NewDeadCell(),

			life.NewDeadCell(),
			life.NewLiveCell(),
		}),
	)

	display := g.String()
	expected := "o  \n  o\n"
	if display != expected {
		t.Errorf("want: %#v, got: %#v", expected, display)
	}
}

func TestLeftEdge(t *testing.T) {
	d := life.Dimension{X: 3, Y: 3}

	leftEdges := []int{0, 3, 6}
	for _, e := range leftEdges {
		result := d.LeftEdge(e)

		if result != true {
			t.Errorf("want: true, got: %v (idx = %v)", result, e)
		}
	}

	nonEdges := []int{1, 2, 4, 5, 7, 8}
	for _, n := range nonEdges {
		result := d.LeftEdge(n)

		if result != false {
			t.Errorf("want: false, got: %v (idx = %v)", result, n)
		}
	}
}

func TestRightEdge(t *testing.T) {
	d := life.Dimension{X: 3, Y: 3}

	rightEdges := []int{2, 5, 8}
	for _, e := range rightEdges {
		result := d.RightEdge(e)

		if result != true {
			t.Errorf("want: true, got: %v (idx = %v)", result, e)
		}
	}

	nonEdges := []int{0, 1, 3, 4, 6, 7}
	for _, n := range nonEdges {
		result := d.RightEdge(n)

		if result != false {
			t.Errorf("want: false, got: %v (idx = %v)", result, n)
		}
	}
}
