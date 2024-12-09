// Wator simulation project by Seán Rourke, C00251168
package Wator

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/xuri/excelize/v2"
)

// Constants
const (
	ScreenWidth       = 500
	ScreenHeight      = 500
	GridSize          = 50
	CellSize          = ScreenWidth / GridSize
	InitialFishCount  = 200
	InitialSharkCount = 50
	FishBreedTime     = 5
	SharkBreedTime    = 8
	SharkStarveTime   = 5
)

type CellType int

const (
	Empty CellType = iota
	Fish
	Shark
)

type Entity struct {
	Type          CellType
	BreedCounter  int
	StarveCounter int
}

type Grid [][]*Entity

type Game struct {
	grid       Grid
	numThreads int
}

/**
 * @brief Initializes the grid with entities.
 *
 * This function creates a new grid of specified size and populates it with
 * initial counts of fish and sharks by calling the `PlaceEntities` function.
 * The grid is represented as a two-dimensional slice of entity pointers,
 * where each cell can either be empty or contain an entity.
 *
 * @return The initialized grid containing entities.
 */
func InitialiseGrid() Grid {
	grid := make(Grid, GridSize)
	for i := range grid {
		grid[i] = make([]*Entity, GridSize)
	}

	PlaceEntities(grid, Fish, InitialFishCount)
	PlaceEntities(grid, Shark, InitialSharkCount)

	return grid
}

/**
 * @brief Places a specified number of entities of a given type in the grid.
 *
 * This function randomly selects empty cells in the grid and places the specified
 * number of entities of the given type (either fish or sharks) into those cells.
 * If the entity type is a shark, it initializes the starvation counter for the shark.
 *
 * @param grid The grid where entities will be placed.
 * @param entityType The type of entity to place in the grid (e.g., Fish or Shark).
 * @param count The number of entities to place in the grid.
 */
func PlaceEntities(grid Grid, entityType CellType, count int) {
	for i := 0; i < count; {
		x, y := rand.Intn(GridSize), rand.Intn(GridSize)
		if grid[x][y] == nil {
			entity := &Entity{Type: entityType}
			if entityType == Shark {
				entity.StarveCounter = SharkStarveTime
			}
			grid[x][y] = entity
			i++
		}
	}
}

/**
 * @brief Moves a fish in the simulation based on its current state and surroundings.
 *
 * This function updates the position of a fish in the grid. The fish can either
 * move to a random empty cell or stay in its current position based on the
 * availability of empty cells in its neighbourhood. The function also handles
 * the breeding mechanics for the fish.
 *
 * @param grid The current state of the grid containing entities.
 * @param newGrid The grid where the updated state will be recorded.
 * @param x The x-coordinate of the fish's current position.
 * @param y The y-coordinate of the fish's current position.
 */
func MoveFish(grid, newGrid Grid, x, y int) {
	cell := grid[x][y]
	cell.BreedCounter++

	// Find empty neighbors
	neighbours := GetNeighbours(x, y)
	emptyCells := FilterEmptyCells(grid, neighbours)

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
	if cell.BreedCounter >= FishBreedTime {
		cell.BreedCounter = 0
		if newGrid[x][y] == nil {
			newGrid[x][y] = &Entity{Type: Fish}
		}
	}
}

/**
 * @brief Moves a shark in the simulation based on its current state and surroundings.
 *
 * This function updates the position of a shark in the grid. The shark can either
 * eat a fish, move to an empty cell, or stay in its current position based on the
 * availability of fish and empty cells in its neighbourhood. The function also
 * handles the breeding and starvation mechanics for the shark.
 *
 * @param grid The current state of the grid containing entities.
 * @param newGrid The grid where the updated state will be recorded.
 * @param x The x-coordinate of the shark's current position.
 * @param y The y-coordinate of the shark's current position.
 */
func MoveShark(grid, newGrid Grid, x, y int) {
	cell := grid[x][y]
	cell.BreedCounter++
	cell.StarveCounter--

	neighbours := GetNeighbours(x, y)
	fishCells := FilterFishCells(grid, neighbours)
	emptyCells := FilterEmptyCells(grid, neighbours)

	if len(fishCells) > 0 {
		// Eat fish
		randomCell := fishCells[rand.Intn(len(fishCells))]
		newX, newY := randomCell[0], randomCell[1]
		newGrid[newX][newY] = cell
		cell.StarveCounter = SharkStarveTime
	} else if len(emptyCells) > 0 {
		// Move to an empty cell
		randomCell := emptyCells[rand.Intn(len(emptyCells))]
		newX, newY := randomCell[0], randomCell[1]
		newGrid[newX][newY] = cell
		// Starve shark
		if cell.StarveCounter <= 0 {
			newGrid[newX][newY] = nil
		}
	} else {
		// Stay in place
		newGrid[x][y] = cell
		// Starve shark
		if cell.StarveCounter <= 0 {
			newGrid[x][y] = nil
		}
	}

	// Breed shark
	if cell.BreedCounter >= SharkBreedTime {
		cell.BreedCounter = 0
		if newGrid[x][y] == nil {
			newGrid[x][y] = &Entity{Type: Shark, StarveCounter: SharkStarveTime}
		}
	}
}

/**
 * @brief Returns the coordinates of the neighbouring cells for a given cell.
 *
 * This function calculates the coordinates of the four direct neighbours (up, down,
 * left, right) of the cell located at (x, y) in a toroidal grid. It ensures that
 * the coordinates wrap around the edges of the grid.
 *
 * @param x The x-coordinate of the cell for which neighbours are to be found.
 * @param y The y-coordinate of the cell for which neighbours are to be found.
 * @return A slice of coordinates representing the neighbours of the specified cell.
 */
func GetNeighbours(x, y int) [][2]int {
	return [][2]int{
		{x, (y - 1 + GridSize) % GridSize},
		{x, (y + 1) % GridSize},
		{(x - 1 + GridSize) % GridSize, y},
		{(x + 1) % GridSize, y},
	}
}

/**
 * @brief Filters and returns the coordinates of empty cells from the given neighbours.
 *
 * This function iterates through the provided list of neighbour coordinates
 * and checks each corresponding cell in the grid. If the cell is nil (empty),
 * the coordinate is added to the result list.
 *
 * @param grid The grid containing the cells to be checked.
 * @param neighbours A slice of coordinates representing the neighbours to filter.
 * @return A slice of coordinates of the empty cells found among the neighbours.
 */
func FilterEmptyCells(grid Grid, neighbours [][2]int) [][2]int {
	var emptyCells [][2]int
	for _, n := range neighbours {
		if grid[n[0]][n[1]] == nil {
			emptyCells = append(emptyCells, n)
		}
	}
	return emptyCells
}

/**
 * @brief Filters and returns the coordinates of fish cells from the given neighbours.
 *
 * This function iterates through the provided list of neighbour coordinates
 * and checks each corresponding cell in the grid. If the cell is not nil
 * and its type is `Fish`, the coordinate is added to the result list.
 *
 * @param grid The grid containing the cells to be checked.
 * @param neighbours A slice of coordinates representing the neighbours to filter.
 * @return A slice of coordinates of the fish cells found among the neighbours.
 */
func FilterFishCells(grid Grid, neighbours [][2]int) [][2]int {
	var fishCells [][2]int
	for _, n := range neighbours {
		if cell := grid[n[0]][n[1]]; cell != nil && cell.Type == Fish {
			fishCells = append(fishCells, n)
		}
	}
	return fishCells
}

/**
 * @brief Updates the simulation of the grid using multiple threads.
 *
 * This function creates a new grid to hold the updated state and divides
 * the work of updating the grid across multiple threads. Each thread processes
 * a specific range of rows, moving fish and sharks according to the rules
 * defined in the simulation.
 *
 * @param grid The current state of the grid containing entities (fish and sharks).
 * @param numThreads The number of threads to use for processing the grid update.
 */
func UpdateSimulation(grid Grid, numThreads int) {
	newGrid := make(Grid, GridSize)
	for i := range newGrid {
		newGrid[i] = make([]*Entity, GridSize)
	}

	var wg sync.WaitGroup
	rowsPerThread := GridSize / numThreads

	for i := 0; i < numThreads; i++ {
		startRow := i * rowsPerThread
		endRow := startRow + rowsPerThread
		if i == numThreads-1 {
			endRow = GridSize
		}

		wg.Add(1)
		go func(startRow, endRow int) {
			defer wg.Done()
			for x := startRow; x < endRow; x++ {
				for y, cell := range grid[x] {
					if cell == nil || (newGrid[x][y] != nil) {
						continue
					}

					if cell.Type == Fish {
						MoveFish(grid, newGrid, x, y)
					} else if cell.Type == Shark {
						MoveShark(grid, newGrid, x, y)
					}
				}
			}
		}(startRow, endRow)
	}

	wg.Wait()
	CopyGrid(grid, newGrid)
}

/**
 * @brief Copies the contents of one grid to another.
 *
 * This function iterates over the source grid and copies each cell's value
 * to the destination grid.
 *
 * @param dest The destination grid where the values will be copied to.
 * @param src The source grid from which values will be copied.
 */
func CopyGrid(dest, src Grid) {
	for i := range src {
		for j := range src[i] {
			dest[i][j] = src[i][j]
		}
	}
}

/**
 * @brief Draws the game grid onto the provided screen.
 *
 * This method fills the screen with a black background and draws each cell
 * of the grid based on its type. Cells representing fish and sharks are
 * drawn in green and red, respectively.
 *
 * @param screen A pointer to an `ebiten.Image` where the game grid will be drawn.
 */
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

/**
 * @brief Updates the simulation for the game.
 *
 * This method invokes the UpdateSimulation function with the game's grid
 * and the number of threads being used. It performs the necessary updates
 * to the game's state.
 *
 * @return nil if the update is successful; an error is returned if
 *         an issue occurs during the update process (note: currently
 *         it always returns nil).
 */func (g *Game) Update() error {
	UpdateSimulation(g.grid, g.numThreads)
	return nil
}

/**
 * @brief Sets the layout dimensions for the game screen.
 *
 * This method defines the layout for the game's screen based on the given
 * dimensions.
 *
 * @param outsideWidth The width of the game screen.
 * @param outsideHeight The height of the game screen.
 * @return The dimensions of the game screen as two integers,
 *         representing the width and height, respectively.
 */
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

/**
 * @brief Benchmarks the performance of a simulation using different thread counts
 *        and outputs the results to a XLSX file.
 *
 * This function performs the following steps:
 * 1. Creates a XLSX file specified by `outputFile`.
 * 2. Writes a header row to the XLSX file.
 * 3. For each thread count specified in `threadCounts`, it initializes
 *    a grid, runs the simulation for the given number of steps, and
 *    measures the elapsed time.
 * 4. Calculates the speedup compared to the baseline (the time taken
 *    with the first thread count).
 * 5. Writes the thread count, execution time, and speedup to the XLSX file.
 *
 * @param steps The number of simulation steps to execute.
 * @param threadCounts A slice of integers representing the number of threads
 *                     to benchmark.
 * @param outputFile The path to the XLSX file where the results will be saved.
 */
func BenchmarkSimulationToXLSX(steps int, threadCounts []int, outputFile string) {
	// Create a new Excel file
	f := excelize.NewFile()

	// Create a new sheet and handle both return values
	index, err := f.NewSheet("Benchmark Results")
	if err != nil {
		fmt.Printf("Error creating sheet: %v\n", err)
		return
	}

	// Write the header
	f.SetCellValue("Benchmark Results", "A1", "Threads")
	f.SetCellValue("Benchmark Results", "B1", "Time (seconds)")
	f.SetCellValue("Benchmark Results", "C1", "Speedup")

	// Benchmark each thread count
	baselineTime := 0.0
	for i, numThreads := range threadCounts {
		grid := InitialiseGrid() // Ensure you have this function defined
		startTime := time.Now()

		for step := 0; step < steps; step++ {
			UpdateSimulation(grid, numThreads) // Ensure you have this function defined
		}

		duration := time.Since(startTime).Seconds()
		if i == 0 {
			baselineTime = duration
		}

		speedup := baselineTime / duration

		// Write the result to Excel
		f.SetCellValue("Benchmark Results", fmt.Sprintf("A%d", i+2), numThreads)
		f.SetCellValue("Benchmark Results", fmt.Sprintf("B%d", i+2), duration)
		f.SetCellValue("Benchmark Results", fmt.Sprintf("C%d", i+2), speedup)
	}

	// Set the active sheet
	f.SetActiveSheet(index)

	// Save the Excel file
	if err := f.SaveAs(outputFile); err != nil {
		fmt.Printf("Failed to save file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Results saved to %s\n", outputFile)
}

/**
 * @brief Main function to start simulation and benchmarking.
 *
 * This function initializes the simulation environment and begins the
 * benchmarking process.
 *
 * @return int Returns 0 on successful completion.
 */
func RunSimulation() {
	rand.Seed(time.Now().UnixNano())

	threadCounts := []int{1, 2, 4, 8}      // Thread configurations to test
	steps := 100                           // Number of steps for benchmarking
	outputFile := "benchmark_results.xlsx" // File to save results to

	fmt.Println("Benchmarking Wator Simulation:")
	BenchmarkSimulationToXLSX(steps, threadCounts, outputFile)

	game := &Game{
		grid:       InitialiseGrid(),
		numThreads: 1,
	}

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Wator Simulation")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
