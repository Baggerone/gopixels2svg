package pixels2svg

import (
	"fmt"
	"sort"
)

const N = int(0)
const NE = int(1)
const E = int(2)
const SE = int(3)
const S = int(4)
const SW = int(5)
const W = int(6)
const NW = int(7)


type Shape struct {
	References map[int][2]int  // Column references as keys and top-most row reference with bottom-most row ref.
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

func isAdjacentCellInBounds(shapeExtr *ShapeExtractor, column, row int, direction int) bool {
	switch direction {
	case N: // North
		return ! shapeExtr.cellIsAtTop(row)

	case NE: // Northeast
		if shapeExtr.cellIsAtTop(row) {
			return false
		}
		return ! shapeExtr.cellIsAtRight(column)

	case E: // East
		return ! shapeExtr.cellIsAtRight(column)

	case SE: // Southeast
		if shapeExtr.cellIsAtBottom(row) {
			return false
		}
		return ! shapeExtr.cellIsAtRight(column)

	case S: // South
		return ! shapeExtr.cellIsAtBottom(row)

	case SW: // Southwest
		if shapeExtr.cellIsAtBottom(row) {
			return false
		}
		return ! shapeExtr.cellIsAtLeft(column)

	case W: // West
		return ! shapeExtr.cellIsAtLeft(column)

	case NW: // Northwest
		if shapeExtr.cellIsAtTop(row) {
			return false
		}

		return ! shapeExtr.cellIsAtLeft(column)
	}

	panic(fmt.Sprintf(
		"Error: direction should be 0, 1, 2, 3, 4, 5, 6, 7.  Was passed: %d",
		direction,
	))
	return false
}

func isSameColorAdjacent(shapeExtr *ShapeExtractor, column, row int, shapeCell GridCell, directions ...int) bool {
	for _, direction := range directions {
		if ! isAdjacentCellInBounds(shapeExtr, column, row, direction) {
			continue
		}
		adjacentColumn, adjacentRow := getCellReferenceInDirection(column, row, direction)
		if shapeCell.doesCellMatch(shapeExtr.grid[adjacentColumn][adjacentRow]) {
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
func makesTriangleToRight(shapeExtr *ShapeExtractor, column, row int, shapeCell GridCell) bool {
	if shapeExtr.cellIsAtRight(column) {
		return false
	}

	_, southRow := getCellReferenceInDirection(column, row, S)
	southCell := shapeExtr.grid[column][southRow]

	eastColumn, _ := getCellReferenceInDirection(column, row, E)
	eastCell := shapeExtr.grid[eastColumn][row]

	southEastCell := shapeExtr.grid[eastColumn][southRow]

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
func makesTriangleToLowerLeft(shapeExtr *ShapeExtractor, column, row int, shapeCell GridCell) bool {
	if shapeExtr.cellIsAtLeft(column) {
		return false
	}

	_, southRow := getCellReferenceInDirection(column, row, S)
	southCell := shapeExtr.grid[column][southRow]

	southWestColumn, _ := getCellReferenceInDirection(column, row, SW)
	southWestCell := shapeExtr.grid[southWestColumn][southRow]

	return shapeCell.doesCellMatch(southCell) && shapeCell.doesCellMatch(southWestCell)
}


// Invalid if there are no same-colored contiguous cells that make a tiny triangle with it
//  and if it has not already been dealt with
// Assumes that there is not a same-colored cell to the left that is useable
func isStartPositionValid(shapeExtr *ShapeExtractor, column, row int, shapeCell GridCell) bool {
	if shapeExtr.grid[column][row].AlreadyUsed {
		return false
	}
	if shapeExtr.cellIsAtBottom(row) {
		return false
	}
	return makesTriangleToRight(shapeExtr, column, row, shapeCell) ||
		makesTriangleToLowerLeft(shapeExtr, column, row, shapeCell)
}

// Assumes we are not on the eastern edge of the grid
// Start with the cell to the south of the starting cell, if it is the same color.
// As long as there is a cell of the same color to the NE of that cell,
// keep going down the column.
// The lowest contiguous same-colored cell that also has a same-colored cell to
// its NE is the one we want.
func findRowOfLowerCellInStartingColumn(shapeExtr *ShapeExtractor, startColumn, startRow int, shapeColor Color) int {
	if shapeExtr.cellIsAtBottom(startRow) {
		return startRow
	}

	cellType := GridCell{Color: shapeColor}
	goodLowRow := startRow
	_, nextRow := getCellReferenceInDirection(startColumn, startRow, S)

	for {
		// Only interested in cells of the same color
		if !cellType.doesCellMatch(shapeExtr.grid[startColumn][nextRow]) {
			return goodLowRow
		}

		if ! isSameColorAdjacent(shapeExtr, startColumn, nextRow, cellType, E, W) {
			return nextRow
		}

		// the cell to the northeast of the cell under consideration is the same color,
		// so if the cell under consideration is at the bottom of the grid, that's the one we want.
		if shapeExtr.cellIsAtBottom(nextRow) {
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

// If the provisional starting cell of the new column has a cell of the same color to its southeast, then
// his finds the highest contiguous same-colored cell in the new column that has
// a cell of the same color to its southeast.
// Otherwise, it finds the highest contiguous same-colored cell below it that has a cell of the same color
// to its east.
//
// Assumes column is not at the right edge of the grid
func findUpperRowForNextColumn(shapeExtr *ShapeExtractor, column, startRow int) int {

	cellInFocus := shapeExtr.grid[column][startRow]

	southeastColumn, southeastRow := getCellReferenceInDirection(column, startRow, SE)
	southeastCell := shapeExtr.grid[southeastColumn][southeastRow]

	oldRow := startRow

	// If the color of the cell to the southeast is the same, then keep rising up the
	// current column finding the highest contiguous same-colored cell that also has a
	// cell of the same color to its southeast
	if cellInFocus.doesCellMatch(southeastCell) {
		for {
			if shapeExtr.cellIsAtTop(oldRow) {
				return oldRow
			}

			_, rowAbove := getCellReferenceInDirection(column, oldRow, N)
			if !cellInFocus.doesCellMatch(shapeExtr.grid[column][rowAbove]) {
				return oldRow
			}

			southEastColumn, southEastRow := getCellReferenceInDirection(column, rowAbove, SE)
			southEastCell := shapeExtr.grid[southEastColumn][southEastRow]

			if !cellInFocus.doesCellMatch(southEastCell) {
				return oldRow
			}

			oldRow = rowAbove
		}
	}

	// Move down the current column until there is a cell of a different color or one that
	// has a cell to the southeast of it of the same color
	for {
		if shapeExtr.cellIsAtBottom(oldRow) {
			return oldRow
		}

		_, rowBelow := getCellReferenceInDirection(column, oldRow, S)
		if !cellInFocus.doesCellMatch(shapeExtr.grid[column][rowBelow]) {
			return oldRow
		}

		if shapeExtr.cellIsAtBottom(rowBelow) {
			return rowBelow
		}

		southEastColumn, southEastRow := getCellReferenceInDirection(column, rowBelow, SE)
		if cellInFocus.doesCellMatch(shapeExtr.grid[southEastColumn][southEastRow]) {
			return rowBelow
		}

		oldRow = rowBelow
	}
}


// Assumes that starting cell is not of the same color.
// Looks down the column from there for a cell of the same color, with some limitations
func getUpperRowOfNextColumnWhenItsLowerThanTheStartCell(
	shapeExtr *ShapeExtractor,
	column, startRow, lowestRow int,
	shapeCell GridCell,
) (int, error) {
	for row := startRow + 1; row <= lowestRow; row++ {
		if shapeCell.doesCellMatch(shapeExtr.grid[column][row]) {
			return row, nil
		}

		if shapeExtr.cellIsAtBottom(row) {
			break
		}
	}

	return 0, fmt.Errorf("No valid cells to the right")
}

func getUpperRowOfNextColumn(
	shapeExtr *ShapeExtractor,
	columnIsToEast bool,
	column, startRow, lowestRow int,
	shapeCell GridCell,
) (int, error) {

	if !shapeCell.doesCellMatch(shapeExtr.grid[column][startRow]) {
		return getUpperRowOfNextColumnWhenItsLowerThanTheStartCell(
			shapeExtr,
			column,
			startRow,
			lowestRow,
			shapeCell,
		)
	}

	cellIsAtEdge := shapeExtr.cellIsAtRightOrLeft(columnIsToEast, column)

	// Starting cell is the same color.  If this column is on the east edge, then go up no more than one row
	if cellIsAtEdge {
		if isSameColorAdjacent(shapeExtr, column, startRow, shapeCell, N) {
			_, rowToNorth := getCellReferenceInDirection(column, startRow, N)
			return rowToNorth, nil
		}

		return startRow, nil
	}


	previousRow := startRow

	// Starting cell is the same color.
	// If the cell above the starting cell is the same color, add it.
	if isSameColorAdjacent(shapeExtr, column, startRow, shapeCell, N) {
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
		if !shapeCell.doesCellMatch(shapeExtr.grid[column][row]) {
			return previousRow, nil
		}

		if !isSameColorAdjacent(shapeExtr, column, row, shapeCell, S) ||
			!isSameColorAdjacent(shapeExtr, column, row, shapeCell, checkDirection) {
			return previousRow, nil
		}

		previousRow = row
	}

	row := 0
	if !isSameColorAdjacent(shapeExtr, column, row, shapeCell, S) ||
		!isSameColorAdjacent(shapeExtr, column, row, shapeCell, checkDirection) {
		return previousRow, nil
	}

	return row, nil
}


func getLowerRowOfNextColumn(
	shapeExtr *ShapeExtractor,
	columnIsToEast bool,
	column, upperRow, lowestRow int,
	shapeCell GridCell,
) int {

	previousRow := upperRow

	for row := upperRow; row <= lowestRow; row++ {
		if ! isSameColorAdjacent(shapeExtr, column, row, shapeCell, S) {
			return row
		}

		previousRow = row
	}

	cellIsAtEdge := shapeExtr.cellIsAtRightOrLeft(columnIsToEast, column)

	// If this column is on the east [or west] edge, then go down no more than one row
	if cellIsAtEdge {
		if isSameColorAdjacent(shapeExtr, column, previousRow, shapeCell, S) {
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
	for row := previousRow + 1; row < shapeExtr.RowCount; row++ {
		if ! isSameColorAdjacent(shapeExtr, column, row, shapeCell, S) {
			return row
		}

		if ! isSameColorAdjacent(shapeExtr, column, row, shapeCell, checkDirection) {
			return row
		}

		newLowestRow = row
	}

	return newLowestRow
}

func getShapeColumnsToOneSide(
	shapeExtr *ShapeExtractor,
	lookingToEast bool,
	previousColumn,
	startRow, lowestRow int,
	shape Shape,
) Shape {

	if startRow >= lowestRow {
		return shape
	}

	cellIsAtEdge := shapeExtr.cellIsAtRightOrLeft(lookingToEast, previousColumn)
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
			shapeExtr,
			lookingToEast,
			nextColumn,
			startRow,
			nextLowerRow,
			shapeCell,
		)
		if err != nil {
			return shape
		}

		nextLowerRow = getLowerRowOfNextColumn(
			shapeExtr,
			lookingToEast,
			nextColumn,
			nextUpperRow,
			nextLowerRow,
			shapeCell,
		)

		shape.References[nextColumn] = [2]int{nextUpperRow, nextLowerRow}

		if shapeExtr.cellIsAtRightOrLeft(lookingToEast, nextColumn) {
			break
		}

		if nextUpperRow >= nextLowerRow {
			break
		}
	}

	return shape
}

func getUpperRowOfColumnToRight(shapeExtr *ShapeExtractor, column, startRow, lowestRow int, shapeCell GridCell) (int, error) {
	return getUpperRowOfNextColumn(shapeExtr, true, column, startRow, lowestRow, shapeCell)
}

func getLowerRowOfColumnToRight(shapeExtr *ShapeExtractor, column, upperRow, lowestRow int, shapeCell GridCell) int {
	return getLowerRowOfNextColumn(shapeExtr, true, column, upperRow, lowestRow, shapeCell)
}

func getUpperRowOfColumnToLeft(shapeExtr *ShapeExtractor, column, startRow, lowestRow int, shapeCell GridCell) (int, error) {
	return getUpperRowOfNextColumn(shapeExtr, false, column, startRow, lowestRow, shapeCell)
}

func getLowerRowOfColumnToLeft(shapeExtr *ShapeExtractor, column, upperRow, lowestRow int, shapeCell GridCell) int {
	return getLowerRowOfNextColumn(shapeExtr, false, column, upperRow, lowestRow, shapeCell)
}


func getShapeColumnsToRight(shapeExtr *ShapeExtractor, previousColumn, startRow, lowestRow int, shape Shape) Shape {
	return getShapeColumnsToOneSide(shapeExtr, true, previousColumn, startRow, lowestRow, shape)
}

func getShapeColumnsToLeft(shapeExtr *ShapeExtractor, previousColumn, startRow, lowestRow int, shape Shape) Shape {
	return getShapeColumnsToOneSide(shapeExtr, false, previousColumn, startRow, lowestRow, shape)
}

func getShapeStartingAtCellReference(shapeExtr *ShapeExtractor, startColumn, startRow int) Shape {
	shape := Shape{
		References: map[int][2]int{},
		Color: shapeExtr.grid[startColumn][startRow].Color,
	}

	shapeCell := GridCell{Color:shape.Color}

	if ! isStartPositionValid(shapeExtr, startColumn, startRow, shapeCell) {
		return shape
	}
	startColumnLowerRow := findRowOfLowerCellInStartingColumn(
		shapeExtr,
		startColumn,
		startRow,
		shape.Color,
	)

	// Record startcolumn in the shape
	shape.References[startColumn] = [2]int{startRow, startColumnLowerRow}

	shape = getShapeColumnsToRight(shapeExtr, startColumn, startRow, startColumnLowerRow, shape)
	shape = getShapeColumnsToLeft(shapeExtr, startColumn, startRow, startColumnLowerRow, shape)

	return shape
}

func getPolygonFromShape(shape Shape) (Polygon, error) {

	polygonRefs := [][2]int{}
	if len(shape.References) <= 1 {
		return Polygon{}, fmt.Errorf("More than one column of cells is required.")
	}

	columnRefs := []int{}
	for column, _ := range shape.References {
		columnRefs = append(columnRefs, column)
	}

	sort.Ints(columnRefs)
	firstColumn := columnRefs[0]
	firstTopRow := shape.References[firstColumn][0]
	lastColumn := columnRefs[len(columnRefs) - 1]

	polygonRefs = append(polygonRefs, [2]int{firstColumn, shape.References[firstColumn][0]})

	nextColumn := columnRefs[1]
	nextTopRow := shape.References[nextColumn][0]

	// If the top of the next column is more than one higher than the top of the first column,
	//   add the point next to the top of the first column but one up
	if nextTopRow > firstTopRow + 1 {
		polygonRefs = append(polygonRefs, [2]int{nextColumn, firstTopRow - 1})
	}

   // Start at left most column and move to the right along the top of the shape
	for colIndex := firstColumn + 1; colIndex < lastColumn; colIndex++ {
		currentTopRow := shape.References[colIndex][0]
		polygonRefs = append(polygonRefs, [2]int{colIndex, currentTopRow})

		nextTopRow := shape.References[colIndex + 1][0]

		// If the top of the next column is more than one lower than the top of the current column,
		//  add the point on the current column that is one row lower than the top of the next column
		if nextTopRow > currentTopRow + 1 {
			polygonRefs = append(polygonRefs, [2]int{colIndex, nextTopRow - 1})

		// If the top of the next column is more than one higher than the top of the current column,
		//   add the point next to the top of the current column but one up
		} else if nextTopRow < currentTopRow - 1 {
			polygonRefs = append(polygonRefs, [2]int{colIndex + 1, currentTopRow - 1})
		}
	}

	// Add the points from the last column
	polygonRefs = append(polygonRefs, [2]int{lastColumn, shape.References[lastColumn][0]})
	polygonRefs = append(polygonRefs, [2]int{lastColumn, shape.References[lastColumn][1]})

	lastBottomRow := shape.References[lastColumn][1]

	nextColumn = columnRefs[len(columnRefs) - 2]
	nextBottomRow := shape.References[nextColumn][1]

	// If the bottom of the second to right column is more than one lower than the bottom of the right most column,
	//   add the point next to the bottom of the right most column but one lower
	if nextBottomRow > lastBottomRow + 1 {
		polygonRefs = append(polygonRefs, [2]int{nextColumn, lastBottomRow + 1})
	}

	// Start at right most column and move to the left along the bottom of the shape
	for colIndex := lastColumn - 1; colIndex > firstColumn; colIndex-- {
		currentBottomRow := shape.References[colIndex][1]
		polygonRefs = append(polygonRefs, [2]int{colIndex, currentBottomRow})

		nextBottomRow := shape.References[colIndex - 1][1]

		// If the bottom of the next column (to the left) is more than one higher than the bottom of the current column,
		//  add the point on the current column that is one row lower than the bottom of the next column
		if nextBottomRow < currentBottomRow - 1 {
			polygonRefs = append(polygonRefs, [2]int{colIndex, nextBottomRow + 1})

			// If the bottom of the next column is more than one lower than the bottom of the current column,
			//   add the point next to the bottom of the current column but one up
		} else if nextBottomRow > currentBottomRow + 1 {
			polygonRefs = append(polygonRefs, [2]int{colIndex - 1, currentBottomRow + 1})
		}
	}

	polygonRefs = append(polygonRefs, [2]int{firstColumn, shape.References[firstColumn][1]})

	polygon := Polygon{
		ColorRGBA: shape.Color,
		Points: polygonRefs,
	}
	return polygon, nil
}