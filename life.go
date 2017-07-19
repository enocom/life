// Package life implements Conway's Game of Life with a CLI interface.
//
// The rules of the game are as follows:
//
// 1) Any live cell with fewer than two live neighbours dies,
//    as if caused by underpopulation.
//
// 2) Any live cell with two or three live neighbours lives
//    on to the next generation.
//
// 3) Any live cell with more than three live neighbours dies,
//    as if by overpopulation.
//
// 4) Any dead cell with exactly three live neighbours becomes
//    a live cell, as if by reproduction.
//
// Typical usage is as follows:
//     g := life.NewGame()
//     g.Start()
//
package life

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

// NewLiveCell creates a living cell
func NewLiveCell() Cell {
	return Cell{alive: true}
}

// NewDeadCell creates a dead cell
func NewDeadCell() Cell {
	return Cell{alive: false}
}

// Cell represents a single living entity
type Cell struct {
	alive bool
}

// Alive returns the state of the cell
func (c Cell) Alive() bool {
	return c.alive
}

// String returns a string representation of a Cell
func (c Cell) String() string {
	if c.Alive() {
		return "o"
	}

	return " "
}

// Dimension models the size of a game of life
type Dimension struct {
	X int
	Y int
}

// LeftEdge returns whether an index is on the left edge of the board
func (d Dimension) LeftEdge(idx int) bool {
	return idx%d.X == 0
}

// RightEdge returns whether an index is on the right edge of the board
func (d Dimension) RightEdge(idx int) bool {
	if idx == 0 {
		return false
	}

	return idx%d.X == d.X-1
}

// LastRowFirstIndex returns the first index of the last row
func (d Dimension) LastRowFirstIndex() int {
	return (d.Y * d.X) - d.X
}

// CellGenerator defines the interface used to generate cells
type CellGenerator interface {
	Generate() Cell
}

// NewFixedCellGenerator is used when users want to configure a deterministic
// collection of cells in a generation
func NewFixedCellGenerator(c []Cell) CellGenerator {
	return &fixedCellGenerator{
		nextIdx: 0,
		cells:   c,
	}
}

type fixedCellGenerator struct {
	nextIdx int
	cells   []Cell
}

func (g *fixedCellGenerator) Generate() Cell {
	if g.nextIdx > len(g.cells)-1 {
		return NewDeadCell()
	}
	cell := g.cells[g.nextIdx]
	g.nextIdx++

	return cell
}

// NewRandomCellGenerator creates a CellGenerator which will return living and
// dead cells randomly
func NewRandomCellGenerator() CellGenerator {
	return &randomCellGenerator{
		r: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

type randomCellGenerator struct {
	r *rand.Rand
}

func (g *randomCellGenerator) Generate() Cell {
	alive := g.r.Intn(2)
	if alive == 0 {
		return NewDeadCell()
	}

	return NewLiveCell()
}

// Option is the underlying type for various configurations of a Generation
type Option func(*Generation)

// WithDimension configures the dimensions within which cells will live and die
func WithDimension(d Dimension) Option {
	return func(g *Generation) {
		g.dimensions = d
	}
}

// WithCells configures a generation to be seeded with the cells passed into
// the function. Use this option when configuring a generation to start at with
// a fixed collection of cells.
func WithCells(c []Cell) Option {
	return func(g *Generation) {
		g.generator = NewFixedCellGenerator(c)
	}
}

// WithRandomCells configures a generation to be seeded with living and dead
// cells randomly.
func WithRandomCells() Option {
	return func(g *Generation) {
		g.generator = NewRandomCellGenerator()
	}
}

// NewGeneration returns a single generation of cells
func NewGeneration(opts ...Option) *Generation {
	g := &Generation{
		dimensions: Dimension{X: 3, Y: 3},
		generator:  NewRandomCellGenerator(),
	}

	for _, o := range opts {
		o(g)
	}

	var cells []Cell
	for i := 0; i < g.dimensions.X*g.dimensions.Y; i++ {
		cells = append(cells, g.generator.Generate())
	}
	g.cells = cells

	return g
}

// Generation represents a collective state of living
// and dead cells
type Generation struct {
	dimensions Dimension
	generator  CellGenerator
	cells      []Cell
}

// Cells returns the generation's cells
func (g *Generation) Cells() []Cell {
	return g.cells
}

// String returns a representation of Generation
func (g *Generation) String() string {
	display := ""
	for row := 0; row < g.dimensions.Y; row++ {
		for column := 0; column < g.dimensions.X; column++ {
			display += fmt.Sprintf("%v", g.cells[column+row*g.dimensions.X])

			if column%g.dimensions.X == g.dimensions.X-1 {
				display += "\n"
			} else {
				display += " "
			}
		}
	}
	return display
}

// Next produces the next generation with some cells living
// and some cells dying
func Next(g1 *Generation) *Generation {
	g1Cells := g1.cells
	var g2Cells []Cell
	for i, cell := range g1Cells {
		nextCell := generate(i, cell, g1Cells, g1.dimensions)
		g2Cells = append(g2Cells, nextCell)
	}
	return &Generation{
		dimensions: g1.dimensions,
		cells:      g2Cells,
	}
}

func generate(idx int, c Cell, cells []Cell, d Dimension) Cell {
	liveNeighbors := leftCell(idx, cells, d.X) +
		rightCell(idx, cells, d.X) +
		aboveCell(idx, cells, d) +
		belowCell(idx, cells, d) +
		aboveDiagonalCells(idx, cells, d) +
		belowDiagonalCells(idx, cells, d)

	if !c.Alive() && liveNeighbors == 3 {
		return NewLiveCell()
	}

	if !c.Alive() {
		return NewDeadCell()
	}

	switch liveNeighbors {
	case 0, 1:
		return NewDeadCell()
	case 2, 3:
		return NewLiveCell()
	default:
		return NewDeadCell()
	}
}

// checkLeft determines if the left cell is alive
func leftCell(idx int, cells []Cell, x int) int {
	if idx%x == 0 {
		return 0
	}

	if cells[idx-1].Alive() {
		return 1
	}

	return 0
}

// checkRight determines if the right cellis alive
func rightCell(idx int, cells []Cell, x int) int {
	if idx%x == x-1 {
		return 0
	}

	if cells[idx+1].Alive() {
		return 1
	}

	return 0
}

func aboveCell(idx int, cells []Cell, d Dimension) int {
	// we're in the first row; there is no above
	if idx < d.X {
		return 0
	}

	if cells[idx-d.X].Alive() {
		return 1
	}

	return 0
}

func belowCell(idx int, cells []Cell, d Dimension) int {
	// we're in the last row; there is no below
	if idx >= d.LastRowFirstIndex() {
		return 0
	}

	if cells[idx+d.X].Alive() {
		return 1
	}

	return 0
}

func aboveDiagonalCells(idx int, cells []Cell, d Dimension) int {
	count := 0

	// we're in the first row; there is no above
	if idx < d.X {
		return count
	}

	// diagonal left
	if !d.LeftEdge(idx) && cells[idx-d.X-1].Alive() {
		count++
	}

	// diagonal right
	if !d.RightEdge(idx) && cells[idx-d.X+1].Alive() {
		count++
	}

	return count
}

func belowDiagonalCells(idx int, cells []Cell, d Dimension) int {
	count := 0
	// we're in the last row; there is no below
	lastRowStartIdx := (d.Y * d.X) - d.X
	if idx >= lastRowStartIdx {
		return 0
	}

	// diagonal left
	if !d.LeftEdge(idx) && cells[idx+d.X-1].Alive() {
		count++
	}

	// diagonal right
	if !d.RightEdge(idx) && cells[idx+d.X+1].Alive() {
		count++
	}

	return count
}

// NewTerminalUI creates a UI whose output is printing to a terminal
func NewTerminalUI(w io.Writer) *TermUI {
	return &TermUI{
		w: w,
	}
}

// TermUI represents a UI runs within a Bash shell
type TermUI struct {
	w io.Writer
}

// ClearScreen provides a means to simulate animation between generations
func (t *TermUI) ClearScreen() {
	_, _ = t.w.Write([]byte("\033[H\033[2J"))
}

// Write prints the frame to the screen
func (t *TermUI) Write(frame string) {
	_, _ = t.w.Write([]byte(frame))
}

// UI represents the interface all implementors must honor
type UI interface {
	ClearScreen()
	Write(string)
}

// GameOption provides a means to configure optional parameters
type GameOption func(*Game)

// WithBoardSize configures the dimensions of the game
func WithBoardSize(size int) GameOption {
	return func(g *Game) {
		g.dimension = Dimension{X: size, Y: size}
	}
}

// WithGenerationRate configures the speed by which one generation gives way to
// another
func WithGenerationRate(rate time.Duration) GameOption {
	return func(g *Game) {
		g.rate = rate
	}
}

// WithUI configures the UI used by the Game
func WithUI(ui UI) GameOption {
	return func(g *Game) {
		g.ui = ui
	}
}

// NewGame creates an unstarted game
func NewGame(opts ...GameOption) *Game {
	g := &Game{
		ui:        NewTerminalUI(os.Stdout),
		dimension: Dimension{X: 10, Y: 10},
		rate:      time.Second,
	}

	for _, o := range opts {
		o(g)
	}

	return g
}

// Game represents a single run of Conway's Game of Life
type Game struct {
	ui        UI
	dimension Dimension
	rate      time.Duration
}

// Start begins the game
func (g *Game) Start() {
	currentGen := NewGeneration(WithDimension(g.dimension))
	g.ui.ClearScreen()
	g.ui.Write(currentGen.String())

	for range time.Tick(g.rate) {
		currentGen = Next(currentGen)
		g.ui.ClearScreen()
		g.ui.Write(currentGen.String())
	}
}
