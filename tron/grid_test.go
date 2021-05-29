package tron

import (
	"fmt"
	"testing"
)

func TestGridInit(t *testing.T) {
	g := newGrid(3, 3)
	g.initialize()
	tests := []gridTestCase{
		{x: 0, y: 0, expected: WALL},
		{x: 1, y: 1, expected: EMPTY},
		{x: 2, y: 2, expected: WALL},
	}
	gridTestCases(tests).run(t, g)
}

type gridTestCase struct {
	x        int
	y        int
	expected uint8
}

func (g gridTestCase) run(t *testing.T, grid *grid) {
	got := grid.get(g.x, g.y)
	if got != uint8(g.expected) {
		t.Errorf("wanted %d at (%d, %d), got %d", g.expected, g.x, g.y, got)
		fmt.Println(grid.String())
	}
}

type gridTestCases []gridTestCase

func (g gridTestCases) run(t *testing.T, grid *grid) {
	for _, c := range g {
		c.run(t, grid)
	}
}
