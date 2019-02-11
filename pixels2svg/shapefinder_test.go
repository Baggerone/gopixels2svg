package pixels2svg

import (
	"testing"
	"strings"
	"fmt"
	"sort"
	)

func Red() GridCell {
	return GridCell{Color: Color{220, 0, 0, 0}}
}

func Green() GridCell {
	return GridCell{Color: Color{0, 220, 0, 0}}
}

func Blue() GridCell {
	return GridCell{Color: Color{0, 0, 220, 0}}
}
func getRGBColorFromString(colorCode string) GridCell {

	colors := map[string]GridCell{
		"r": Red(),
		"g": Green(),
		"b": Blue(),
	}

	rgbColor, OK := colors[strings.Trim(colorCode, " ")]; if ! OK {
		return GridCell{Color: Color{0, 0, 0, 0}}
	}

	return rgbColor
}

func transposeGrid(grid Grid) Grid {
	newGrid := Grid{}
	rowCount := len(grid)
	colCount := len(grid[0])

	for i := 0; i < colCount; i++ {
		newCol := []GridCell{}
		for j := 0; j < rowCount; j++ {
			newCol = append(newCol, GridCell{})
		}
		newGrid = append(newGrid, newCol)
	}

	for rowIndex, row := range grid {
		for colIndex, colValue := range row {
			newGrid[colIndex][rowIndex].Color = colValue.Color
		}
	}
	return newGrid
}
func transposeGridUsed(gridUsed [][]bool) [][]bool {
	newGrid := [][]bool{}
	rowCount := len(gridUsed)
	colCount := len(gridUsed[0])

	for i := 0; i < colCount; i++ {
		newCol := []bool{}
		for j := 0; j < rowCount; j++ {
			newCol = append(newCol, false)
		}
		newGrid = append(newGrid, newCol)
	}

	for rowIndex, row := range gridUsed {
		for colIndex, colValue := range row {
			newGrid[colIndex][rowIndex] = colValue
		}
	}

	return newGrid
}

// initGrid expects a multi-line string of single characters separated by spaces
// representing colored cells
// Its return grid is a slice of columns each having a slice of colors for the rows
func initGrid(textGrid string, gridUsed [][]bool, t *testing.T) [][]GridCell {
	grid := [][]GridCell{}
	rows := strings.Split(textGrid, "\n")

	for _, row := range rows{
		gridRow := []GridCell{}
		cells := strings.Split(row, ".")
		if len(cells) <= 1 {
			continue
		}
		for _, cell := range cells {
			gridRow = append(gridRow, getRGBColorFromString(cell))
		}
		grid = append(grid, gridRow)
	}

	grid = transposeGrid(grid)

	if gridUsed == nil {
		return grid
	}

	gridUsed = transposeGridUsed(gridUsed)
	if len(grid) != len(gridUsed) || len(grid[0]) != len(gridUsed[0]) {
		t.Errorf(
			"initGrid got grids that don't correspond to each other:\n  %s \nvs.\n  %+v",
			textGrid,
			gridUsed,
		)
		t.Fail()
	}

	for colIndex, column := range gridUsed {
		for rowIndex, cell := range column {
			grid[colIndex][rowIndex].AlreadyUsed = cell
		}
	}

	return grid
}

func compareGrids(results, expected [][]GridCell) string {
	if len(results) != len(expected) {
		return fmt.Sprintf(
			"Wrong number of columns. Expected %d, but got %d.\n%v",
			len(expected),
			len(results),
			results,
			)
	}

	if len(expected) > 0 && len(expected[0]) != len(results[0]) {
		return fmt.Sprintf(
			"Wrong number of rows. Expected %d, but got %d.\n%v",
			len(expected[0]),
			len(results[0]),
			results,
		)
	}

	for colIndex, expectedCol := range expected {
		for rowIndex, expectedRowValue := range expectedCol {
			if expectedRowValue.Color != results[colIndex][rowIndex].Color {
				return fmt.Sprintf(
					"Wrong value at column %d row %d. Expected %v, but got %v.\n%v",
					colIndex,
					rowIndex,
					expectedRowValue,
					results[colIndex][rowIndex],
					results,
				)
			}
		}
	}

	return ""
}

func getShapeFromGrid(grid Grid, color Color) Shape {

	bottomRowIndex := len(grid[0]) - 1

	shapeRefs := map[int][2]int{}
	for colIndex, column := range grid {
		foundTop := false
		foundBottom := false
		topRow := -1
		bottomRow := -1

		for rowIndex, cell := range column {
			if !foundTop && cell.Color == color {
				foundTop = true
				topRow = rowIndex
				continue
			}

			if foundTop && cell.Color != color {
				foundBottom = true
				bottomRow = rowIndex - 1
				break
			}
		}

		if ! foundBottom {
			bottomRow = bottomRowIndex
		}

		if topRow >= 0 && bottomRow >= 0 {
			shapeRefs[colIndex] = [2]int{topRow, bottomRow}
		}

	}

	shape := Shape{
		References: shapeRefs,
		Color: color,
	}

	return shape
}

func compareShapeReferences(results, expected Shape) string {

	prettyOutput := "{"
	keys := []int{}
	for k := range results.References {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	for _, k := range keys {
		prettyOutput = fmt.Sprintf("%s\n %d: %v", prettyOutput, k, results.References[k])
	}

	prettyOutput = fmt.Sprintf("%s\n}", prettyOutput)

	if len(results.References) != len(expected.References) {
		return fmt.Sprintf(
			"Wrong number of columns. Expected %d, but got %d.\nGot: %s",
			len(expected.References),
			len(results.References),
			prettyOutput,
		)
	}

	for key, expectedCells := range expected.References {
		resultsCells := results.References[key]
		if len(expectedCells) != len(resultsCells) {
			return fmt.Sprintf(
				"Wrong number of cells in column %d. Expected %d, but got %d.\nGot: %s",
			    key,
				len(expectedCells),
				len(resultsCells),
				prettyOutput,
			)
		}
		for cellIndex, expectedCellValue := range expectedCells {
			if expectedCellValue != resultsCells[cellIndex] {
				return fmt.Sprintf(
					"Wrong cell value in column %d, row %d. Expected %d, but got %d.\nGot: %s",
					key,
					cellIndex,
					expectedCellValue,
					resultsCells[cellIndex],
					prettyOutput,
				)
			}
		}
	}

	return ""
}

func compareGridsAlreadyUsed(results, expected Grid) string {

	if len(results) != len(expected) {
		return fmt.Sprintf(
			"Wrong number of columns. Expected %d, but got %d.",
			len(expected),
			len(results),
		)
	}

	if len(expected) > 0 && len(results[0]) != len(expected[0]) {
		return fmt.Sprintf(
			"Wrong number of rows. Expected %d, but got %d.",
			len(expected[0]),
			len(results[0]),
		)
	}

	for colIndex, expectedColumn := range expected {
		for rowIndex, expectedCell := range expectedColumn {
			resultsCell := results[colIndex][rowIndex]
			if expectedCell.AlreadyUsed != resultsCell.AlreadyUsed {
				return fmt.Sprintf(
					"Wrong cell value in column %d, row %d. Expected %v, but got %v.",
					colIndex,
					rowIndex,
					expectedCell.AlreadyUsed,
					resultsCell.AlreadyUsed,
				)
			}
		}
	}
	return ""
}


func TestInitGrid(t *testing.T) {
	textGrid := `
r.g.b
r.r.g`

	gridUsed := [][]bool{
		{true, false, false},
		{true, true, false},
	}

    // The grid gets transposed to a list of columns that each have a list of row cell values
    results := initGrid(textGrid, gridUsed, t)

    trueRed := Red()
    trueRed.AlreadyUsed = true

    expected := Grid{
    	{trueRed, trueRed},
    	{Green(), trueRed},
    	{Blue(), Green()},
	}

	errorMsg := compareGrids(results, expected)
	if errorMsg != "" {
		t.Error(errorMsg)
	}
}

func TestIsAdjacentCellInBoundsFalse(t *testing.T) {
	textGrid := `
g.g.b
g.r.b
b.g.r`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	type testData struct {
		startCol int
		startRow int
		direction int
	}

	allTestData := []testData{
		{1, 0, N},
		{2, 2, NE},
		{0, 0, NE},
		{2, 1, E},
		{2, 0, SE},
		{0, 2, SE},
		{1, 2, S},
		{0, 0, SW},
		{2, 2, SW},
		{0, 1, W},
		{0, 2, NW},
		{2, 0, NW},
	}

	for index, data := range allTestData {
		results := isAdjacentCellInBounds(&shapeExtr, data.startCol, data.startRow, data.direction)
		if results {
			t.Errorf("Expected false result at index %d for %+v", index, data)
		}
	}
}

func TestIsAdjacentCellInBoundsTrue(t *testing.T) {
	textGrid := `
g.g.b
g.r.b
b.g.r`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	type testData struct {
		startCol int
		startRow int
		direction int
	}

	allTestData := []testData{
		{1, 1, N},
		{1, 2, NE},
		{0, 1, NE},
		{1, 2, E},
		{1, 0, SE},
		{0, 1, SE},
		{0, 1, S},
		{1, 0, SW},
		{2, 1, SW},
		{1, 0, W},
		{1, 2, NW},
		{2, 1, NW},
	}

	for index, data := range allTestData {
		results := isAdjacentCellInBounds(&shapeExtr, data.startCol, data.startRow, data.direction)
		if ! results {
			t.Errorf("Expected true result at index %d for %+v", index, data)
		}
	}
}

func TestIsSameColorAdjacentFalse(t *testing.T) {
	textGrid := `
g.r.b
g.r.b
b.g.r`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	results := isSameColorAdjacent(&shapeExtr, 1, 1, grid[1][1], NE, E, S, SW, W, NW)
	if results {
		t.Error("Expected false result.")
	}
}

func TestIsSameColorAdjacentFalseDueToAlreadyUsed(t *testing.T) {
	textGrid := `
g.r.b
g.r.b
b.g.r`

	gridUsed := [][]bool{
		{true, true, false},
		{true, false, false},
		{false, false, false},
	}

	grid := initGrid(textGrid, gridUsed, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	results := isSameColorAdjacent(&shapeExtr, 1, 1, grid[1][1], N)
	if results {
		t.Error("Expected false result.")
	}
}

func TestIsSameColorAdjacentTrue(t *testing.T) {
	textGrid := `
g.r.b
r.r.b
b.g.r`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	results := isSameColorAdjacent(&shapeExtr, 1, 1, grid[1][1], E, W)
	if ! results {
		t.Error("Expected true result.")
		return
	}

	textGrid = `
g.g.b
b.r.r
b.g.b`

	grid = initGrid(textGrid, nil, t)

	shapeExtr = ShapeExtractor{}
	shapeExtr.Init(grid)

	results = isSameColorAdjacent(&shapeExtr, 1, 1, grid[1][1], E, W)
	if ! results {
		t.Error("Expected true result.")
		return
	}
}

// Given a start position on a Red cell with non-Red cells to its south and southwest
// the function should return false
func TestIsStartPositionValidFalseBadCellsToSouthAndSouthWest(t *testing.T) {

	textGrid := `
g.g.b
g.r.b
b.g.r`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	startCol := 0
	startRow := 1

	results := isStartPositionValid(&shapeExtr, startCol, startRow, grid[startCol][startRow])
	if results {
		t.Error("Expected a false return value but got true")
	}
}

// Given a start position on a Red cell with used Red cells to its south and southwest
// the function should return false
func TestIsStartPositionValidFalseAlreadyUsedCellsToSouthAndSouthWest(t *testing.T) {

	textGrid := `
g.g.b
g.r.b
r.r.r`

	gridUsed := [][]bool{
		{true, true, false},
		{true, false, false},
		{true, true, false},
	}

	grid := initGrid(textGrid, gridUsed, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	startCol := 0
	startRow := 1

	results := isStartPositionValid(&shapeExtr, startCol, startRow, grid[startCol][startRow])
	if results {
		t.Error("Expected a false return value but got true")
	}
}


// Given a start position on a Red cell with non-Red cells to its south and southeast
// the function should return false
func TestIsStartPositionValidFalseBadCellsToSouthAndSouthEast(t *testing.T) {

	textGrid := `
g.g.b
g.r.b
r.g.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	startCol := 0
	startRow := 1

	results := isStartPositionValid(&shapeExtr, startCol, startRow, grid[startCol][startRow])
	if results {
		t.Error("Expected a false return value but got true")
	}
}

// Given a start position on a the eastern edge and a cell with a different color to the south
// the function should return false
func TestIsStartPositionValidFalseOnEasternEdge(t *testing.T) {

	textGrid := `
r.r.b
r.g.r
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	startCol := 2
	startRow := 1

	results := isStartPositionValid(&shapeExtr, startCol, startRow, grid[startCol][startRow])
	if results {
		t.Error("Expected a false return value but got true")
	}
}

// Given a start position on a the southern edge
// the function should return false
func TestIsStartPositionValidFalseOnSouthernEdge(t *testing.T) {

	textGrid := `
r.r.b
r.g.r
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	startCol := 1
	startRow := 2

	results := isStartPositionValid(&shapeExtr, startCol, startRow, grid[startCol][startRow])
	if results {
		t.Error("Expected a false return value but got true")
	}
}

// Given a start position with the same color to the south and east
// the function should return true
func TestIsStartPositionValidTrueGoodCellsToTheSouthAndEast(t *testing.T) {

	textGrid := `
r.r.b
r.g.r
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	startCol := 0
	startRow := 0

	results := isStartPositionValid(&shapeExtr, startCol, startRow, grid[startCol][startRow])
	if ! results {
		t.Error("Expected a true return value but got false")
	}
}

// Given a start position with the same color to the south and southeast
// the function should return true
func TestIsStartPositionValidTrueGoodCellsToTheSouthAndSouthEast(t *testing.T) {

	textGrid := `
g.r.b
g.r.r
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	startCol := 1
	startRow := 0

	results := isStartPositionValid(&shapeExtr, startCol, startRow, grid[startCol][startRow])
	if ! results {
		t.Error("Expected a true return value but got false")
	}
}

// Given a start position with the same color to the east and southeast
// the function should return true
func TestIsStartPositionValidTrueGoodCellsToTheEastAndSouthEast(t *testing.T) {

	textGrid := `
g.g.b
g.r.r
b.g.r`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	startCol := 1
	startRow := 1

	results := isStartPositionValid(&shapeExtr, startCol, startRow, grid[startCol][startRow])
	if ! results {
		t.Error("Expected a true return value but got false")
	}
}

// Given a start position with the same color to the south and southwest
// the function should return true
func TestIsStartPositionValidTrueGoodCellsToTheSouthAndSouthWest(t *testing.T) {

	textGrid := `
g.g.b
g.g.r
b.r.r`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	startCol := 2
	startRow := 1

	results := isStartPositionValid(&shapeExtr, startCol, startRow, grid[startCol][startRow])
	if ! results {
		t.Error("Expected a true return value but got false")
	}
}

func TestFindRowOfLowerCellInStartingColumnSameAsStartRowBecauseDifferentColor(t *testing.T) {

	textGrid := `
b.b.b
r.r.r
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	startCol := 0
	startRow := 1

	results := findRowOfLowerCellInStartingColumn(
		&shapeExtr,
		startCol,
		startRow,
		grid[startCol][startRow].Color,
	)
	expected := startRow

	if results != expected {
		t.Errorf("Bad row index. Expected %d, but got %d", expected, results)
	}
}

func TestFindRowOfLowerCellInStartingColumnRightBelowStartRowAtBottomOfGrid(t *testing.T) {

	textGrid := `
b.b.b
r.r.r
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	startCol := 1
	startRow := 1

	results := findRowOfLowerCellInStartingColumn(
		&shapeExtr,
		startCol,
		startRow,
		grid[startCol][startRow].Color,
	)
	expected := 2

	if results != expected {
		t.Errorf("Bad row index. Expected %d, but got %d", expected, results)
	}
}

func TestFindRowOfLowerCellInStartingColumnFartherDownTheColumnButNotAtBottom(t *testing.T) {

	textGrid := `
b.r.r
b.r.r
b.r.r
b.r.g
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	startCol := 1
	startRow := 0

	results := findRowOfLowerCellInStartingColumn(
		&shapeExtr,
		startCol,
		startRow,
		grid[startCol][startRow].Color,
	)
	expected := 3

	if results != expected {
		t.Errorf("Bad row index. Expected %d, but got %d", expected, results)
	}
}

func TestFindRowOfLowerCellInStartingColumnFartherDownTheColumnAtBottom(t *testing.T) {

	textGrid := `
b.r.r
b.r.r
b.r.r
r.r.g
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 1
	startRow := 0

	results := findRowOfLowerCellInStartingColumn(
		&shapeExtr,
		startCol,
		startRow,
		grid[startCol][startRow].Color,
	)
	expected := 4

	if results != expected {
		t.Errorf("Bad row index. Expected %d, but got %d", expected, results)
	}
}

func TestIsSubColumnOneColorTrueAdjacentRows(t *testing.T) {
	textGrid := `
b.r.r
r.r.r
b.b.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	results := isSubColumnOneColor(1, 0, 1, grid)
	if ! results {
		t.Errorf("Expected a true return value but got false.")
	}

}

func TestIsSubColumnOneColorTrueSeparatedRows(t *testing.T) {
	textGrid := `
b.r.r
r.r.r
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	results := isSubColumnOneColor(1, 0, 2, grid)
	if ! results {
		t.Errorf("Expected a true return value but got false.")
	}

}
func TestIsSubColumnOneColorFalse(t *testing.T) {
	textGrid := `
b.r.r
r.b.r
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	results := isSubColumnOneColor(1, 0, 2, grid)
	if results {
		t.Errorf("Expected a false return value but got true.")
	}

}
func TestFindNewStartRowGoingUpToTop(t *testing.T) {
	textGrid := `
b.r.r
b.r.r
r.r.r
r.r.r
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	results := findUpperRowForNextColumn(&shapeExtr,1, 2)
	expected := 0

	if results != expected {
		t.Errorf("Bad row index. Expected %d, but got %d.", expected, results)
	}
}

func TestFindNewStartRowGoingUpOne(t *testing.T) {
	textGrid := `
b.r.b
b.r.b
r.r.r
r.r.r
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	results := findUpperRowForNextColumn(&shapeExtr,1, 2)
	expected := 1

	if results != expected {
		t.Errorf("Bad row index. Expected %d, but got %d.", expected, results)
	}
}

func TestFindNewStartRowSameAsStart(t *testing.T) {
	textGrid := `
b.r.b
b.r.b
r.r.b
r.r.r
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	results := findUpperRowForNextColumn(&shapeExtr,1, 2)
	expected := 2

	if results != expected {
		t.Errorf("Bad row index. Expected %d, but got %d.", expected, results)
	}
}

func TestFindNewStartRowGoingDownAtBottom(t *testing.T) {
	textGrid := `
b.r.b
b.r.b
r.r.b
r.r.b
b.r.r`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	results := findUpperRowForNextColumn(&shapeExtr,1, 2)
	expected := 3

	if results != expected {
		t.Errorf("Bad row index. Expected %d, but got %d.", expected, results)
	}
}

func TestFindNewStartRowGoingDownOne(t *testing.T) {
	textGrid := `
b.r.b
b.r.b
r.r.b
r.r.r
b.r.r`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)

	results := findUpperRowForNextColumn(&shapeExtr,1, 1)
	expected := 2

	if results != expected {
		t.Errorf("Bad row index. Expected %d, but got %d.", expected, results)
	}
}
func TestGetUpperRowOfNextColumnWhenItsLowerThanTheStartCell__SlightlyLower(t *testing.T) {

	textGrid := `
b.r.b
b.r.r
b.r.r
r.r.g
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 0
	results, err := getUpperRowOfNextColumnWhenItsLowerThanTheStartCell(
		&shapeExtr,
		startCol,
		startRow,
		4,
		Red(),
	)
    if err != nil {
    	t.Errorf("Got unexpected error: %s", err)
    	return
	}

	expected := 1
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetUpperRowOfNextColumnWhenItsLowerThanTheStartCell__AtBottom(t *testing.T) {

	textGrid := `
b.r.b
b.r.g
b.r.g
r.r.g
b.r.r`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 0
	results, err := getUpperRowOfNextColumnWhenItsLowerThanTheStartCell(
		&shapeExtr,
		startCol,
		startRow,
		4,
		Red(),
	)
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 4
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetUpperRowOfNextColumnWhenItsLowerThanTheStartCell__AtLowerRow(t *testing.T) {

	textGrid := `
b.r.b
b.r.g
b.r.g
r.r.r
b.b.r`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 0
	results, err := getUpperRowOfNextColumnWhenItsLowerThanTheStartCell(&shapeExtr,
		startCol,
		startRow,
		4,
		Red(),
	)
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 3
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetUpperRowOfNextColumnWhenItsLowerThanTheStartCell__None(t *testing.T) {

	textGrid := `
b.r.b
b.r.b
b.r.b
r.r.g
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 0
	_, err := getUpperRowOfNextColumnWhenItsLowerThanTheStartCell(
		&shapeExtr,
		startCol,
		startRow,
		4,
		Red(),
	)
	if err == nil {
		t.Errorf("Expected error for no valid cells but didn't get one.")
		return
	}

}

func TestGetUpperRowOfNextColumnWhenItsLowerThanTheStartCell__PastLowestRow(t *testing.T) {

	textGrid := `
b.r.b
b.r.b
b.r.b
r.r.g
b.b.r`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 0
	_, err := getUpperRowOfNextColumnWhenItsLowerThanTheStartCell(
		&shapeExtr,
		startCol,
		startRow,
		3,
		Red(),
	)
	if err == nil {
		t.Errorf("Expected error for no valid cells but didn't get one.")
		return
	}

}

func TestGetUpperRowOfColumnToRight__None(t *testing.T) {

	textGrid := `
b.b.r
b.r.b
b.r.b
r.r.g
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 1
	_, err := getUpperRowOfColumnToRight(&shapeExtr, startCol, startRow, 4, Red())
	if err == nil {
		t.Errorf("Expected error for no valid cells but didn't get one.")
		return
	}

}

func TestGetUpperRowOfColumnToRight__TopRow1(t *testing.T) {

	textGrid := `
b.r.r
b.r.r
b.r.b
r.r.g
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 0
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToRight(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 0
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetUpperRowOfColumnToRight__TopRow2(t *testing.T) {

	textGrid := `
b.r.r
b.r.b
b.r.b
r.r.g
b.r.g`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 0
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToRight(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 0
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetUpperRowOfColumnToRight__AboveAndAtRight(t *testing.T) {

	textGrid := `
b.b.r
b.b.r
b.r.r
r.r.g
b.r.g`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 2
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToRight(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 1
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetUpperRowOfColumnToRight__AboveAtTop(t *testing.T) {

	textGrid := `
b.b.r.r
b.b.r.r
b.r.r.r
r.r.g.g
b.r.g.g`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 2
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToRight(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 0
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetUpperRowOfColumnToRight__AboveButNotAtTop__WithRedToRight(t *testing.T) {

	textGrid := `
b.b.r.b
b.b.r.b
b.b.r.r
r.r.r.r
b.r.g.g`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 3
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToRight(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 1
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetUpperRowOfColumnToRight__AboveButNotAtTop__OneRedToRight(t *testing.T) {

	textGrid := `
b.b.r.b
b.b.r.b
b.b.r.r
r.r.r.g
b.r.g.g`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 3
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToRight(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 1
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetUpperRowOfColumnToRight__AboveButNotAtTop__NoRedToRight(t *testing.T) {

	textGrid := `
b.b.r.b
b.b.r.b
b.b.r.b
r.r.r.g
b.r.g.g`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	startRow := 3
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToRight(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 2
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetLowerRowOfColumnToRight__SameAsUpperRow(t *testing.T) {

	textGrid := `
b.r.b
b.r.g
b.r.r
r.r.b
b.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	upperRow := 2
	lowestRow := 3
	results := getLowerRowOfColumnToRight(&shapeExtr, startCol, upperRow, lowestRow, Red())

	expected := 2
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetLowerRowOfColumnToRight__SameAsLowerRow(t *testing.T) {

	textGrid := `
b.b.r
b.r.r
b.r.r
r.r.r
b.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	upperRow := 0
	lowestRow := 3
	results := getLowerRowOfColumnToRight(&shapeExtr, startCol, upperRow, lowestRow, Red())

	expected := 3
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetLowerRowOfColumnToRight__AtRight__AvoidStalactite(t *testing.T) {
	textGrid := `
b.b.b
b.r.b
b.r.r
r.b.r
b.b.r`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	upperRow := 2
	lowestRow := 2
	results := getLowerRowOfColumnToRight(&shapeExtr, startCol, upperRow, lowestRow, Red())

	expected := 3
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetLowerRowOfColumnToRight__NotAtRight__ToBottom(t *testing.T) {

	textGrid := `
b.b.b.b
b.r.b.b
b.r.r.r
r.b.r.r
b.b.r.r`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	upperRow := 2
	lowestRow := 2
	results := getLowerRowOfColumnToRight(&shapeExtr, startCol, upperRow, lowestRow, Red())

	expected := 4
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetLowerRowOfColumnToRight__NotAtRight__NotToBottom(t *testing.T) {

	textGrid := `
b.b.b.b
b.r.r.b
b.r.r.r
r.b.r.b
b.b.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	upperRow := 1
	lowestRow := 2
	results := getLowerRowOfColumnToRight(&shapeExtr, startCol, upperRow, lowestRow, Red())

	expected := 3
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetLowerRowOfColumnToRight__NotAtRight__AvoidStalactite(t *testing.T) {

	textGrid := `
b.b.b.b
b.r.r.b
b.r.r.r
r.b.r.b
b.b.r.b`
		grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	upperRow := 1
	lowestRow := 2
	results := getLowerRowOfColumnToRight(&shapeExtr, startCol, upperRow, lowestRow, Red())

	expected := 3
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetShapeColumnsToRight__NotAtRight(t *testing.T) {

	textGrid := `
b.b.b.b
b.r.r.b
b.r.r.r
r.b.r.b
b.b.r.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 1
	upperRow := 1
	lowestRow := 2
	shape := Shape{
		References: map[int][2]int{},
		Color: Red().Color,
	}
	results := getShapeColumnsToRight(&shapeExtr, startCol, upperRow, lowestRow, shape)

	expected := Shape{
		References: map[int][2]int{
			2: {1, 3},
			3: {2, 2},
		},
	}

	errMsg := compareShapeReferences(results, expected)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}
}
func TestGetShapeColumnsToRight__None(t *testing.T) {

	textGrid := `
b.b.b.b
b.r.b.b
b.r.b.r
r.b.b.b
b.b.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 1
	upperRow := 1
	lowestRow := 2
	shape := Shape{
		Color: Red().Color,
	}
	results := getShapeColumnsToRight(&shapeExtr, startCol, upperRow, lowestRow, shape)

	expected := Shape{
		References: map[int][2]int{},
	}

	errMsg := compareShapeReferences(results, expected)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}
}
func TestGetUpperRowOfColumnToLeft__None(t *testing.T) {

	textGrid := `
r.b.b
b.r.b
b.r.b
b.r.r
b.r.g`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 0
	startRow := 1
	_, err := getUpperRowOfColumnToLeft(&shapeExtr, startCol, startRow, 4, Red())
	if err == nil {
		t.Errorf("Expected error for no valid cells but didn't get one.")
		return
	}

}

func TestGetUpperRowOfColumnToLeft__TopRow1(t *testing.T) {

	textGrid := `
r.r.b
r.r.b
b.r.b
g.r.r
g.r.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 0
	startRow := 0
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToLeft(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 0
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetUpperRowOfColumnToLeft__TopRow2(t *testing.T) {

	textGrid := `
r.r.b
b.r.b
b.r.b
g.r.r
g.r.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 0
	startRow := 0
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToLeft(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 0
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetUpperRowOfColumnToLeft__AboveAndAtLeft(t *testing.T) {

	textGrid := `
r.b.b
r.b.b
r.r.b
g.r.r
g.r.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 0
	startRow := 2
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToLeft(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 1
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetUpperRowOfColumnToLeft__AboveAtTop(t *testing.T) {

	textGrid := `
r.r.b.b
r.r.b.b
r.r.r.b
g.g.r.r
g.g.r.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 1
	startRow := 2
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToLeft(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 0
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetUpperRowOfColumnToLeft__AboveButNotAtTop__WithRedToLeft(t *testing.T) {

	textGrid := `
b.r.b.b
b.r.b.b
r.r.b.b
r.r.r.r
g.g.r.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 1
	startRow := 3
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToLeft(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 1
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetUpperRowOfColumnToLeft__AboveButNotAtTop__OneRedToLeft(t *testing.T) {

	textGrid := `
b.r.b.b
b.r.b.b
r.r.b.b
g.r.r.r
g.g.r.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 1
	startRow := 3
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToLeft(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 1
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetUpperRowOfColumnToLeft__AboveButNotAtTop__NoRedToLeft(t *testing.T) {

	textGrid := `
b.r.b.b
b.r.b.b
b.r.b.b
g.r.r.r
g.g.r.g`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 1
	startRow := 3
	lowestRow := len(grid[0])
	results, err := getUpperRowOfColumnToLeft(&shapeExtr, startCol, startRow, lowestRow, Red())
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		return
	}

	expected := 2
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetLowerRowOfColumnToLeft__SameAsUpperRow(t *testing.T) {

	textGrid := `
b.r.b
g.r.b
r.r.b
b.r.r
b.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 0
	upperRow := 2
	lowestRow := 3
	results := getLowerRowOfColumnToLeft(&shapeExtr, startCol, upperRow, lowestRow, Red())

	expected := 2
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetLowerRowOfColumnToLeft__SameAsLowerRow(t *testing.T) {

	textGrid := `
r.b.b
r.r.b
r.r.b
r.r.r
b.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 0
	upperRow := 0
	lowestRow := 3
	results := getLowerRowOfColumnToLeft(&shapeExtr, startCol, upperRow, lowestRow, Red())

	expected := 3
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetLowerRowOfColumnToLeft__AtRight__AvoidStalactite(t *testing.T) {

	textGrid := `
b.b.b
b.r.b
r.r.b
r.b.r
r.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 0
	upperRow := 2
	lowestRow := 2
	results := getLowerRowOfColumnToLeft(&shapeExtr, startCol, upperRow, lowestRow, Red())

	expected := 3
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}

func TestGetLowerRowOfColumnToLeft__NotAtRight__ToBottom(t *testing.T) {

	textGrid := `
b.b.b.b
b.b.r.b
r.r.r.b
r.r.b.r
r.r.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 1
	upperRow := 2
	lowestRow := 2
	results := getLowerRowOfColumnToLeft(&shapeExtr, startCol, upperRow, lowestRow, Red())

	expected := 4
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetLowerRowOfColumnToLeft__NotAtRight__NotToBottom(t *testing.T) {

	textGrid := `
b.b.b.b
b.r.r.b
r.r.r.b
b.r.b.r
b.b.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 1
	upperRow := 1
	lowestRow := 2
	results := getLowerRowOfColumnToLeft(&shapeExtr, startCol, upperRow, lowestRow, Red())

	expected := 3
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetLowerRowOfColumnToLeft__NotAtLeft__AvoidStalactite(t *testing.T) {

	textGrid := `
b.b.b.b
b.r.r.b
r.r.r.b
b.r.b.r
b.r.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 1
	upperRow := 1
	lowestRow := 2
	results := getLowerRowOfColumnToLeft(&shapeExtr, startCol, upperRow, lowestRow, Red())

	expected := 3
	if results != expected {
		t.Errorf("Bad row. Expected %d.  Got %d", expected, results)
	}
}
func TestGetShapeColumnsToLeft__NotAtLeft(t *testing.T) {

	textGrid := `
b.b.b.b
b.r.r.b
r.r.r.b
b.r.b.r
b.r.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	upperRow := 1
	lowestRow := 2
	shape := Shape{
		References: map[int][2]int{},
		Color: Red().Color,
	}
	results := getShapeColumnsToLeft(&shapeExtr, startCol, upperRow, lowestRow, shape)

	expected := Shape{
		References: map[int][2]int{
			0: {2, 2},
			1: {1, 3},
		},
	}

	errMsg := compareShapeReferences(results, expected)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}
}
func TestGetShapeColumnsToLeft__None(t *testing.T) {

	textGrid := `
b.b.b.b
b.b.r.b
r.b.r.b
b.b.b.r
b.b.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	upperRow := 1
	lowestRow := 2
	shape := Shape{
		Color: Red().Color,
	}
	results := getShapeColumnsToLeft(&shapeExtr, startCol, upperRow, lowestRow, shape)

	expected := Shape{References: map[int][2]int{}}

	errMsg := compareShapeReferences(results, expected)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}
}

func TestGetShapeStartingAtCellReference__None(t *testing.T) {

	textGrid := `
b.b.b.b
b.b.r.b
r.b.r.b
b.b.b.r
b.b.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	upperRow := 1

	results := getShapeStartingAtCellReference(&shapeExtr, startCol, upperRow)

	expected := Shape{References: map[int][2]int{}}

	errMsg := compareShapeReferences(results, expected)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}
}

func TestGetShapeStartingAtCellReference__Ladder(t *testing.T) {

	textGrid := `
b.b.b.b.b
b.b.r.g.g
r.b.r.r.g
b.r.r.b.r
b.b.r.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	upperRow := 1

	results := getShapeStartingAtCellReference(&shapeExtr, startCol, upperRow)

	expected := Shape{
		References: map[int][2]int{
			1: {3, 3},
			2: {1, 4},
			3: {2, 2},
		},
	}

	errMsg := compareShapeReferences(results, expected)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}

	usedTrue := GridCell{AlreadyUsed: true}
	usedFalse := GridCell{AlreadyUsed: false}

	resultsGrid := shapeExtr.grid
	expectedGrid := Grid{
		// Column 0
		{usedFalse, usedFalse, usedFalse, usedFalse, usedFalse},
	 	// Column 1
		{usedFalse, usedFalse, usedFalse, usedTrue, usedFalse},
		 // Column 2
		{usedFalse, usedTrue, usedTrue, usedTrue, usedTrue},
		// Column 3
		{usedFalse, usedFalse, usedTrue, usedFalse, usedFalse},
		// Column 4
		{usedFalse, usedFalse, usedFalse, usedFalse, usedFalse},
	}

	errMsg = compareGridsAlreadyUsed(resultsGrid, expectedGrid)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}

}


func TestGetShapeStartingAtCellReference__WithHole(t *testing.T) {

	textGrid := `
b.b.b.b.b
b.b.r.g.g
r.b.r.r.g
b.r.r.b.r
b.b.r.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 2
	upperRow := 1

	results := getShapeStartingAtCellReference(&shapeExtr, startCol, upperRow)

	expected := Shape{
		References: map[int][2]int{
			1: {3, 3},
			2: {1, 4},
			3: {2, 2},
		},
	}

	errMsg := compareShapeReferences(results, expected)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}

	usedTrue := GridCell{AlreadyUsed: true}
	usedFalse := GridCell{AlreadyUsed: false}

	resultsGrid := shapeExtr.grid
	expectedGrid := Grid{
		// Column 0
		{usedFalse, usedFalse, usedFalse, usedFalse, usedFalse},
		// Column 1
		{usedFalse, usedFalse, usedFalse, usedTrue, usedFalse},
		// Column 2
		{usedFalse, usedTrue, usedTrue, usedTrue, usedTrue},
		// Column 3
		{usedFalse, usedFalse, usedTrue, usedFalse, usedFalse},
		// Column 4
		{usedFalse, usedFalse, usedFalse, usedFalse, usedFalse},
	}

	errMsg = compareGridsAlreadyUsed(resultsGrid, expectedGrid)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}
}



func TestGetShapeStartingAtCellReference__Triangle(t *testing.T) {

	textGrid := `
b.r.r.r.r
b.b.r.r.r
b.b.b.r.r
b.b.b.b.r
b.b.b.b.b`
	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 1
	upperRow := 0

	results := getShapeStartingAtCellReference(&shapeExtr, startCol, upperRow)

	expected := Shape{
		References: map[int][2]int{
			1: {0, 0},
			2: {0, 1},
			3: {0, 2},
			4: {0, 3},
		},
	}

	errMsg := compareShapeReferences(results, expected)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}

	usedTrue := GridCell{AlreadyUsed: true}
	usedFalse := GridCell{AlreadyUsed: false}

	resultsGrid := shapeExtr.grid
	expectedGrid := Grid{
		// Column 0
		{usedFalse, usedFalse, usedFalse, usedFalse, usedFalse},
		// Column 1
		{usedTrue, usedFalse, usedFalse, usedFalse, usedFalse},
		// Column 2
		{usedTrue, usedTrue, usedFalse, usedFalse, usedFalse},
		// Column 3
		{usedTrue, usedTrue, usedTrue, usedFalse, usedFalse},
		// Column 4
		{usedTrue, usedTrue, usedTrue, usedTrue, usedFalse},
	}

	errMsg = compareGridsAlreadyUsed(resultsGrid, expectedGrid)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}

}


func TestGetShapeStartingAtCellReference__BigComplicatedOne(t *testing.T) {
/*
0.1.2.3.4.5.6.7.8.9.0.1.2.3
*/
	textGrid := `
b.b.b.r.b.b.b.b.b.g.g.g.g.g
b.b.r.r.b.r.g.r.g.g.g.r.g.g
b.b.r.r.b.r.r.r.g.r.r.r.g.g
r.r.r.r.r.r.r.r.r.r.r.r.r.r
r.b.r.r.r.r.r.r.r.g.g.r.g.g
b.b.r.r.r.r.r.r.r.g.g.r.r.g
b.b.r.b.r.r.b.r.r.r.r.r.b.r
b.b.r.g.b.r.b.r.r.r.r.r.g.g
r.b.r.r.b.r.r.b.r.b.b.b.b.g
b.b.b.b.r.r.b.r.r.b.r.r.b.r
b.b.b.b.b.r.b.b.r.b.b.r.b.b`

	grid := initGrid(textGrid, nil, t)

	shapeExtr := ShapeExtractor{}
	shapeExtr.Init(grid)
	startCol := 3
	startRow := 0

	results := getShapeStartingAtCellReference(&shapeExtr, startCol, startRow)

	expected := Shape{
		References: map[int][2]int{
			1: {3, 3},  // Column 0 cells don't get included because of this choke point (single cell)
			2: {1, 6},
			3: {0, 5},
			4: {3, 6},
			5: {1, 7},
			6: {2, 5},
			7: {1, 7},
			8: {3, 8},
			9: {2, 3},   // Lower cells don't get included because of the other colors blocking them
			10: {2, 3},
			11: {1, 4},
			12: {3, 3},
		},
	}

	errMsg := compareShapeReferences(results, expected)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}

	uTrue := GridCell{AlreadyUsed: true}
	uFalse := GridCell{AlreadyUsed: false}

	resultsGrid := shapeExtr.grid
	expectedGrid := Grid{
		// Column 0
		{uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse},
		// Column 1
		{uFalse, uFalse, uFalse, uTrue, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse},
		// Column 2
		{uFalse, uTrue, uTrue, uTrue, uTrue, uTrue, uTrue, uFalse, uFalse, uFalse, uFalse},
		// Column 3
		{uTrue, uTrue, uTrue, uTrue, uTrue, uTrue, uFalse, uFalse, uFalse, uFalse, uFalse},
		// Column 4
		{uFalse, uFalse, uFalse, uTrue, uTrue, uTrue, uTrue, uFalse, uFalse, uFalse, uFalse},
		// Column 5
		{uFalse, uTrue, uTrue, uTrue, uTrue, uTrue, uTrue, uTrue, uFalse, uFalse, uFalse},
		// Column 6
		{uFalse, uFalse, uTrue, uTrue, uTrue, uTrue, uFalse, uFalse, uFalse, uFalse, uFalse},
		// Column 7
		{uFalse, uTrue, uTrue, uTrue, uTrue, uTrue, uTrue, uTrue, uFalse, uFalse, uFalse},
		// Column 8
		{uFalse, uFalse, uFalse, uTrue, uTrue, uTrue, uTrue, uTrue, uTrue, uFalse, uFalse},
		// Column 9
		{uFalse, uFalse, uTrue, uTrue, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse},
		// Column 10
		{uFalse, uFalse, uTrue, uTrue, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse},
		// Column 11
		{uFalse, uTrue, uTrue, uTrue, uTrue, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse},
		// Column 12
		{uFalse, uFalse, uFalse, uTrue, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse},
		// Column 13
		{uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse, uFalse},
	}

	errMsg = compareGridsAlreadyUsed(resultsGrid, expectedGrid)
	if errMsg != "" {
		t.Errorf(errMsg)
		return
	}
}


func TestGetPolygonFromShape__SideColumnsTaller(t *testing.T) {

	textGrid := `
r.b.b.b.b
r.r.r.b.r
r.r.r.r.r
r.r.b.r.r
r.b.b.b.r`

	redColor := Red().Color

	grid := initGrid(textGrid, nil, t)

	shape := getShapeFromGrid(grid, redColor)

	resultPolygon, _ := getPolygonFromShape(shape)
	results := resultPolygon.Points

	expected := [][2]int{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 1},
		{4, 4},
		{3, 3},
		{2, 2},
		{1, 3},
		{0, 4},
	}

	errMsg := compareOutlinePoints(results, expected)
	if errMsg != "" {
		t.Errorf("%s\n Expected: \n %v\n  But got: \n%v", errMsg, expected, results)
		return
	}

}


func TestGetPolygonFromShape__ExtremeVariationsInColumns(t *testing.T) {

	textGrid := `
b.r.r.b.b
b.r.r.b.r
r.r.r.r.r
b.r.r.r.b
b.r.r.r.b`

	redColor := Red().Color

	grid := initGrid(textGrid, nil, t)

	shape := getShapeFromGrid(grid, redColor)

	resultPolygon, _ := getPolygonFromShape(shape)
	results := resultPolygon.Points

	expected := [][2]int{
		{0, 2},
		{1, 1},
		{1, 0},
		{2, 0},
		{2, 1},
		{3, 2},
		{4, 1},
		{4, 2},
		{3, 3},
		{3, 4},
		{2, 4},
		{1, 4},
		{1, 3},
	}

	errMsg := compareOutlinePoints(results, expected)
	if errMsg != "" {
		t.Errorf(errMsg)
		//t.Errorf("%s\n Expected: \n %v\n  But got: \n%v", errMsg, expected, results)
		return
	}

}


func TestGetPolygonFromShape__WithHole(t *testing.T) {

	textGrid := `
b.r.b.b.b
b.r.r.b.b
b.r.r.r.b
b.r.b.r.r
b.b.r.r.r`

	redColor := Red().Color

	grid := initGrid(textGrid, nil, t)

	shape := getShapeFromGrid(grid, redColor)

	resultPolygon, _ := getPolygonFromShape(shape)
	results := resultPolygon.Points

	expected := [][2]int{
		{1, 0},
		{2, 1},
		{3, 2},
		{4, 3},
		{4, 4},
		{3, 4},
		{3, 3},
		{2, 2},
		{1, 3},
	}

	errMsg := compareOutlinePoints(results, expected)
	if errMsg != "" {
		t.Errorf(errMsg)
		//t.Errorf("%s\n Expected: \n %v\n  But got: \n%v", errMsg, expected, results)
		return
	}

}


func TestGetPolygonFromShape__Triangle(t *testing.T) {

	textGrid := `
r.r.r.r.b
b.r.r.r.b
b.b.r.r.b
b.b.b.r.b
b.b.b.b.b`

	redColor := Red().Color

	grid := initGrid(textGrid, nil, t)

	shape := getShapeFromGrid(grid, redColor)

	resultPolygon, _ := getPolygonFromShape(shape)
	results := resultPolygon.Points

	expected := [][2]int{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
		{3, 3},
		{2, 2},
		{1, 1},
	}

	errMsg := compareOutlinePoints(results, expected)
	if errMsg != "" {
		t.Errorf(errMsg)
		//t.Errorf("%s\n Expected: \n %v\n  But got: \n%v", errMsg, expected, results)
		return
	}
}