package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var gol Life

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

type pos [2]int

type Cell struct {
	row   int
	col   int
	alive bool
	char  string
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
	percentStartAlive   int
	colorsActive        bool
}

type Life struct {
	Board  *board
	params *GameParams
}

func (c *Cell) birth() {
	c.alive = true
	if gol.params.colorsActive == true {
		c.char = strings.Join([]string{"\033[35m", "o", "\033[0m"}, "")
	} else {
		c.char = "o"
	}
}

func (c *Cell) live() {
	c.alive = true
	if gol.params.colorsActive == true {
		c.char = strings.Join([]string{"\033[34m", "o", "\033[0m"}, "")
	} else {
		c.char = "o"
	}
}

func (c *Cell) kill() {
	c.alive = false
	if gol.params.colorsActive == true {
		c.char = strings.Join([]string{"\033[31m", " ", "\033[0m"}, "")
	} else {
		c.char = " "
	}
}

func (c *Cell) CheckRules() string {
	adjacentAliveCount := 0

	for _, p := range adjacent {
		if outOfBounds := checkArrayBounds(c, &p, gol.Board); outOfBounds == true {
			continue
		} else if gol.Board.currentGen[c.row+p[0]][c.col+p[1]].alive == true {
			adjacentAliveCount++
		}
	}

	if c.alive == true && adjacentAliveCount < 2 {
		return "dies" //underpopulation
	} else if c.alive == true && (adjacentAliveCount == 2 || adjacentAliveCount == 3) {
		return "lives" //life
	} else if c.alive == true && adjacentAliveCount > 3 {
		return "dies" //overpopulation
	} else if c.alive == false && adjacentAliveCount == 3 {
		return "reborn" //reproduction
	} else {
		return "dies"
	}
}

func (l *Life) Sleep() {
	time.Sleep(time.Duration(gol.params.generationalDelayMs) * time.Millisecond)
}

func (l *Life) Tick() {
	for rowIdx, cellArray := range l.Board.currentGen {
		for colIdx, cell := range cellArray {
			if result := cell.CheckRules(); result == "lives" {
				l.Board.nextGen[rowIdx][colIdx].live()
			} else if result == "dies" {
				l.Board.nextGen[rowIdx][colIdx].kill()
			} else if result == "reborn" {
				l.Board.nextGen[rowIdx][colIdx].birth()
			}
		}
	}

	for rowIdx := range l.Board.nextGen {
		copy(l.Board.currentGen[rowIdx], l.Board.nextGen[rowIdx])
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

			if l.params.percentStartAlive > 0 {
				randInt := rand.Intn(100)
				if randInt <= l.params.percentStartAlive {
					cell.live()
				} else {
					cell.kill()
				}
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

func main() {
	p := &GameParams{
		height:              50,
		width:               180,
		generationalDelayMs: 300,
		percentStartAlive:   0,
		colorsActive:        true,
	}

	InitGame(p)

	sc := ShapeCreator{}
	sc.CreateGosperGliderGun(gol.Board, [2]int{10, 20})
	sc.CreateGosperGliderGun(gol.Board, [2]int{30, 50})

	gol.Print()

	for {
		gol.Tick()
		gol.Print()
		gol.Sleep()
	}
}

type ShapeCreator struct{}

func (s *ShapeCreator) CreateGlider(b *board, p pos) {
	points := []pos{
		[2]int{p[0], p[1]},
		[2]int{p[0] + 1, p[1] + 1},
		[2]int{p[0] + 1, p[1] + 2},
		[2]int{p[0], p[1] + 2},
		[2]int{p[0] - 1, p[1] + 2},
	}

	for _, position := range points {
		b.currentGen[position[0]][position[1]].live()
	}
}

func (s *ShapeCreator) CreateGosperGliderGun(b *board, p pos) {
	points := []pos{
		[2]int{p[0], p[1]},
		[2]int{p[0] + 1, p[1]},
		[2]int{p[0], p[1] + 1},
		[2]int{p[0] + 1, p[1] + 1},

		[2]int{p[0] - 2, p[1] + 13},
		[2]int{p[0] - 2, p[1] + 12},
		[2]int{p[0] - 1, p[1] + 11},
		[2]int{p[0], p[1] + 10},
		[2]int{p[0] + 1, p[1] + 10},
		[2]int{p[0] + 1, p[1] + 14},
		[2]int{p[0] + 2, p[1] + 10},
		[2]int{p[0] + 3, p[1] + 11},
		[2]int{p[0] + 4, p[1] + 12},
		[2]int{p[0] + 4, p[1] + 13},

		[2]int{p[0] + 4, p[1] + 13},
		[2]int{p[0] + 1, p[1] + 14},
		[2]int{p[0] - 1, p[1] + 15},
		[2]int{p[0] + 3, p[1] + 15},

		[2]int{p[0] + 1, p[1] + 16},
		[2]int{p[0] + 2, p[1] + 16},
		[2]int{p[0], p[1] + 16},
		[2]int{p[0] + 1, p[1] + 17},

		[2]int{p[0], p[1] + 20},
		[2]int{p[0], p[1] + 21},
		[2]int{p[0] - 1, p[1] + 20},
		[2]int{p[0] - 1, p[1] + 21},
		[2]int{p[0] - 2, p[1] + 20},
		[2]int{p[0] - 2, p[1] + 21},

		[2]int{p[0] + 1, p[1] + 22},
		[2]int{p[0] - 3, p[1] + 22},

		[2]int{p[0] + 1, p[1] + 24},
		[2]int{p[0] + 2, p[1] + 24},
		[2]int{p[0] - 3, p[1] + 24},
		[2]int{p[0] - 4, p[1] + 24},

		[2]int{p[0] - 2, p[1] + 34},
		[2]int{p[0] - 1, p[1] + 34},
		[2]int{p[0] - 2, p[1] + 35},
		[2]int{p[0] - 1, p[1] + 35},
	}

	for _, position := range points {
		b.currentGen[position[0]][position[1]].live()
	}
}
