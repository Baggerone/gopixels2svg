package main

import (
	pixels2svg "github.com/baggerone/gopixels2svg/pixels2svg"
	// "fmt"
)

func getBaseGrid(columnCount, rowCount int, color [4]uint8) [][][4]uint8 {
	grid := [][][4]uint8{}

	for col := 0; col < columnCount; col++ {
		nextCol := [][4]uint8{}
		for row := 0; row < rowCount; row++ {
			nextCol = append(nextCol, color)
		}
		grid = append(grid, nextCol)
	}

	return grid
}

func transposeGrid(grid [][][4]uint8) [][][4]uint8 {
	newGrid := [][][4]uint8{}

	// initialize new grid
	for _ = range grid[0] {
		newColumn := [][4]uint8{}
		for _ = range grid {
			newColumn = append(newColumn, [4]uint8{})
		}
		newGrid = append(newGrid, newColumn)
	}

	// populate new grid
	for rowY, row := range grid {
		for colX, cell := range row {
			newGrid[colX][rowY] = cell
		}
	}
	return newGrid
}

func assignColorsToGrid(image []string, colors map[string][4]uint8) [][][4]uint8 {
	grid := [][][4]uint8{}

	for _, row := range image {
		newRow := [][4]uint8{}
		for _, colorCode := range row {
			newRow = append(newRow, colors[string(colorCode)])
		}
		grid = append(grid, newRow)
	}

	return transposeGrid(grid)
}

func sailboat() [][][4]uint8 {
	image := []string{
		"                    ",
		"           m        ",
		"          sm        ",
		"         ssms       ",
		"        sssmss      ",
		"       ssssmss      ",
		"      sssssmsss     ",
		"     ssssssmsss     ",
		"    sssssssmssss    ",
		"   ssssssssmssss    ",
		"           m        ",
		"  hhhhhhhhhhhhhhhhh ",
		"  hhhhhhhhhhhhhhh   ",
		"   hhhhhhhhhhhhh    ",
		"                    ",
		"                    ",
		"                    ",
	}

	colors := map[string][4]uint8{
		" ": [4]uint8{0, 0, 150, 0},
		"s": [4]uint8{250, 250, 245, 0},
		"m": [4]uint8{150, 150, 0, 0},
		"h": [4]uint8{220, 50, 0, 0},
	}

	return assignColorsToGrid(image, colors)
}

/*
 *  Manually check the file that has been written
 */
func main() {
	var s pixels2svg.ShapeExtractor

	s.Init(sailboat())

	s.WriteSVGToFile("example_sailboat.xml")
}
