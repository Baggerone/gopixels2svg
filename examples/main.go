package main

import (
	pixels2svg "github.com/baggerone/gopixels2svg/pixels2svg"
	// "fmt"
)

// Transpose a grid from row by column to column by row
func transposeGrid(grid [][][4]uint8) [][][4]uint8 {
	newGrid := [][][4]uint8{}

	// initialize new grid
	for range grid[0] {
		newColumn := [][4]uint8{}
		for range grid {
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

	// the colors to be used for the different letters of the "image"
	colors := map[string][4]uint8{
		" ": {0, 0, 150, 0},
		"s": {250, 250, 245, 0},
		"m": {150, 150, 0, 0},
		"h": {220, 50, 0, 0},
	}
	// convert the letters in the image strings (above) to colors of a column by row grid
	return assignColorsToGrid(image, colors)
}

/*
 * In order to see the SVG as an image,
 * copy the file named example_sailboat.xml as an *.html file and then open it in a browser
 */
func main() {
	var s pixels2svg.ShapeExtractor

	s.Init(sailboat())

	s.WriteSVGToFile("example_sailboat.xml")
}
