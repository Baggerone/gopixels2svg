package pixels2svg

import (
	"bytes"
	"fmt"
	"os"
)

type evaluatorFunc func(int, int, [4]uint8) bool

type Line struct {
	ColorRGBA [4]uint8
	ColX1     int
	RowY1     int
	ColX2     int
	RowY2     int
}

type Polygon struct {
	ColorRGBA [4]uint8
	Points    [][2]int
}

type ShapeExtractor struct {
	badDirection       int
	grid               [][][4]uint8
	alreadyDone        [][]bool
	ColCount           int
	RowCount           int
	neighborEvaluators [8]evaluatorFunc
}

func (s *ShapeExtractor) setNeighborEvaluators() {
	s.neighborEvaluators = [8]evaluatorFunc{
		s.isNorthCellGood,
		s.isNorthEastCellGood,
		s.isEastCellGood,
		s.isSouthEastCellGood,
		s.isSouthCellGood,
		s.isSouthWestCellGood,
		s.isWestCellGood,
		s.isNorthWestCellGood,
	}
}

func (s *ShapeExtractor) isCellDoneOrDifferent(
	nextCol, nextRow int,
	color [4]uint8,
) bool {
	// True if different color or alreadyDone
	return s.grid[nextCol][nextRow] != color || s.alreadyDone[nextCol][nextRow]
}

/*
 *  Given a starting cell, get the furthest cell to its right
 *  (i.e. the column (X) index) that still has its same color.
 */
func (s *ShapeExtractor) getColorRow(
	startCol, rowY int,
	color [4]uint8,
) int {
	for colX := startCol + 1; colX < s.ColCount; colX++ {
		if s.isCellDoneOrDifferent(colX, rowY, color) {
			return colX - 1
		}
	}
	return s.ColCount - 1
}

func (s *ShapeExtractor) cellIsAtLeft(colX int) bool {
	return colX <= 0
}

func (s *ShapeExtractor) cellIsAtTop(rowY int) bool {
	return rowY <= 0
}

func (s *ShapeExtractor) cellIsAtRight(colX int) bool {
	return colX >= s.ColCount-1
}

func (s *ShapeExtractor) cellIsAtBottom(rowY int) bool {
	return rowY >= s.RowCount-1
}

func (s *ShapeExtractor) isNorthCellGood(
	colX, rowY int,
	color [4]uint8,
) bool {
	if s.cellIsAtTop(rowY) {
		return false
	}
	return !s.isCellDoneOrDifferent(colX, rowY-1, color)
}

func (s *ShapeExtractor) isNorthEastCellGood(
	colX, rowY int,
	color [4]uint8,
) bool {
	if s.cellIsAtTop(rowY) || s.cellIsAtRight(colX) {
		return false
	}
	return !s.isCellDoneOrDifferent(colX+1, rowY-1, color)
}

func (s *ShapeExtractor) isEastCellGood(
	colX, rowY int,
	color [4]uint8,
) bool {
	if s.cellIsAtRight(colX) {
		return false
	}
	return !s.isCellDoneOrDifferent(colX+1, rowY, color)
}

func (s *ShapeExtractor) isSouthEastCellGood(
	colX, rowY int,
	color [4]uint8,
) bool {
	if s.cellIsAtRight(colX) || s.cellIsAtBottom(rowY) {
		return false
	}
	return !s.isCellDoneOrDifferent(colX+1, rowY+1, color)
}

func (s *ShapeExtractor) isSouthCellGood(
	colX, rowY int,
	color [4]uint8,
) bool {
	if s.cellIsAtBottom(rowY) {
		return false
	}
	return !s.isCellDoneOrDifferent(colX, rowY+1, color)
}

func (s *ShapeExtractor) isSouthWestCellGood(
	colX, rowY int,
	color [4]uint8,
) bool {
	if s.cellIsAtBottom(rowY) || s.cellIsAtLeft(colX) {
		return false
	}
	return !s.isCellDoneOrDifferent(colX-1, rowY+1, color)
}

func (s *ShapeExtractor) isWestCellGood(
	colX, rowY int,
	color [4]uint8,
) bool {
	if s.cellIsAtLeft(colX) {
		return false
	}
	return !s.isCellDoneOrDifferent(colX-1, rowY, color)
}

func (s *ShapeExtractor) isNorthWestCellGood(
	colX, rowY int,
	color [4]uint8,
) bool {
	if s.cellIsAtLeft(colX) || s.cellIsAtTop(rowY) {
		return false
	}
	return !s.isCellDoneOrDifferent(colX-1, rowY-1, color)
}

/*
 *  Assuming an outline walker is headed a certain direction,
 *  get the direction to its left.
 */
func (s *ShapeExtractor) getLeftDirection(direction int) int {
	if direction <= 1 {
		return 6 + direction
	}
	return direction - 2
}

/*
 *  Assuming an outline walker is headed a certain direction,
 *  get the direction slightly to its right.
 */
func (s *ShapeExtractor) getAngledRightDirection(direction int) int {
	if direction >= 7 {
		return 0
	}
	return direction + 1
}

/*
 * Given a cell and a outline walker's direction, what is the new direction
 * of the first good neighboring cell,
 * starting from the left of the walker and going clockwise.
 */
func (s *ShapeExtractor) directionToGoodNeighboringCell(
	colX, rowY, direction int,
	color [4]uint8,
) int {

	if s.neighborEvaluators[0] == nil {
		s.setNeighborEvaluators()
	}

	newDirection := s.getLeftDirection(direction)
	for index := 0; index < 7; index++ {
		evaluator := s.neighborEvaluators[newDirection]
		if evaluator(colX, rowY, color) {
			return newDirection
		}
		newDirection = s.getAngledRightDirection(newDirection)
	}

	return s.badDirection
}

/*
 * Given a cell and a direction what are the cell cordinates of the neighboring
 * cell in that direction.
 * Assumes that boundary checks are already done, so that the starting cell
 * won't be on a boundary
 */
func (s *ShapeExtractor) getCellInDirection(colX, rowY, nextDirection int) (int, int) {
	switch nextDirection {
	case 0: // North
		return colX, rowY - 1
	case 1: // Northeast
		return colX + 1, rowY - 1
	case 2: // East
		return colX + 1, rowY
	case 3: // Southeast
		return colX + 1, rowY + 1
	case 4: // South
		return colX, rowY + 1
	case 5: // Southwest
		return colX - 1, rowY + 1
	case 6: // West
		return colX - 1, rowY
	case 7: // Northwest
		return colX - 1, rowY - 1
	}

	panic(fmt.Sprintf(
		"Error: direction should be 0, 1, 2, 3, 4, 5, 6, 7.  Was passed: %d",
		nextDirection,
	))
	return 0, 0
}

/*
 * Get the outline of a polygon from a starting cell.
 * Don't use cells of a different color or that have been used previously
 *
 * Trace the outline as if you're a person walking a certain direction
 * clockwise around the edge of the polygon, keeping the outside of
 * the polygon on your left. Stop when you get back to the starting point.
 *
 * Direction codes
 *  0: North
 *  1: Northeast
 *  2: East
 *  3: Southeast
 *  4: South
 *  5: Southwest
 *  6: West
 *  7: Northwest
 */
func (s *ShapeExtractor) OutlinePolygon(
	colX, rowY, direction int,
	color [4]uint8,
) [][2]int {
	outlinePoints := [][2]int{{colX, rowY}}
	if s.alreadyDone[colX][rowY] {
		return nil
	}
	for {
		newDirection := s.directionToGoodNeighboringCell(colX, rowY, direction, color)

		if newDirection >= s.badDirection {
			if len(outlinePoints) <= 2 {
				return nil
			}
		}

		newCol, newRow := s.getCellInDirection(colX, rowY, newDirection)

		if newCol == outlinePoints[0][0] && newRow == outlinePoints[0][1] {
			return outlinePoints
		}

		colX = newCol
		rowY = newRow
		direction = newDirection
		outlinePoints = append(outlinePoints, [2]int{colX, rowY})
	}
}

/*
 * Given a starting cell, get the range of cells to its right and below
 * that have the same color and form a line
 *
 */
func (s *ShapeExtractor) GetLine(startCol, startRow int) Line {
	color := s.grid[startCol][startRow]

	newLine := Line{
		ColorRGBA: color,
		ColX1:     startCol,
		RowY1:     startRow,
	}
	// I want it to start looking to the East
	// That function starts looking to the "left", so tell it I'm facing South
	direction := s.directionToGoodNeighboringCell(startCol, startRow, 4, color)

	if direction >= s.badDirection {
		s.alreadyDone[startCol][startRow] = true
		newLine.ColX2 = startCol
		newLine.RowY2 = startRow // Don't worry if it's just a dot
		return newLine
	}
	prevCol := startCol
	prevRow := startRow

	neighborEvaluator := s.neighborEvaluators[direction]

	for {
		nextCol, nextRow := s.getCellInDirection(prevCol, prevRow, direction)
		if !neighborEvaluator(prevCol, prevRow, color) {
			break
		}

		s.alreadyDone[prevCol][prevRow] = true
		prevCol = nextCol
		prevRow = nextRow
	}

	s.alreadyDone[prevCol][prevRow] = true
	newLine.ColX2 = prevCol
	newLine.RowY2 = prevRow

	return newLine
}

func (s *ShapeExtractor) Init(colorGrid [][][4]uint8) {

	s.badDirection = 8
	s.ColCount = len(colorGrid)
	s.RowCount = len(colorGrid[0])
	s.grid = colorGrid
	s.alreadyDone = [][]bool{}

	// set the alreadyDone grid values
	for colX := 0; colX < s.ColCount; colX++ {
		nextColumn := []bool{}
		for rowY := 0; rowY < s.RowCount; rowY++ {
			nextColumn = append(nextColumn, false)
		}
		s.alreadyDone = append(s.alreadyDone, nextColumn)
	}
}

/*
 * Given a cell on the grid. Get all the polygon outlines that follow
 * from that cell.
 */
func (s *ShapeExtractor) GetPolygonsFromCell(
	colX, rowY, direction int,
	color [4]uint8,
) [][][2]int {

	allPolygons := [][][2]int{}
	outlinePoints := s.OutlinePolygon(colX, rowY, direction, color)
	if len(outlinePoints) < 3 {
		return allPolygons
	}

	cleanedUpPolygons := CleanUpPolygonOutline(
		outlinePoints,
		[][][2]int{},
		0,
	)

	for _, nextPolygon := range cleanedUpPolygons {
		_, reducedPolygon := ReducePolygonOutline(nextPolygon)
		if len(reducedPolygon) > 2 {
			s.markPolygonAlreadyDone(nextPolygon)
			allPolygons = append(allPolygons, reducedPolygon)
		}
	}

	return allPolygons
}

/*
 * Given a polygon with an outline starting at a certain cell, mark
 * all its points (outline and internal) as "alreadyDone".
 *
 * Assumes that the polygon outline is contiguous right-angle points.
 * Assumes first point is the highest row of the outline and the
 *   left-most cell of that row.
 *
 * From each point on the outline of the polygon, marks it and all cells
 * below it as already done until it runs into another cell of the outline
 * or a cell of a different color
 *
 */
func (s *ShapeExtractor) markPolygonAlreadyDone(polygonOutline [][2]int) {
	firstCol, firstRow := split2Int(polygonOutline[0])
	color := s.grid[firstCol][firstRow]

	for _, nextPoint := range polygonOutline {
		nextCol, nextRow := split2Int(nextPoint)
		s.alreadyDone[nextCol][nextRow] = true

		for lowerRow := nextRow + 1; lowerRow < s.RowCount-1; lowerRow++ {

			if color != s.grid[nextCol][lowerRow] {
				continue
			}

			if IsPointIn2IntArray(nextCol, lowerRow, polygonOutline) {
				continue
			}

			s.alreadyDone[nextCol][lowerRow] = true
		}
	}
}

func (s *ShapeExtractor) ProcessAllPolygons() []Polygon {
	startDirection := 2
	allPolygons := []Polygon{}

	// Start at top left and move to the right, then down a row, then right ...
	for rowIndex := 0; rowIndex < s.RowCount; rowIndex++ {
		for colIndex := 0; colIndex < s.ColCount; colIndex++ {
			color := s.grid[colIndex][rowIndex]
			nextPolygons := s.GetPolygonsFromCell(
				colIndex,
				rowIndex,
				startDirection,
				color,
			)
			for _, nextPoly := range nextPolygons {
				newPoly := Polygon{
					ColorRGBA: color,
					Points:    nextPoly,
				}
				allPolygons = append(allPolygons, newPoly)
			}
		}
	}

	return allPolygons
}

/*
 * Goes through the grid, cell by cell, and gets all the lines of the same color.
 *
 */
func (s *ShapeExtractor) ProcessAllLines() []Line {
	allLines := []Line{}

	// Start at top left and move to the right, then down a row, then right ...
	for rowIndex := 0; rowIndex < s.RowCount; rowIndex++ {
		for colIndex := 0; colIndex < s.ColCount; colIndex++ {
			if !s.alreadyDone[colIndex][rowIndex] {
				nextLine := s.GetLine(colIndex, rowIndex)
				allLines = append(allLines, nextLine)
			}
		}
	}

	return allLines
}

func (s *ShapeExtractor) GetAllShapes() ([]Polygon, []Line) {
	s.setNeighborEvaluators()
	allPolygons := s.ProcessAllPolygons()
	allLines := s.ProcessAllLines()

	return allPolygons, allLines
}

func (s *ShapeExtractor) GetSVGText() string {
	allPolygons, allLines := s.GetAllShapes()

	svgWidth := s.ColCount
	svgHeight := s.RowCount

	var svgBuffer bytes.Buffer // Concatenation is more economical with a Buffer
	svgBuffer.WriteString(
		fmt.Sprintf(`<svg width="%d" height="%d">`, svgWidth, svgHeight),
	)
	svgBuffer.WriteString("\n <g>\n")

	for _, next := range allPolygons {
		hexColor := GetHexColor(next.ColorRGBA)
		svgBuffer.WriteString(fmt.Sprintf(`  <polygon class="%s" points="`, hexColor))

		for _, nextPoint := range next.Points {
			svgBuffer.WriteString(fmt.Sprintf("%d,%d ", nextPoint[0], nextPoint[1]))
		}
		svgBuffer.WriteString(fmt.Sprintf(`" stroke="%s" fill="%s" />`, hexColor, hexColor))
		svgBuffer.WriteString("\n")
	}

	for _, next := range allLines {
		hexColor := GetHexColor(next.ColorRGBA)

		svgBuffer.WriteString(fmt.Sprintf(`  <line class="%s" `, hexColor))
		svgBuffer.WriteString(fmt.Sprintf(`x1="%d" y1="%d" `, next.ColX1, next.RowY1))
		svgBuffer.WriteString(fmt.Sprintf(`x2="%d" y2="%d" `, next.ColX2, next.RowY2))
		svgBuffer.WriteString(fmt.Sprintf(`stroke="%s" fill="%s" />`, hexColor, hexColor))
		svgBuffer.WriteString("\n")
	}

	svgBuffer.WriteString(" </g>\n</svg>")

	return svgBuffer.String()
}

func (s *ShapeExtractor) WriteSVGToFile(filePath string) error {
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.WriteString(s.GetSVGText())
	if err == nil {
		println("\nWrote SVG to ", filePath)
	}
	return err
}

func GetHexColor(colorRGBA [4]uint8) string {
	return fmt.Sprintf("#%02X%02X%02X", colorRGBA[0], colorRGBA[1], colorRGBA[2])
}

/*
 * Given an outline of a polygon on a grid, find the pair of indexes
 *  where there is the first overlap and return that pair.
 * Returns empty array if no overlap
 */
func FindOutlineOverlap(outlinePoints [][2]int) [2]int {
	if len(outlinePoints) == 0 {
		return [2]int{}
	}

	for outerIndex, outerPoint := range outlinePoints[1:] {
		outerColX := outerPoint[0]
		outerRowY := outerPoint[1]

		for innerIndex, innerPoint := range outlinePoints[0:outerIndex] {
			innerColX := innerPoint[0]
			innerRowY := innerPoint[1]

			if outerColX == innerColX && outerRowY == innerRowY {
				return [2]int{innerIndex, outerIndex + 1} // Add to outerIndex, since it starts at 1
			}
		}
	}

	return [2]int{}
}

/*
 * Given an outline of a polygon on a grid, purge out the sections that
 * loop back on themselves and return a slice of those sections.
 */
func CleanUpPolygonOutline(
	outlinePoints [][2]int,
	purgedOutlines [][][2]int,
	startIndex int,
) [][][2]int {

	allPolygons := [][][2]int{outlinePoints}
	allPolygons = append(allPolygons, purgedOutlines...)

	if len(outlinePoints) <= startIndex {
		return allPolygons
	}

	overlapPoints := FindOutlineOverlap(outlinePoints)
	if overlapPoints == [2]int{} {
		return allPolygons
	}

	newOutline := [][2]int{} // Create new slice to avoid  modifying original
	newOutline = append(newOutline, outlinePoints[:overlapPoints[0]]...)
	newOutline = append(newOutline, outlinePoints[overlapPoints[1]:]...)

	newPurge := [][2]int{}
	newPurge = append(newPurge, outlinePoints[overlapPoints[0]:overlapPoints[1]]...)

	// Only include sections that have at least two points
	if len(newPurge) > 2 {
		purgedOutlines = append(purgedOutlines, newPurge)
	}

	return CleanUpPolygonOutline(
		newOutline,
		purgedOutlines,
		overlapPoints[0]+1,
	)
}

func getDirection(num1, num2 int, increasingLetter, decreasingLetter string) string {
	if num2 == num1 {
		return ""
	}
	if num2 > num1 {
		return increasingLetter
	}
	return decreasingLetter
}

func getUpDown(row1, row2 int) string {
	return getDirection(row1, row2, "D", "U")
}

func getLeftRight(col1, col2 int) string {
	return getDirection(col1, col2, "R", "L")
}

func getOutlineDirection(colRow1, colRow2 [2]int) string {
	return getLeftRight(colRow1[0], colRow2[0]) + getUpDown(colRow1[1], colRow2[1])
}

func getOutlineDirectionSet(outlinePoints [][2]int) string {
	if len(outlinePoints) <= 1 {
		return ""
	}
	directions := ""

	for index := 0; index < len(outlinePoints)-1; index++ {
		nextDir := getOutlineDirection(
			outlinePoints[index],
			outlinePoints[index+1],
		)
		directions += nextDir
	}
	return directions
}

/*
 * Given a slice of points returns the index of the last point before
 * a different direction is taken (i.e. different than the first direction,
 * e.g. right-right-right-down, returns the index of the point after the
 * third "right").
 */
func getIndexOfLastRepeatDirection(outlinePoints [][2]int) int {

	if len(outlinePoints) <= 1 {
		return 0
	}
	firstDir := getOutlineDirection(
		outlinePoints[0],
		outlinePoints[1],
	)

	for index := 1; index < len(outlinePoints)-1; index++ {
		nextDir := getOutlineDirection(
			outlinePoints[index],
			outlinePoints[index+1],
		)
		if nextDir != firstDir {
			return index
		}
	}

	return len(outlinePoints) - 1
}

/*
 * Given a slice of points, return the index of the last repeat pattern where the
 * pattern is the points that are in the same direction followed by one in a
 * different direction.  For example, right-right-right-down, right-right-right-down, ...
 *
 */
func getIndexOfLastRepeatDirectionPattern(outlinePoints [][2]int) int {
	endOfFirstPattern := getIndexOfLastRepeatDirection(outlinePoints) + 1
	if endOfFirstPattern <= 1 || endOfFirstPattern >= len(outlinePoints)-1 {
		return 0
	}

	directionPattern := getOutlineDirectionSet(outlinePoints[:endOfFirstPattern+1])
	lastIndex := 0
	foundRepeat := false

	for index := endOfFirstPattern; index < len(outlinePoints)-endOfFirstPattern; index += endOfFirstPattern {
		endOfCurrentPattern := index + endOfFirstPattern
		nextPattern := getOutlineDirectionSet(
			outlinePoints[index : endOfCurrentPattern+1])
		if nextPattern != directionPattern {
			if foundRepeat {
				return index
			}
			return 0
		}
		lastIndex = endOfCurrentPattern
		foundRepeat = true
	}
	return lastIndex
}

/*
 * Given the outline of a polygon, find recurring patterns in terms
 * of direction of change (e.g. up, right, down, left) and remove
 * redundant intermediate points.
 *
 */
func ReducePolygonOutline(
	outlinePoints [][2]int,
) (int, [][2]int) {
	if len(outlinePoints) < 3 {
		return 0, outlinePoints
	}

	reductionCount := 0
	index := 0
	newPoints := [][2]int{outlinePoints[0]}

	for index < len(outlinePoints)-4 {

		// First, check for a repeated slope pattern d1-d1-d1-d2, d1-d1-d1-d2
		sameDirectionIndex := getIndexOfLastRepeatDirectionPattern(outlinePoints[index:])
		if sameDirectionIndex > 3 {
			newPoints = append(newPoints, outlinePoints[index+sameDirectionIndex])
			index += sameDirectionIndex
			reductionCount++
			continue
		}

		// Otherwise, get long straight lines (same direction)
		sameDirectionIndex = getIndexOfLastRepeatDirection(outlinePoints[index:])
		if sameDirectionIndex > 1 {
			newPoints = append(newPoints, outlinePoints[index+sameDirectionIndex])
			index += sameDirectionIndex
			reductionCount++
			continue
		}

		// No reduction
		newPoints = append(newPoints, outlinePoints[index+1])
		index++
	}

	// get remaining straight lines (same direction)
	sameDirectionIndex := getIndexOfLastRepeatDirection(outlinePoints[index:])
	if sameDirectionIndex > 1 {
		newPoints = append(newPoints, outlinePoints[index+sameDirectionIndex])
		index += sameDirectionIndex
		reductionCount++
	}

	for endIndex := index + 1; endIndex < len(outlinePoints); endIndex++ {
		newPoints = append(newPoints, outlinePoints[endIndex])
	}

	return reductionCount, newPoints
}

func IsPointIn2IntArray(colX, rowY int, outlinePoints [][2]int) bool {
	for _, nextPoint := range outlinePoints {
		nextCol, nextRow := split2Int(nextPoint)
		if colX == nextCol && rowY == nextRow {
			return true
		}
	}
	return false
}

func split2Int(inArray [2]int) (int, int) {
	return inArray[0], inArray[1]
}
