package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Constants
const (
	ScreenWidth       = 500
	ScreenHeight      = 500
	GridSize          = 50
	CellSize          = ScreenWidth / GridSize
	InitialFishCount  = 100
	InitialSharkCount = 20
)

// EntityType Enum
type CellType int

const (
	Empty CellType = iota
	Fish
	Shark
)

// Entity struct to represent fish and sharks with their breeding and starving times
type Entity struct {
	Type          CellType
	BreedCounter  int
	StarveCounter int
	BreedTime     int
	StarveTime    int
}

type Grid [][]*Entity

type Game struct {
	grid Grid
}

// initialize the grid with fish and sharks
func initializeGrid() Grid {
	grid := make(Grid, GridSize)
	for i := range grid {
		grid[i] = make([]*Entity, GridSize)
	}

	placeEntities(grid, Fish, InitialFishCount, 5, 0)   // Fish breed every 5 steps
	placeEntities(grid, Shark, InitialSharkCount, 8, 5) // Sharks breed every 8 steps and starve after 5

	return grid
}

// randomly place entities on grid
func placeEntities(grid Grid, entityType CellType, count int, breedTime, starveTime int) {
	for i := 0; i < count; {
		x, y := rand.Intn(GridSize), rand.Intn(GridSize)
		if grid[x][y] == nil {
			entity := &Entity{
				Type:       entityType,
				BreedTime:  breedTime,
				StarveTime: starveTime,
			}
			if entityType == Shark {
				entity.StarveCounter = starveTime
			}
			grid[x][y] = entity
			i++
		}
	}
}

// Move fish
func moveFish(grid, newGrid Grid, x, y int) {
	cell := grid[x][y]
	cell.BreedCounter++

	neighbours := getNeighbours(x, y)
	emptyCells := filterEmptyCells(grid, neighbours)

	if len(emptyCells) > 0 {
		// Move to a random empty cell
		randomCell := emptyCells[rand.Intn(len(emptyCells))]
		newX, newY := randomCell[0], randomCell[1]
		newGrid[newX][newY] = cell
	} else {
		// Stay in place
		newGrid[x][y] = cell
	}

	// Breed fish
	if cell.BreedCounter >= cell.BreedTime {
		cell.BreedCounter = 0
		newGrid[x][y] = &Entity{Type: Fish, BreedTime: 5} // New fish has its own breed time
	}
}

// Move shark
func moveShark(grid, newGrid Grid, x, y int) {
	cell := grid[x][y]
	cell.BreedCounter++
	cell.StarveCounter--

	neighbours := getNeighbours(x, y)
	fishCells := filterFishCells(grid, neighbours)
	emptyCells := filterEmptyCells(grid, neighbours)

	if len(fishCells) > 0 {
		// Eat fish
		randomCell := fishCells[rand.Intn(len(fishCells))]
		newX, newY := randomCell[0], randomCell[1]
		newGrid[newX][newY] = cell
		cell.StarveCounter = cell.StarveTime
	} else if len(emptyCells) > 0 {
		// Move to an empty cell
		randomCell := emptyCells[rand.Intn(len(emptyCells))]
		newX, newY := randomCell[0], randomCell[1]
		newGrid[newX][newY] = cell
	} else {
		// Stay in place
		newGrid[x][y] = cell
	}

	// Starve shark
	if cell.StarveCounter <= 0 {
		newGrid[x][y] = nil
		return
	}

	// Breed shark
	if cell.BreedCounter >= cell.BreedTime {
		cell.BreedCounter = 0
		newGrid[x][y] = &Entity{Type: Shark, BreedTime: 8, StarveTime: 5} // New shark has its own breed and starve times
	}
}

// get neighbours
func getNeighbours(x, y int) [][2]int {
	return [][2]int{
		{x, (y - 1 + GridSize) % GridSize},
		{x, (y + 1) % GridSize},
		{(x - 1 + GridSize) % GridSize, y},
		{(x + 1) % GridSize, y},
	}
}

// filter empty cells
func filterEmptyCells(grid Grid, neighbours [][2]int) [][2]int {
	var emptyCells [][2]int
	for _, n := range neighbours {
		if grid[n[0]][n[1]] == nil {
			emptyCells = append(emptyCells, n)
		}
	}
	return emptyCells
}

// filter fish cells
func filterFishCells(grid Grid, neighbours [][2]int) [][2]int {
	var fishCells [][2]int
	for _, n := range neighbours {
		if cell := grid[n[0]][n[1]]; cell != nil && cell.Type == Fish {
			fishCells = append(fishCells, n)
		}
	}
	return fishCells
}

// update simulation
func updateSimulation(grid Grid) {
	newGrid := make(Grid, GridSize)
	for i := range newGrid {
		newGrid[i] = make([]*Entity, GridSize)
	}

	for x := range grid {
		for y, cell := range grid[x] {
			if cell == nil || (newGrid[x][y] != nil) {
				continue
			}

			if cell.Type == Fish {
				moveFish(grid, newGrid, x, y)
			} else if cell.Type == Shark {
				moveShark(grid, newGrid, x, y)
			}
		}
	}

	copyGrid(grid, newGrid)
}

// copy state
func copyGrid(dest, src Grid) {
	for i := range src {
		for j := range src[i] {
			dest[i][j] = src[i][j]
		}
	}
}

// draw grid
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	for x := 0; x < GridSize; x++ {
		for y := 0; y < GridSize; y++ {
			cell := g.grid[x][y]
			if cell == nil {
				continue
			}

			var colour color.RGBA
			if cell.Type == Fish {
				colour = color.RGBA{0, 255, 0, 255}
			} else if cell.Type == Shark {
				colour = color.RGBA{255, 0, 0, 255}
			}

			ebitenutil.DrawRect(screen, float64(y*CellSize), float64(x*CellSize), float64(CellSize), float64(CellSize), colour)
		}
	}
}

// Update the game state
func (g *Game) Update() error {
	updateSimulation(g.grid)
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		grid: initializeGrid(),
	}

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Wator Simulation")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
