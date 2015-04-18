package main

import (
	"math/rand"
	"strings"
)

const alive = "o"
const dead = " "

type Generation struct {
	cells  [][]string
	height int
	width  int
}

func (g *Generation) Awaken() {
	cellTypes := []string{dead, alive}

	g.cells = make([][]string, g.height)

	for row := 0; row < g.height; row++ {
		cols := make([]string, g.width)

		for col := 0; col < g.width; col++ {
			cols[col] = cellTypes[random(0, 2)]
		}
		g.cells[row] = cols
	}
}

func (g *Generation) Reproduce() {
	nextGeneration := make([][]string, g.height)

	for rIndex, row := range g.cells {
		nextGeneration[rIndex] = make([]string, g.width)

		for cIndex, cell := range row {
			neighbors := g.findNeighbors(rIndex, cIndex, g.cells)

			if willSurvive(cell, neighbors) {
				nextGeneration[rIndex][cIndex] = alive
			} else {
				nextGeneration[rIndex][cIndex] = dead
			}

		}
	}

	g.cells = nextGeneration
}

func (g *Generation) findNeighbors(rowIndex int, colIndex int, cells [][]string) []string {
	var neighbors []string

	// add top
	if rowIndex != 0 {
		if cells[rowIndex-1][colIndex] == alive {
			neighbors = append(neighbors, cells[rowIndex-1][colIndex])
		}
	}

	// add bottom
	if rowIndex != g.height-1 {
		if cells[rowIndex+1][colIndex] == alive {
			neighbors = append(neighbors, cells[rowIndex+1][colIndex])
		}
	}

	// add left
	if colIndex != 0 {
		if cells[rowIndex][colIndex-1] == alive {
			neighbors = append(neighbors, cells[rowIndex][colIndex-1])
		}
	}

	// add right
	if colIndex != g.width-1 {
		if cells[rowIndex][colIndex+1] == alive {
			neighbors = append(neighbors, cells[rowIndex][colIndex+1])
		}
	}

	// add upper-left-diagonal
	if rowIndex != 0 {
		if colIndex != 0 {
			if cells[rowIndex-1][colIndex-1] == alive {
				neighbors = append(neighbors, cells[rowIndex-1][colIndex-1])
			}
		}
	}

	// add upper-right-diagonal
	if rowIndex != 0 {
		if colIndex != g.width-1 {
			if cells[rowIndex-1][colIndex+1] == alive {
				neighbors = append(neighbors, cells[rowIndex-1][colIndex+1])
			}
		}
	}

	// add lower-left-diagonal
	if rowIndex != g.height-1 {
		if colIndex != 0 {
			if cells[rowIndex+1][colIndex-1] == alive {
				neighbors = append(neighbors, cells[rowIndex+1][colIndex-1])
			}
		}
	}

	// add lower-right-diagonal
	if rowIndex != g.height-1 {
		if colIndex != g.width-1 {
			if cells[rowIndex+1][colIndex+1] == alive {
				neighbors = append(neighbors, cells[rowIndex+1][colIndex+1])
			}
		}
	}

	return neighbors
}

func willSurvive(cell string, neighbors []string) bool {
	livingNeighbors := len(neighbors)

	if cell == alive {

		if livingNeighbors < 2 {
			return false
		}

		if livingNeighbors == 2 || livingNeighbors == 3 {
			return true
		}

		if livingNeighbors > 3 {
			return false
		}

	}

	if cell == dead && livingNeighbors == 3 {
		return true
	}

	return false
}

func (g *Generation) ToString() string {
	rowStrings := make([]string, g.height)

	for index, row := range g.cells {
		rowStrings[index] = strings.Join(row, " ")
	}

	return strings.Join(rowStrings, "\n")
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
