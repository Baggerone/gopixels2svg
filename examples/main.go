package main

import (
	// "fmt"
	pixels2svg "github.com/baggerone/gopixels2svg/pixels2svg"
	"image"
	_ "image/png" // needed for reading a PNG file, even though it's not explicitly used
	"os"
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
		"           m        ",
		"           m        ",
		"          sm        ",
		"         ssms       ",
		"        sssmss      ",
		"       ssssmss      ",
		"      sssssmsss     ",
		"    sssssssmsss     ",
		"  sssssssssmssss    ",
		"sssssssssssmssss    ",
		"           m        ",
		"  hhhhhhhhhhhhhhhhh ",
		"  hhhhhhhhhhhhhhhh  ",
		"   hhhhhhhhhhhhhh   ",
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

// See https://jimdoescode.github.io/2015/05/22/manipulating-colors-in-go.html
func convertToUint8(rgbTone uint32) uint8 {
	return uint8(rgbTone / 0x101)
}

func ReadPNGPixels(filePath string) [][][4]uint8 {
	// fmt.Println("\nReading \n", filePath)
	infile, err := os.Open(filePath)
	if err != nil {
		// replace this with good error handling
		panic(err)
	}

	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, err := image.Decode(infile)
	if err != nil {
		// replace this with good error handling
		panic(err)
	}

	colorGrid := [][][4]uint8{}
	bounds := src.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	for x := 0; x < width; x++ {
		newCol := [][4]uint8{}
		for y := 0; y < height; y++ {
			red, green, blue, _ := src.At(x, y).RGBA()
			red8 := convertToUint8(red)
			green8 := convertToUint8(green)
			blue8 := convertToUint8(blue)

			colorRGBA := [4]uint8{red8, green8, blue8, 255}
			newCol = append(newCol, colorRGBA)
		}
		colorGrid = append(colorGrid, newCol)
	}

	return colorGrid
}

/*
 * In order to see the SVG as an image,
 * open the *.html files in a browser
 */
func main() {
	var s pixels2svg.ShapeExtractor

	s.Init(sailboat())
	s.WriteSVGToFile("example_sailboat.html")

	s.Init(ReadPNGPixels("test1.png"))
	s.WriteSVGToFile("example_test1.html")
}
