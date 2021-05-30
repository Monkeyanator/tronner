package tron

import (
	"fmt"
	"strings"
)

type grid struct {
	height uint
	width  uint
	data   [][]uint8
}

const (
	EMPTY = 8 // don't ask
	WALL  = 9
)

func newGrid(w, h uint) *grid {
	d := make([][]uint8, h)
	for i := uint(0); i < h; i++ {
		d[i] = make([]uint8, w)
	}
	return &grid{
		height: h,
		width:  w,
		data:   d,
	}
}

func (g *grid) get(x, y int) uint8 {
	return g.data[y][x]
}

func (g *grid) set(x, y int, val uint8) {
	g.data[y][x] = val
}

func (g *grid) initialize() {
	for i := 0; i < int(g.height); i++ {
		for j := 0; j < int(g.width); j++ {
			g.data[i][j] = EMPTY
		}
	}
	for i := 0; i < int(g.height); i++ {
		g.data[i][0] = WALL
		g.data[i][g.width-1] = WALL
	}
	for i := 0; i < int(g.width); i++ {
		g.data[0][i] = WALL
		g.data[g.height-1][i] = WALL
	}
}

func (g *grid) serialize() string {
	sb := strings.Builder{}
	for y := 0; y < int(g.height); y++ {
		for x := 0; x < int(g.width); x++ {
			// this is dummy slow for something that happens 10 times a second and where
			// almost nothing changes lol
			sb.WriteString(fmt.Sprintf("%d", g.data[y][x]))
		}
	}
	return sb.String()
}

func (g *grid) String() string {
	sb := strings.Builder{}
	for x := 0; x < int(g.width); x++ {
		for y := 0; y < int(g.height); y++ {
			sb.WriteString(fmt.Sprintf("%d", g.get(x, y)))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
