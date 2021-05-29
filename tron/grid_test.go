package tron

import (
	"fmt"
	"testing"
)

func TestGridInit(t *testing.T) {
	g := newGrid(3, 3)
	g.initialize()
	tests := []gridTestCase{
		{x: 0, y: 0, expected: 1},
		{x: 1, y: 1, expected: 0},
		{x: 2, y: 2, expected: 1},
	}
	gridTestCases(tests).run(t, g)
}

type gridTestCase struct {
	x        int
	y        int
	expected uint8
}

type gridTestCases []gridTestCase

func (g gridTestCases) run(t *testing.T, grid *grid) {
	for _, c := range g {
		got := grid.get(c.x, c.y)
		if got != uint8(c.expected) {
			t.Errorf("wanted %d at (%d, %d), got %d", c.expected, c.x, c.y, got)
			fmt.Println(grid.String())
		}
	}
}
