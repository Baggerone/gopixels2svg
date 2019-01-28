package pixels2svg

import (
	"fmt"
	"testing"
)

func compareBoolSlices(results, expected []bool) string {
	resultsCount := len(results)
	expectedCount := len(expected)

	if resultsCount != expectedCount {
		return fmt.Sprintf("\n For length: Expected %d, but got %d", expectedCount, resultsCount)
	}

	for index, nextResult := range results {
		if nextResult != expected[index] {
			return fmt.Sprintf("\n  Expected %v, \n   but got %v", expected, results)
		}
	}
	return ""
}

func compareBoolGrids(results, expected [][]bool) string {
	resultsCount := len(results)
	expectedCount := len(expected)

	if resultsCount != expectedCount {
		return fmt.Sprintf("\n For length: Expected %d, but got %d", expectedCount, resultsCount)
	}

	for index, nextResult := range results {
		err := compareBoolSlices(nextResult, expected[index])
		if err != "" {
			return fmt.Sprintf("\n  For column %d. %s", index, err)
		}
	}
	return ""
}

func comparePolygonPointsSlices(results, expected [][][2]int) string {
	resultsCount := len(results)
	expectedCount := len(expected)

	if resultsCount != expectedCount {
		return fmt.Sprintf("\n For length of all Polygons: Expected %d, but got %d", expectedCount, resultsCount)
	}

	for index, nextResult := range results {
		err := compareOutlinePoints(nextResult, expected[index])
		if err != "" {
			return fmt.Sprintf("\n For Polygon: %d, %s", index, err)
		}
	}
	return ""
}

func comparePolygons(results, expected []Polygon) string {
	resultsCount := len(results)
	expectedCount := len(expected)

	if resultsCount != expectedCount {
		return fmt.Sprintf("\n For length of all Polygons: Expected %d, but got %d", expectedCount, resultsCount)
	}

	for index, nextPolygon := range results {
		err := compareOutlinePoints(nextPolygon.Points, expected[index].Points)
		if err != "" {
			return fmt.Sprintf("\n For Polygon: %d, %s", index, err)
		}
		expectedColor := expected[index].ColorRGBA
		resultsColor := nextPolygon.ColorRGBA
		if resultsColor != expectedColor {
			return fmt.Sprintf(
				"\n For ColorRGBA: Expected %d, but got %d",
				expectedColor,
				resultsColor,
			)
		}
	}
	return ""
}

func compareLines(results, expected []Line) string {
	resultsCount := len(results)
	expectedCount := len(expected)

	if resultsCount != expectedCount {
		return fmt.Sprintf(
			"\n For length of all Lines: Expected %d, but got %d",
			expectedCount,
			resultsCount,
		)
	}

	for index, nextResult := range results {
		expectedLine := expected[index]
		expectedColor := expectedLine.ColorRGBA
		resultsColor := nextResult.ColorRGBA
		if resultsColor != expectedColor {
			return fmt.Sprintf(
				"\n Line %d. For ColorRGBA: Expected %d, but got %d",
				index,
				expectedColor,
				resultsColor,
			)
		}

		labels := [4]string{"ColX1", "RowY1", "ColX2", "RowY2"}
		resultCoords := [4]int{
			nextResult.ColX1,
			nextResult.RowY1,
			nextResult.ColX2,
			nextResult.RowY2,
		}
		expectedCoords := [4]int{
			expectedLine.ColX1,
			expectedLine.RowY1,
			expectedLine.ColX2,
			expectedLine.RowY2,
		}

		for index, label := range labels {
			if resultCoords[index] != expectedCoords[index] {
				return fmt.Sprintf(
					"\n Line %d. %s: Expected %d, but got %d",
					index,
					label,
					expectedCoords[index],
					resultCoords[index],
				)
			}
		}
	}
	return ""
}

func compareOutlinePoints(results, expected [][2]int) string {
	resultsCount := len(results)
	expectedCount := len(expected)

	if resultsCount != expectedCount {
		return fmt.Sprintf("\n For length: Expected %d, but got %d", expectedCount, resultsCount)
	}

	for index, nextResult := range results {
		if nextResult != expected[index] {
			return fmt.Sprintf("\n  Expected %v, \n   but got %v", expected, results)
		}
	}
	return ""
}

func compareSliceofStrings(results, expected []string) string {
	resultsCount := len(results)
	expectedCount := len(expected)

	if resultsCount != expectedCount {
		return fmt.Sprintf("\n For length: Expected %d, but got %d", expectedCount, resultsCount)
	}

	for index, nextResult := range results {
		if nextResult != expected[index] {
			return fmt.Sprintf("\n  Expected %v, \n   but got %v", expected, results)
		}
	}
	return ""
}

/*
 *  Five columns, four rows - all almost black
 */
func getColorGrid() Grid {
	black := [4]uint8{1, 1, 1, 1}
	grid := Grid{
		{{Color: black}, {Color: black}, {Color: black}, {Color: black}}, // Column 0
		{{Color: black}, {Color: black}, {Color: black}, {Color: black}}, // Column 1
		{{Color: black}, {Color: black}, {Color: black}, {Color: black}}, // Column 2
		{{Color: black}, {Color: black}, {Color: black}, {Color: black}}, // Column 3
		{{Color: black}, {Color: black}, {Color: black}, {Color: black}}, // Column 4
	}
	return grid
}

func getAlreadyUsedFromGrid(grid Grid) [][]bool {
	used := [][]bool{}
	for _, column := range grid {
		newColumn := []bool{}
		for _, cell := range column {
			newColumn = append(newColumn, cell.AlreadyUsed)
		}
		used = append(used, newColumn)
	}
	return used
}

// TODO: This should cause a failure because of the internal reds from 7, 7, but it doesn't
/*
 *    0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16 17
 *  0 r  r  r  r  r  r  y  y  y  y  y  y  g  g  g  g  g  g
 *  1 r  r  r  r  r  r  b  y  y  y  y  y  g  g  g  g  g  g
 *  2 r  r  r  r  r  r  b  b  y  y  y  y  g  g  g  g  g  g
 *  3 r  r  r  r  r  r  b  b  b  y  y  y  g  g  g  g  g  g
 *  4 r  r  r  b  r  r  b  b  b  b  y  y  g  g  r  r  g  g
 *  5 r  r  r  b  r  r  b  b  b  b  b  y  g  g  r  r  g  g
 *  6 r  r  r  b  r  r  b  b  b  b  b  b  g  g  r  r  g  g
 *  7 r  r  r  b  r  r  b  r  r  b  b  b  g  g  r  r  g  g
 *  8 r  r  r  r  r  r  b  r  r  b  b  b  g  g  g  g  g  g
 *  9 r  r  r  r  r  r  b  r  r  b  b  b  g  g  g  g  g  g
 * 10 r  r  r  r  r  r  b  r  r  b  b  b  g  g  g  g  g  g
 * 11 r  r  r  r  r  r  b  b  b  b  b  b  g  g  g  g  g  g
 * 12
 */
func getBigColorGrid() Grid {
	grid := Grid{}

	red := Color{235, 0, 0, 0}
	yellow := Color{235, 235, 0, 0}
	blue := Color{0, 0, 220, 0}
	green := Color{0, 180, 60, 0}

	// Left Rectangle
	for colX := 0; colX < 6; colX++ {
		nextCol := []GridCell{}
		for rowY := 0; rowY < 12; rowY++ {
			nextCol = append(nextCol, GridCell{Color:red})
		}
		grid = append(grid, nextCol)
	}

	// Middle polygons
	for colX := 6; colX < 12; colX++ {
		nextCol := []GridCell{}
		for rowY := 0; rowY < 12; rowY++ {
			if rowY < 2+colX-6 {
				nextCol = append(nextCol, GridCell{Color: yellow})
			} else {
				nextCol = append(nextCol, GridCell{Color: blue})
			}
		}
		grid = append(grid, nextCol)
	}

	// Right Rectangle
	for colX := 12; colX < 18; colX++ {
		nextCol := []GridCell{}
		for rowY := 0; rowY < 12; rowY++ {
			nextCol = append(nextCol, GridCell{Color: green})
		}
		grid = append(grid, nextCol)
	}

	// Add rectangle in middle of Right Rectangle
	for colX := 14; colX < 16; colX++ {
		for rowY := 4; rowY < 8; rowY++ {
			grid[colX][rowY] = GridCell{Color: red}
		}
	}

	// Add a line in middle of Left Rectangle
	for rowY := 4; rowY < 8; rowY++ {
		grid[3][rowY] = GridCell{Color: blue}
	}

	return grid
}


func TestGetLineHorizontalWhole(t *testing.T) {
	var s ShapeExtractor
	s.Init(getColorGrid())

	startCol := 0
	startRow := 0

	results := s.GetLine(startCol, startRow)

	expected := Line{
		ColorRGBA: [4]uint8{1, 1, 1, 1},
		ColX1:     0,
		RowY1:     0,
		ColX2:     s.ColCount - 1,
		RowY2:     0,
	}

	err := compareLines([]Line{results}, []Line{expected})
	if err != "" {
		t.Errorf(err)
	}
}

func TestGetLineHorizontalPartial(t *testing.T) {
	var s ShapeExtractor
	s.Init(getColorGrid())
	s.grid[s.ColCount-1][1] = GridCell{Color: [4]uint8{9, 9, 9, 9}}
	startCol := 1
	startRow := 1

	results := s.GetLine(startCol, startRow)

	expected := Line{
		ColorRGBA: [4]uint8{1, 1, 1, 1},
		ColX1:     1,
		RowY1:     1,
		ColX2:     s.ColCount - 2,
		RowY2:     1,
	}

	err := compareLines([]Line{results}, []Line{expected})
	if err != "" {
		t.Errorf(err)
	}
}

func TestGetLineLastColumnWhole(t *testing.T) {
	var s ShapeExtractor
	s.Init(getColorGrid())

	startCol := s.ColCount - 1
	startRow := 0

	results := s.GetLine(startCol, startRow)

	expected := Line{
		ColorRGBA: [4]uint8{1, 1, 1, 1},
		ColX1:     startCol,
		RowY1:     0,
		ColX2:     startCol,
		RowY2:     s.RowCount - 1,
	}

	err := compareLines([]Line{results}, []Line{expected})
	if err != "" {
		t.Errorf(err)
	}
}

func TestGetLineVerticalPartial(t *testing.T) {
	var s ShapeExtractor
	s.Init(getColorGrid())
	s.grid[1][0].AlreadyUsed =  true
	s.grid[2][0].AlreadyUsed =  true
	s.grid[2][1] = GridCell{Color: Color{2, 2, 2, 2}}
	s.grid[2][2] = GridCell{Color: Color{2, 2, 2, 2}}
	s.grid[1][s.RowCount-1] = GridCell{Color: Color{2, 2, 2, 2}}

	startCol := 1
	startRow := 1

	results := s.GetLine(startCol, startRow)

	expected := Line{
		ColorRGBA: [4]uint8{1, 1, 1, 1},
		ColX1:     startCol,
		RowY1:     startRow,
		ColX2:     startCol,
		RowY2:     s.RowCount - 2,
	}

	err := compareLines([]Line{results}, []Line{expected})
	if err != "" {
		t.Errorf(err)
	}
}

// Line toward Southwest
func TestGetLineAngledPartial(t *testing.T) {
	var s ShapeExtractor
	s.Init(getColorGrid())
	s.grid[3][0].AlreadyUsed = true
	s.grid[4][0].AlreadyUsed = true
	s.grid[4][1] = GridCell{Color: Color{2, 2, 2, 2}}
	s.grid[4][2] = GridCell{Color: Color{2, 2, 2, 2}}
	s.grid[4][3] = GridCell{Color: Color{2, 2, 2, 2}}
	s.grid[3][2] = GridCell{Color: Color{2, 2, 2, 2}}
	s.grid[1][s.RowCount-1] = GridCell{Color: Color{2, 2, 2, 2}}

	startCol := 3
	startRow := 1

	results := s.GetLine(startCol, startRow)

	expected := Line{
		ColorRGBA: Color{1, 1, 1, 1},
		ColX1:     startCol,
		RowY1:     startRow,
		ColX2:     2,
		RowY2:     2,
	}

	err := compareLines([]Line{results}, []Line{expected})
	if err != "" {
		t.Errorf(err)
	}
}

func TestGetLineOneCell(t *testing.T) {
	var s ShapeExtractor
	s.Init(getColorGrid())

	startCol := 0
	startRow := 0

	s.grid[0][1] = GridCell{Color: Color{9, 9, 9, 9}}
	s.grid[1][0].AlreadyUsed = true
	s.grid[1][1].AlreadyUsed = true

	results := s.GetLine(startCol, startRow)

	expected := Line{
		ColorRGBA: [4]uint8{1, 1, 1, 1},
		ColX1:     startCol,
		RowY1:     startRow,
		ColX2:     startCol,
		RowY2:     startRow,
	}

	err := compareLines([]Line{results}, []Line{expected})
	if err != "" {
		t.Errorf(err)
	}
}


func TestGetLeftDirectionFromNorth(t *testing.T) {
	var s ShapeExtractor

	results := s.getLeftDirection(0)
	expected := 6

	if results != expected {
		t.Errorf("Expected %d, but got %d", expected, results)
	}
}

func TestGetLeftDirectionFromEast(t *testing.T) {
	var s ShapeExtractor

	results := s.getLeftDirection(2)
	expected := 0

	if results != expected {
		t.Errorf("Expected %d, but got %d", expected, results)
	}
}

func TestGetRightDirectionFromNorth(t *testing.T) {
	var s ShapeExtractor

	results := s.getAngledRightDirection(0)
	expected := 1

	if results != expected {
		t.Errorf("Expected %d, but got %d", expected, results)
	}
}

func TestGetRightDirectionFromWest(t *testing.T) {
	var s ShapeExtractor

	results := s.getAngledRightDirection(6)
	expected := 7

	if results != expected {
		t.Errorf("Expected %d, but got %d", expected, results)
	}
}

func TestGetRightDirectionFromNorthWest(t *testing.T) {
	var s ShapeExtractor

	results := s.getAngledRightDirection(7)
	expected := 0

	if results != expected {
		t.Errorf("Expected %d, but got %d", expected, results)
	}
}

func TestFindOutlineOverlap(t *testing.T) {
	outlinePoints := [][2]int{
		{1, 1},
		{1, 2},
		{1, 3},
		{2, 3},
		{3, 3},
		{3, 2},
		{2, 2},
		{1, 2},
		{1, 1},
		{2, 1},
	}

	overlap := FindOutlineOverlap(outlinePoints)

	results := overlap[0]
	expected := 1

	if results != expected {
		t.Errorf("Expected %d, but got %d", expected, results)
	}

	results = overlap[1]
	expected = 7

	if results != expected {
		t.Errorf("Expected %d, but got %d", expected, results)
	}
}


/*
 *     0   1   2   3   4
 * 0 | A | A | B | B | B |
 * 1 | A | B | B | B | C |
 * 2 | C | B | B | B | C |
 * 3 | C | C | B | C | C |
 */
func TestProcessAllPolygons(t *testing.T) {
	var s ShapeExtractor
	s.Init(getColorGrid())

	expectedDone := getAlreadyUsedFromGrid(s.grid)

	// Color the B cells
	BColRows := [][]int{
		{},
		{1, 2},
		{0, 1, 2, 3},
		{0, 1, 2},
		{0},
	}

	for colIndex, nextCol := range BColRows {
		for _, nextRow := range nextCol {
			s.grid[colIndex][nextRow] = GridCell{Color: Color{2, 2, 2, 2}}
		}
	}

	// Color the C cells
	CColRows := [][]int{
		{2, 3},
		{3},
		{},
		{3},
		{1, 2, 3},
	}

	for colIndex, nextCol := range CColRows {
		for _, nextRow := range nextCol {
			s.grid[colIndex][nextRow] = GridCell{Color: Color{3, 3, 3, 3}}
		}
	}

	allPolygons := s.ProcessAllPolygons()

	results := allPolygons
	expected := []Polygon{
		{ // A's
			ColorRGBA: Color{1, 1, 1, 1},
			Points:    [][2]int{{0, 0}, {1, 0}, {0, 1}},
		},
		{ // B's
			ColorRGBA: Color{2, 2, 2, 2},
			Points: [][2]int{
				{2, 0},
				// [2]int{3, 0}, // reduced into next one, since in a line
				{4, 0},
				{3, 1}, // gets reduced into next one, since in a line
				{3, 2},
				{2, 3},
				{1, 2},
				{1, 1},
			},
		},
		{ // C1
			ColorRGBA: Color{3, 3, 3, 3},
			Points:    [][2]int{{4, 2}, {4, 3}, {3, 3}},
		},
		{ // C2
			ColorRGBA: Color{3, 3, 3, 3},
			Points:    [][2]int{{0, 2}, {1, 3}, {0, 3}},
		},
	}

	err := comparePolygons(results, expected)
	if err != "" {
		t.Errorf("\nPolygons. %s", err)
		return
	}

	resultsDone := getAlreadyUsedFromGrid(s.grid)
	expectedDone[1] = []bool{false, true, true, false}
	expectedDone[2] = []bool{true, true, true, false}
	expectedDone[3] = []bool{true, true, true, false}

	err = compareBoolGrids(resultsDone, expectedDone)
	if err != "" {
		t.Errorf("Already done. %s", err)
		return
	}
}

/*
 *     0   1   2   3   4
 * 0 | B | A | A | A | A |
 * 1 | A | A | A | A | A |
 * 2 | C | A | A | A | C |
 * 3 | C | C | A | A | C |
 */
func TestGetAllShapes(t *testing.T) {
	var s ShapeExtractor
	s.Init(getColorGrid())

	// Color B cell
	s.grid[0][0] = GridCell{Color: Color{2, 2, 2, 2}}

	// Color C cells
	s.grid[0][2] = GridCell{Color: Color{3, 3, 3, 3}}
	s.grid[0][3] = GridCell{Color: Color{3, 3, 3, 3}}
	s.grid[1][3] = GridCell{Color: Color{3, 3, 3, 3}}
	s.grid[4][2] = GridCell{Color: Color{3, 3, 3, 3}}
	s.grid[4][3] = GridCell{Color: Color{3, 3, 3, 3}}

	allPolygons, allLines := s.GetAllShapes()

	results := allPolygons
	expected := []Polygon{
		{ // A's
			ColorRGBA: [4]uint8{1, 1, 1, 1},
			Points: [][2]int{
				{1, 0},
				{4, 0},
				{4, 1},
				{3, 2},
				{3, 3},
				{2, 3},
				{1, 2},
				{0, 1},
			},
		},
		{ // C1
			ColorRGBA: [4]uint8{3, 3, 3, 3},
			Points:    [][2]int{{0, 2}, {1, 3}, {0, 3}},
		},
	}

	err := comparePolygons(results, expected)
	if err != "" {
		t.Errorf("\nPolygons. %s", err)
		return
	}

	expectedLines := []Line{
		{ // B
			ColorRGBA: [4]uint8{2, 2, 2, 2},
			ColX1:     0,
			RowY1:     0,
			ColX2:     0,
			RowY2:     0,
		},
		{ // C2
			ColorRGBA: [4]uint8{3, 3, 3, 3},
			ColX1:     4,
			RowY1:     2,
			ColX2:     4,
			RowY2:     3,
		},
	}

	err = compareLines(allLines, expectedLines)
	if err != "" {
		t.Errorf(err)
		return
	}
}

/*
 *     0   1   2   3   4
 * 0 | B | A | A | A | A |
 * 1 | A | A | A | A | A |
 * 2 | C | A | A | A | C |
 * 3 | C | C | A | A | C |
 */
func TestGetSVGText(t *testing.T) {
	var s ShapeExtractor
	s.Init(getColorGrid())

	// Color B cell
	s.grid[0][0] = GridCell{Color: Color{2, 2, 222, 2}}

	// Color C cells
	s.grid[0][2] = GridCell{Color: Color{223, 3, 3, 3}}
	s.grid[0][3] = GridCell{Color: Color{223, 3, 3, 3}}
	s.grid[1][3] = GridCell{Color: Color{223, 3, 3, 3}}
	s.grid[4][2] = GridCell{Color: Color{223, 3, 3, 3}}
	s.grid[4][3] = GridCell{Color: Color{223, 3, 3, 3}}

	results := s.GetSVGText()
	expected := `<svg width="5" height="4">
 <g>
  <polygon class="#010101" points="1,0 4,0 4,1 3,2 3,3 2,3 1,2 0,1 " stroke="#010101" fill="#010101" />
  <polygon class="#DF0303" points="0,2 1,3 0,3 " stroke="#DF0303" fill="#DF0303" />
  <line class="#0202DE" x1="0" y1="0" x2="0" y2="0" stroke="#0202DE" fill="#0202DE" />
  <line class="#DF0303" x1="4" y1="2" x2="4" y2="3" stroke="#DF0303" fill="#DF0303" />
 </g>
</svg>`

	if results != expected {
		t.Errorf("\nExpected \n%s, \nbut got \n%s", expected, results)
		return
	}
}

/*
 *  Two rectangles on the sides two polygons
 *
 */
func TestGetSVGTextLarge(t *testing.T) {
	var s ShapeExtractor

	gridColors := getBigColorGrid()
	s.Init(gridColors)

	results := s.GetSVGText()
	expected := `<svg width="18" height="12">
 <g>
  <polygon class="#EB0000" points="0,0 5,0 5,11 0,11 0,1 " stroke="#EB0000" fill="#EB0000" />
  <polygon class="#EBEB00" points="6,0 11,0 11,6 6,1 " stroke="#EBEB00" fill="#EBEB00" />
  <polygon class="#00B43C" points="12,0 17,0 17,11 12,11 12,1 " stroke="#00B43C" fill="#00B43C" />
  <polygon class="#0000DC" points="6,2 11,7 11,11 6,11 6,3 " stroke="#0000DC" fill="#0000DC" />
  <polygon class="#EB0000" points="14,4 15,4 15,7 14,7 14,6 14,5 " stroke="#EB0000" fill="#EB0000" />
  <line class="#0000DC" x1="3" y1="4" x2="3" y2="7" stroke="#0000DC" fill="#0000DC" />
 </g>
</svg>`

	if results != expected {
		t.Errorf("\nExpected \n%s, \nbut got \n%s", expected, results)
		return
	}
}

/*
 *  Manually check the file that has been written
 */
func TestWriteSVGToFile(t *testing.T) {
	var s ShapeExtractor

	gridColors := getBigColorGrid()
	s.Init(gridColors)

	s.WriteSVGToFile("test_svg_output.xml")
}

func TestGetIndexOfLastRepeatDirectionWhole(t *testing.T) {
	outlinePoints := [][2]int{
		{3, 3}, // South
		{3, 4}, // South
		{3, 5}, // South
		{3, 6}, // South
		{3, 7}, // South
	}
	results := getIndexOfLastRepeatDirection(outlinePoints)
	expected := len(outlinePoints) - 1
	if results != expected {
		t.Errorf("\nFor index. Expected %d, but got %d", expected, results)
		return
	}
}

func TestGetIndexOfLastRepeatDirectionPartial(t *testing.T) {
	outlinePoints := [][2]int{
		{3, 3}, // South
		{3, 4}, // South
		{3, 5}, // last South
		{4, 5}, // East
		{4, 6}, // South
	}
	results := getIndexOfLastRepeatDirection(outlinePoints)
	expected := 2
	if results != expected {
		t.Errorf("\nFor index. Expected %d, but got %d", expected, results)
		return
	}
}

func TestGetIndexOfLastRepeatDirectionPatternPairWhole(t *testing.T) {
	outlinePoints := [][2]int{
		{3, 3},
		{4, 3}, // East
		{4, 4}, // South
		{5, 4}, // East
		{5, 5}, // South
		{6, 5}, // East
		{6, 6}, // South
	}

	results := getIndexOfLastRepeatDirectionPattern(outlinePoints)
	expected := len(outlinePoints) - 1
	if results != expected {
		t.Errorf("\nFor index. Expected %d, but got %d", expected, results)
		return
	}
}

func TestGetIndexOfLastRepeatDirectionPatternPairAlmostWhole(t *testing.T) {
	outlinePoints := [][2]int{
		{3, 3},
		{4, 3}, // East
		{4, 4}, // South
		{5, 4}, // East
		{5, 5}, // South
		{6, 5}, // East
		{6, 6}, // South
		{7, 6}, // East
	}

	results := getIndexOfLastRepeatDirectionPattern(outlinePoints)
	expected := 6
	if results != expected {
		t.Errorf("\nFor index. Expected %d, but got %d", expected, results)
		return
	}
}

func TestGetIndexOfLastRepeatDirectionPatternPairPartial(t *testing.T) {
	outlinePoints := [][2]int{
		{3, 3},
		{4, 3}, // East
		{4, 4}, // South
		{5, 4}, // East
		{5, 5}, // South - last of pattern
		{6, 5}, // East
		{7, 5}, // East
	}

	results := getIndexOfLastRepeatDirectionPattern(outlinePoints)
	expected := 4
	if results != expected {
		t.Errorf("\nFor index. Expected %d, but got %d", expected, results)
		return
	}
}

func TestGetIndexOfLastRepeatDirectionPatternWhole(t *testing.T) {
	outlinePoints := [][2]int{
		{3, 3},
		{4, 3}, // East
		{5, 3}, // East
		{6, 3}, // East
		{6, 4}, // South
		{7, 4}, // East
		{8, 4}, // East
		{9, 4}, // East
		{9, 5}, // South
	}

	results := getIndexOfLastRepeatDirectionPattern(outlinePoints)
	expected := len(outlinePoints) - 1
	if results != expected {
		t.Errorf("\nFor index. Expected %d, but got %d", expected, results)
		return
	}
}

func TestGetIndexOfLastRepeatDirectionPatternAlmostWhole(t *testing.T) {
	outlinePoints := [][2]int{
		{3, 3},
		{4, 3},  // East
		{5, 3},  // East
		{6, 3},  // East
		{6, 4},  // South
		{7, 4},  // East
		{8, 4},  // East
		{9, 4},  // East
		{9, 5},  // South - last pattern
		{10, 5}, // East
	}

	results := getIndexOfLastRepeatDirectionPattern(outlinePoints)
	expected := 8
	if results != expected {
		t.Errorf("\nFor index. Expected %d, but got %d", expected, results)
		return
	}
}

func TestGetIndexOfLastRepeatDirectionPatternPartial(t *testing.T) {
	outlinePoints := [][2]int{
		{3, 3},
		{4, 3},  // East
		{5, 3},  // East
		{6, 3},  // East
		{6, 4},  // South
		{7, 4},  // East
		{8, 4},  // East
		{9, 4},  // East
		{9, 5},  // South - last pattern
		{10, 5}, // East
		{11, 5}, // East
		{12, 5}, // East
		{13, 5}, // East
		{14, 5}, // East
	}

	results := getIndexOfLastRepeatDirectionPattern(outlinePoints)
	expected := 8
	if results != expected {
		t.Errorf("\nFor index. Expected %d, but got %d", expected, results)
		return
	}
}

func TestGetIndexOfLastRepeatDirectionPatternShort(t *testing.T) {
	outlinePoints := [][2]int{
		{3, 3},
		{4, 3}, // East
		{5, 3}, // East
		{6, 3}, // East
		{7, 3}, // East
		{7, 4}, // South
		{8, 4}, // East
	}

	results := getIndexOfLastRepeatDirectionPattern(outlinePoints)
	expected := 0
	if results != expected {
		t.Errorf("\nFor index. Expected %d, but got %d", expected, results)
		return
	}
}

func TestGetIndexOfLastRepeatDirectionPatternTiny(t *testing.T) {
	outlinePoints := [][2]int{{3, 3}}

	results := getIndexOfLastRepeatDirectionPattern(outlinePoints)
	expected := 0
	if results != expected {
		t.Errorf("\nFor index. Expected %d, but got %d", expected, results)
		return
	}
}

func TestReducePolygonOutlineHorizontalLine(t *testing.T) {

	outlinePoints := [][2]int{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
		{4, 0},
		{4, 1},
	}

	resultsCount, resultsPoints := ReducePolygonOutline(outlinePoints)
	expectedCount := 1

	if resultsCount != expectedCount {
		t.Errorf("Expected %d, but got %d", expectedCount, resultsCount)
	}

	expectedPoints := [][2]int{
		{0, 0},
		{4, 0},
		{4, 1},
	}

	err := compareOutlinePoints(resultsPoints, expectedPoints)
	if err != "" {
		t.Errorf("Reduced Polygon outline. %s", err)
	}
}

func TestReducePolygonOutlineVerticalLine(t *testing.T) {

	outlinePoints := [][2]int{
		{0, 0},
		{0, 1},
		{0, 2},
		{0, 3},
		{0, 4},
		{1, 4},
	}

	resultsCount, resultsPoints := ReducePolygonOutline(outlinePoints)
	expectedCount := 1

	if resultsCount != expectedCount {
		t.Errorf("Expected %d, but got %d", expectedCount, resultsCount)
		return
	}

	expectedPoints := [][2]int{
		{0, 0},
		{0, 4},
		{1, 4},
	}

	err := compareOutlinePoints(resultsPoints, expectedPoints)
	if err != "" {
		t.Errorf("\nReduced Polygon outline. %s", err)
	}
}

func TestReducePolygonOutlineFallingSlope(t *testing.T) {
	outlinePoints := [][2]int{
		{0, 0},
		{0, 1}, // South 1
		{1, 1}, // East 1
		{1, 2}, // South 1
		{2, 2}, // East 1 ... Reduce preceding three into this one
		{2, 3}, // South 1
		{1, 3}, // West 1
	}

	resultsCount, resultsPoints := ReducePolygonOutline(outlinePoints)
	expectedCount := 1

	if resultsCount != expectedCount {
		t.Errorf("\nExpected %d, but got %d", expectedCount, resultsCount)
		return
	}

	expectedPoints := [][2]int{
		{0, 0},
		{2, 2},
		{2, 3},
		{1, 3},
	}

	err := compareOutlinePoints(resultsPoints, expectedPoints)
	if err != "" {
		t.Errorf("\nReduced Polygon outline. %s", err)
	}
}

func TestReducePolygonOutlineRisingSlope(t *testing.T) {

	outlinePoints := [][2]int{
		{0, 4},
		{0, 3}, // North 1
		{1, 3}, // East 1
		{1, 2}, // North 1
		{2, 2}, // East 1 ... Reduce preceding three into this one
		{2, 1}, // North 1
		{1, 1}, // West 1
	}

	resultsCount, resultsPoints := ReducePolygonOutline(outlinePoints)
	expectedCount := 1

	if resultsCount != expectedCount {
		t.Errorf("\nExpected %d, but got %d", expectedCount, resultsCount)
		return
	}

	expectedPoints := [][2]int{
		{0, 4},
		{2, 2},
		{2, 1},
		{1, 1},
	}

	err := compareOutlinePoints(resultsPoints, expectedPoints)
	if err != "" {
		t.Errorf("Reduced Polygon outline. %s", err)
	}
}

func TestReducePolygonOutlineComplicated(t *testing.T) {

	outlinePoints := [][2]int{
		{0, 0},
		{0, 1}, // South
		{1, 1}, // East
		{1, 2}, // South
		{2, 2}, // East - reduce previous three
		{2, 3}, // South
		{2, 4}, // South
		{2, 5}, // South - reduce previous two
		{3, 5}, // East
		{4, 5}, // East - reduce previous one
		{4, 4}, // North
		{5, 4}, // East
		{5, 5}, // South
		{5, 6}, // South
		{5, 7}, // South - reduce previous two
		{4, 7}, // West
		{3, 7}, // West  - reduce previous one
		{3, 6}, // North
		{2, 6}, // West
		{2, 5}, // North
		{1, 5}, // West - reduce previous three
		{0, 5}, // West
		{0, 4}, // North
		{0, 3}, // North
		{0, 2}, // North - reduce previous two
	}

	_, resultsPoints := ReducePolygonOutline(outlinePoints)

	expectedPoints := [][2]int{
		{0, 0},
		{2, 2}, // East
		{2, 5}, // South
		{4, 5}, // East
		{4, 4}, // North
		{5, 4}, // East
		{5, 7}, // South
		{3, 7}, // West
		{1, 5}, // West
		{0, 5}, // West
		{0, 2}, // North
	}

	err := compareOutlinePoints(resultsPoints, expectedPoints)
	if err != "" {
		t.Errorf("Reduced Polygon (multi pass) outline. %s", err)
	}
}
