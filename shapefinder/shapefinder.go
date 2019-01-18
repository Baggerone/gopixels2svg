package shapefinder

import "fmt"

const N = int(0)
const NE = int(1)
const E = int(2)
const SE = int(3)
const S = int(4)
const SW = int(5)
const W = int(6)
const NW = int(7)


type Color [4]uint8

type GridCell struct {
	AlreadyUsed bool
	Color       Color
}

func (g GridCell) doesCellMatch(g2 GridCell) bool {
	return !g2.AlreadyUsed && g.Color == g2.Color
}

type Grid [][]GridCell


func (g Grid) cellIsAtRightOrLeft(checkToTheRight bool, column, row int) bool {
	cellIsAtEdge := false
	if checkToTheRight {
		cellIsAtEdge = g.cellIsAtRight(column, row)
	} else {
		cellIsAtEdge = g.cellIsAtLeft(column, row)
	}
	return cellIsAtEdge
}

func (g Grid)cellIsAtRight(column, row int) bool {
	return column >= len(g) - 1
}

func (g Grid)cellIsAtLeft(column, row int) bool {
	return column <= 0
}

func (g Grid)cellIsAtTop(column, row int) bool {
	return row <= 0
}

func (g Grid)cellIsAtBottom(column, row int) bool {
	return row >= len(g[0]) - 1
}


type Shape struct {
	References map[int][]int
	Color Color
}

/*
 * Given a cell and a direction what are the cell cordinates of the neighboring
 * cell in that direction.
 * Assumes that boundary checks are already done, so that the starting cell
 * won't be on a boundary
 */
func getCellReferenceInDirection(columnX, rowY, nextDirection int) (int, int) {
	switch nextDirection {
	case N: // North
		return columnX, rowY - 1
	case NE: // Northeast
		return columnX + 1, rowY - 1
	case E: // East
		return columnX + 1, rowY
	case SE: // Southeast
		return columnX + 1, rowY + 1
	case S: // South
		return columnX, rowY + 1
	case SW: // Southwest
		return columnX - 1, rowY + 1
	case W: // West
		return columnX - 1, rowY
	case NW: // Northwest
		return columnX - 1, rowY - 1
	}

	panic(fmt.Sprintf(
		"Error: direction should be 0, 1, 2, 3, 4, 5, 6, 7.  Was passed: %d",
		nextDirection,
	))
	return 0, 0
}

func isAdjacentCellInBounds(column, row int, grid Grid, direction int) bool {
	switch direction {
	case N: // North
		return ! grid.cellIsAtTop(column, row)

	case NE: // Northeast
		if grid.cellIsAtTop(column, row) {
			return false
		}
		return ! grid.cellIsAtRight(column, row)

	case E: // East
		return ! grid.cellIsAtRight(column, row)

	case SE: // Southeast
		if grid.cellIsAtBottom(column, row) {
			return false
		}
		return ! grid.cellIsAtRight(column, row)

	case S: // South
		return ! grid.cellIsAtBottom(column, row)

	case SW: // Southwest
		if grid.cellIsAtBottom(column, row) {
			return false
		}
		return ! grid.cellIsAtLeft(column, row)

	case W: // West
		return ! grid.cellIsAtLeft(column, row)

	case NW: // Northwest
		if grid.cellIsAtTop(column, row) {
			return false
		}

		return ! grid.cellIsAtLeft(column, row)
	}

	panic(fmt.Sprintf(
		"Error: direction should be 0, 1, 2, 3, 4, 5, 6, 7.  Was passed: %d",
		direction,
	))
	return false
}

func isSameColorAdjacent(column, row int, grid Grid, shapeCell GridCell, directions ...int) bool {
	for _, direction := range directions {
		if ! isAdjacentCellInBounds(column, row, grid, direction) {
			continue
		}
		adjacentColumn, adjacentRow := getCellReferenceInDirection(column, row, direction)
		if shapeCell.doesCellMatch(grid[adjacentColumn][adjacentRow]) {
			return true
		}
	}

	return false
}


//
// _ _ _ _ _      _ _ _ _ _     _ _ _ _ _
// _ _ G G _      _ _ G _ _     _ _ G G _
// _ _ G _ _  or  _ _ G G _ or  _ _ _ G _
// _ _ _ _ _      _ _ _ _ _     _ _ _ _ _
//
// True if cells either
//   - to the south and east are the same color or
//   - to the south and southeast are the same color or
//   - to the east and southeast are the same color
// Assumes the start position is not on the bottom row
func makesTriangleToRight(column, row int, grid Grid, shapeCell GridCell) bool {
	if grid.cellIsAtRight(column, row) {
		return false
	}

	_, southRow := getCellReferenceInDirection(column, row, S)
	southCell := grid[column][southRow]

	eastColumn, _ := getCellReferenceInDirection(column, row, E)
	eastCell := grid[eastColumn][row]

	southEastCell := grid[eastColumn][southRow]

	if shapeCell.doesCellMatch(southCell) {
		return shapeCell.doesCellMatch(eastCell) || shapeCell.doesCellMatch(southEastCell)
	}

	return shapeCell.doesCellMatch(eastCell) && shapeCell.doesCellMatch(southEastCell)
}

//
// _ _ _ _ _
// _ _ G _ _
// _ G G _ _
// _ _ _ _ _
//
// True if the cells to the south and southwest are the same color
// Assumes cell to the west has already been dealt with and is not usable
// Assumes the start position is not on the bottom row
func makesTriangleToLowerLeft(column, row int, grid Grid, shapeCell GridCell) bool {
	if grid.cellIsAtLeft(column, row) {
		return false
	}

	_, southRow := getCellReferenceInDirection(column, row, S)
	southCell := grid[column][southRow]

	southWestColumn, _ := getCellReferenceInDirection(column, row, SW)
	southWestCell := grid[southWestColumn][southRow]

	return shapeCell.doesCellMatch(southCell) && shapeCell.doesCellMatch(southWestCell)
}


// Invalid if there are no same-colored contiguous cells that make a tiny triangle with it
//  and if it has not already been dealt with
// Assumes that there is not a same-colored cell to the left that is useable
func isStartPositionValid(column, row int, grid Grid, shapeCell GridCell) bool {
	if grid[column][row].AlreadyUsed {
		return false
	}
	if grid.cellIsAtBottom(column, row) {
		return false
	}
	return makesTriangleToRight(column, row, grid, shapeCell) ||
		makesTriangleToLowerLeft(column, row, grid, shapeCell)
}

// Assumes we are not on the eastern edge of the grid
// Start with the cell to the south of the starting cell, if it is the same color.
// As long as there is a cell of the same color to the NE of that cell,
// keep going down the column.
// The lowest contiguous same-colored cell that also has a same-colored cell to
// its NE is the one we want.
func findRowOfLowerCellInStartingColumn(startColumn, startRow int, grid Grid, shapeColor Color) int {
	if grid.cellIsAtBottom(startColumn, startRow) {
		return startRow
	}

	cellType := GridCell{Color: shapeColor}
	goodLowRow := startRow
	_, nextRow := getCellReferenceInDirection(startColumn, startRow, S)

	for {
		// Only interested in cells of the same color
		if !cellType.doesCellMatch(grid[startColumn][nextRow]) {
			return goodLowRow
		}

		if ! isSameColorAdjacent(startColumn, nextRow, grid, cellType, E, W) {
			return nextRow
		}

		// the cell to the northeast of the cell under consideration is the same color,
		// so if the cell under consideration is at the bottom of the grid, that's the one we want.
		if grid.cellIsAtBottom(startColumn, nextRow) {
		    return nextRow
		}

		goodLowRow = nextRow
		_, nextRow = getCellReferenceInDirection(startColumn, nextRow, S)
	}
}

// Assumes upperRow is above lowerRow and the cells in column,upperRow and column,lowerRow are
// the same color
func isSubColumnOneColor(column, upperRow, lowerRow int, grid Grid) bool {
	colorInFocus := grid[column][upperRow]

	for i := upperRow + 1; i < lowerRow; i++ {
		if grid[column][i] != colorInFocus {
			return false
		}
	}
	return true
}

func addSubColumnToShape(column, startRow, endRow int, shape Shape) Shape {
	_, ok := shape.References[column]
	if !ok {
		shape.References[column] = []int{}
	}

	for i := startRow; i <= endRow; i++ {
		shape.References[column] = append(shape.References[column], i)
	}

	return shape
}

// If the provisional starting cell of the new column has a cell of the same color to its southeast, then
// his finds the highest contiguous same-colored cell in the new column that has
// a cell of the same color to its southeast.
// Otherwise, it finds the highest contiguous same-colored cell below it that has a cell of the same color
// to its east.
//
// Assumes column is not at the right edge of the grid
func findUpperRowForNextColumn(column, startRow int, grid Grid) int {

	cellInFocus := grid[column][startRow]

	southeastColumn, southeastRow := getCellReferenceInDirection(column, startRow, SE)
	southeastCell := grid[southeastColumn][southeastRow]

	oldRow := startRow

	// If the color of the cell to the southeast is the same, then keep rising up the
	// current column finding the highest contiguous same-colored cell that also has a
	// cell of the same color to its southeast
	if cellInFocus.doesCellMatch(southeastCell) {
		for {
			if grid.cellIsAtTop(column, oldRow) {
				return oldRow
			}

			_, rowAbove := getCellReferenceInDirection(column, oldRow, N)
			if !cellInFocus.doesCellMatch(grid[column][rowAbove]) {
				return oldRow
			}

			southEastColumn, southEastRow := getCellReferenceInDirection(column, rowAbove, SE)
			southEastCell := grid[southEastColumn][southEastRow]

			if !cellInFocus.doesCellMatch(southEastCell) {
				return oldRow
			}

			oldRow = rowAbove
		}
	}

	// Move down the current column until there is a cell of a different color or one that
	// has a cell to the southeast of it of the same color
	for {
		if grid.cellIsAtBottom(column, oldRow) {
			return oldRow
		}

		_, rowBelow := getCellReferenceInDirection(column, oldRow, S)
		if !cellInFocus.doesCellMatch(grid[column][rowBelow]) {
			return oldRow
		}

		if grid.cellIsAtBottom(column, rowBelow) {
			return rowBelow
		}

		southEastColumn, southEastRow := getCellReferenceInDirection(column, rowBelow, SE)
		if cellInFocus.doesCellMatch(grid[southEastColumn][southEastRow]) {
			return rowBelow
		}

		oldRow = rowBelow
	}
}


// Assumes that starting cell is not of the same color.
// Looks down the column from there for a cell of the same color, with some limitations
func getUpperRowOfNextColumnWhenItsLowerThanTheStartCell(
	column,
	startRow,
	lowestRow int,
	grid Grid,
	shapeCell GridCell,
) (int, error) {
	for row := startRow + 1; row <= lowestRow; row++ {
		if shapeCell.doesCellMatch(grid[column][row]) {
			return row, nil
		}

		if grid.cellIsAtBottom(column, row) {
			break
		}
	}

	return 0, fmt.Errorf("No valid cells to the right")
}

func getUpperRowOfNextColumn(
	columnIsToEast bool,
	column, startRow, lowestRow int,
	grid Grid,
	shapeCell GridCell,
) (int, error) {

	if !shapeCell.doesCellMatch(grid[column][startRow]) {
		return getUpperRowOfNextColumnWhenItsLowerThanTheStartCell(
			column,
			startRow,
			lowestRow,
			grid,
			shapeCell,
		)
	}

	cellIsAtEdge := grid.cellIsAtRightOrLeft(columnIsToEast, column, startRow)

	// Starting cell is the same color.  If this column is on the east edge, then go up no more than one row
	if cellIsAtEdge {
		if isSameColorAdjacent(column, startRow, grid, shapeCell, N) {
			_, rowToNorth := getCellReferenceInDirection(column, startRow, N)
			return rowToNorth, nil
		}

		return startRow, nil
	}


	previousRow := startRow

	// Starting cell is the same color.
	// If the cell above the starting cell is the same color, add it.
	if isSameColorAdjacent(column, startRow, grid, shapeCell, N) {
		previousRow = startRow - 1
	}

	checkDirection := SW
	if columnIsToEast {
		checkDirection = SE
	}

	// Starting cell is the same color and column is not on east edge, so look above it up the column
	// until you get a different color or until the cell to the right is a different color
	// (Avoid a single column "chimney")
	for row := startRow - 2; row > 0; row-- {
		if !shapeCell.doesCellMatch(grid[column][row]) {
			return previousRow, nil
		}

		if !isSameColorAdjacent(column, row, grid, shapeCell, S) ||
			!isSameColorAdjacent(column, row, grid, shapeCell, checkDirection) {
			return previousRow, nil
		}

		previousRow = row
	}

	row := 0
	if !isSameColorAdjacent(column, row, grid, shapeCell, S) ||
		!isSameColorAdjacent(column, row, grid, shapeCell, checkDirection) {
		return previousRow, nil
	}

	return row, nil
}


func getLowerRowOfNextColumn(
	columnIsToEast bool,
	column, upperRow, lowestRow int,
	grid Grid,
	shapeCell GridCell,
) int {

	previousRow := upperRow

	for row := upperRow; row <= lowestRow; row++ {
		if ! isSameColorAdjacent(column, row, grid, shapeCell, S) {
			return row
		}

		previousRow = row
	}

	cellIsAtEdge := grid.cellIsAtRightOrLeft(columnIsToEast, column, previousRow)

	// If this column is on the east [or west] edge, then go down no more than one row
	if cellIsAtEdge {
		if isSameColorAdjacent(column, previousRow, grid, shapeCell, S) {
			_, rowToSouth := getCellReferenceInDirection(column, previousRow, S)
			return rowToSouth
		}

		return previousRow
	}

	newLowestRow := previousRow

	checkDirection := W
	if columnIsToEast {
		checkDirection = E
	}

	// This column is not on the east [or west] edge
	// All same color from upperRow down to the lowestRow of the same-colored subcolumn to the left [or right],
	// so look below it down the column until you get a different color or until the cell to the right [or left]
	// is a different color.
	// (Avoid a single column "stalactite")
	for row := previousRow + 1; row < len(grid[0]); row++ {
		if ! isSameColorAdjacent(column, row, grid, shapeCell, S) {
			return row
		}

		if ! isSameColorAdjacent(column, row, grid, shapeCell, checkDirection) {
			return row
		}

		newLowestRow = row
	}

	return newLowestRow
}

func getShapeColumnsToOneSide(
	lookingToEast bool,
	previousColumn,
	startRow, lowestRow int,
	grid Grid,
	shape Shape,
) Shape {

	if startRow >= lowestRow {
		return shape
	}

	cellIsAtEdge := grid.cellIsAtRightOrLeft(lookingToEast, previousColumn, startRow)
	if cellIsAtEdge {
		return shape
	}

	nextColumn := previousColumn
	direction := W
	if lookingToEast {
		direction = E
	}
	nextLowerRow := lowestRow
	shapeCell := GridCell{Color: shape.Color}

	for {
		nextColumn, _ = getCellReferenceInDirection(nextColumn, startRow, direction)

		nextUpperRow, err := getUpperRowOfNextColumn(
			lookingToEast,
			nextColumn,
			startRow,
			nextLowerRow,
			grid,
			shapeCell,
		)
		if err != nil {
			return shape
		}

		nextLowerRow = getLowerRowOfNextColumn(
			lookingToEast,
			nextColumn,
			nextUpperRow,
			nextLowerRow,
			grid,
			shapeCell,
		)
		shape = addSubColumnToShape(nextColumn, nextUpperRow, nextLowerRow, shape)

		if grid.cellIsAtRightOrLeft(lookingToEast, nextColumn, startRow) {
			break
		}

		if nextUpperRow >= nextLowerRow {
			break
		}
	}

	return shape
}

func getUpperRowOfColumnToRight(column, startRow, lowestRow int, grid Grid, shapeCell GridCell) (int, error) {
	return getUpperRowOfNextColumn(true, column, startRow, lowestRow, grid, shapeCell)
}

func getLowerRowOfColumnToRight(column, upperRow, lowestRow int, grid Grid, shapeCell GridCell) int {
	return getLowerRowOfNextColumn(true, column, upperRow, lowestRow, grid, shapeCell)
}

func getUpperRowOfColumnToLeft(column, startRow, lowestRow int, grid Grid, shapeCell GridCell) (int, error) {
	return getUpperRowOfNextColumn(false, column, startRow, lowestRow, grid, shapeCell)
}

func getLowerRowOfColumnToLeft(column, upperRow, lowestRow int, grid Grid, shapeCell GridCell) int {
	return getLowerRowOfNextColumn(false, column, upperRow, lowestRow, grid, shapeCell)
}


func getShapeColumnsToRight(previousColumn, startRow, lowestRow int, grid Grid, shape Shape) Shape {
	return getShapeColumnsToOneSide(true, previousColumn, startRow, lowestRow, grid, shape)
}

func getShapeColumnsToLeft(previousColumn, startRow, lowestRow int, grid Grid, shape Shape) Shape {
	return getShapeColumnsToOneSide(false, previousColumn, startRow, lowestRow, grid, shape)
}

func getShapeStartingAtCellReference(startColumn, startRow int, grid Grid) Shape {
	shape := Shape{
		References: map[int][]int{},
		Color: grid[startColumn][startRow].Color,
	}

	shapeCell := GridCell{Color:shape.Color}

	if ! isStartPositionValid(startColumn, startRow, grid, shapeCell) {
		return shape
	}
	startColumnLowerRow := findRowOfLowerCellInStartingColumn(
		startColumn,
		startRow,
		grid,
		shape.Color,
	)

	// Record startcolumn in the shape
	shape = addSubColumnToShape(startColumn, startRow, startColumnLowerRow, shape)

	shape = getShapeColumnsToRight(startColumn, startRow, startColumnLowerRow, grid, shape)
	shape = getShapeColumnsToLeft(startColumn, startRow, startColumnLowerRow, grid, shape)

	return shape
}