package main

import (
	"fmt"
	"strings"
	"time"
)

type pos [2]int

type Cell struct {
	row   int
	col   int
	alive bool
	char  string
}

func (c *Cell) kill() {
	c.alive = false
	c.char = ". "
}

func (c *Cell) revive() {
	c.alive = true
	c.char = " o"
}

type board struct {
	currentGen [][]Cell
	nextGen    [][]Cell
	height     int
	width      int
}

type GameParams struct {
	height              int
	width               int
	generationalDelayMs int
}

type Life struct {
	Board  *board
	params *GameParams
}

var adjacent = []pos{
	[2]int{-1, -1},
	[2]int{-1, 0},
	[2]int{-1, 1},
	[2]int{0, 1},
	[2]int{0, -1},
	[2]int{1, 1},
	[2]int{1, 0},
	[2]int{1, -1},
}

var gol Life

func (l *Life) Sleep() {
	time.Sleep(time.Duration(gol.params.generationalDelayMs) * time.Millisecond)
}

func (l *Life) Tick() {
	for rowIdx, cellArray := range l.Board.currentGen {
		for colIdx, cell := range cellArray {
			if alive := cell.CheckRules(); alive == true {
				l.Board.nextGen[rowIdx][colIdx].revive()
			} else {
				l.Board.nextGen[rowIdx][colIdx].kill()
			}
		}
	}

	for rowIdx := range l.Board.nextGen {
		copy(l.Board.currentGen[rowIdx], l.Board.nextGen[rowIdx])
	}
}

func (c *Cell) CheckRules() bool {
	adjacentAliveCount := 0

	for _, p := range adjacent {
		if outOfBounds := checkArrayBounds(c, &p, gol.Board); outOfBounds == true {
			continue
		} else if gol.Board.currentGen[c.row+p[0]][c.col+p[1]].alive == true {
			adjacentAliveCount++
		}
	}

	if c.alive == true && adjacentAliveCount < 2 {
		return false //underpopulation
	} else if c.alive == true && (adjacentAliveCount == 2 || adjacentAliveCount == 3) {
		return true //life
	} else if c.alive == true && adjacentAliveCount > 3 {
		return false //overpopulation
	} else if c.alive == false && adjacentAliveCount == 3 {
		return true //reproduction
	} else {
		return false
	}
}

func checkArrayBounds(c *Cell, p *pos, g *board) bool {
	y := c.row + p[0]
	x := c.col + p[1]
	return y < 0 || x < 0 || y >= g.height || x >= g.width
}

func (l *Life) Print() {
	fmt.Println("\033[2J")
	sb := strings.Builder{}

	for _, cellArray := range l.Board.currentGen {
		for _, cell := range cellArray {
			sb.WriteString(string(cell.char))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	fmt.Println(sb.String())
}

func (l *Life) createGlider(p pos) {
	l.Board.currentGen[p[0]][p[1]].revive()
	l.Board.currentGen[p[0]+1][p[1]+1].revive()
	l.Board.currentGen[p[0]+1][p[1]+2].revive()
	l.Board.currentGen[p[0]][p[1]+2].revive()
	l.Board.currentGen[p[0]-1][p[1]+2].revive()
}

func InitGame(p *GameParams) *Life {
	gol = Life{
		Board: &board{
			currentGen: make([][]Cell, p.height),
			nextGen:    make([][]Cell, p.height),
			height:     p.height,
			width:      p.width,
		},
		params: p,
	}
	return gol.InitCells()
}

func (l *Life) InitCells() *Life {
	for rowIdx := range l.Board.currentGen {
		l.Board.currentGen[rowIdx] = make([]Cell, l.Board.width)
	}

	for rowIdx := range l.Board.currentGen {
		for colIdx := range l.Board.currentGen[rowIdx] {
			cell := Cell{
				alive: false,
				row:   rowIdx,
				col:   colIdx,
				char:  " .",
			}

			l.Board.currentGen[rowIdx][colIdx] = cell
		}
	}

	for rowIdx := range l.Board.nextGen {
		l.Board.nextGen[rowIdx] = make([]Cell, l.Board.width)
		copy(l.Board.nextGen[rowIdx], l.Board.currentGen[rowIdx])
	}

	return l
}

func main() {
	p := &GameParams{
		height:              20,
		width:               30,
		generationalDelayMs: 300,
	}

	InitGame(p)

	gol.createGlider(pos{3, 3})
	gol.createGlider(pos{8, 6})
	gol.createGlider(pos{15, 15})

	gol.Print()

	for {
		gol.Tick()
		gol.Print()
		gol.Sleep()
	}
}
