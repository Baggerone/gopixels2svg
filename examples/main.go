package main

import (
	// "fmt"
	pixels2svg "github.com/baggerone/gopixels2svg/pixels2svg"
	"image"
	_ "image/png" // needed for reading a PNG file, even though it's not explicitly used
	"os"
	"strings"
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

func ReadPNGPixels(filePath string) ([][][4]uint8, error) {

	var infile *os.File
	var err error
	var src image.Image

	// fmt.Println("\nReading \n", filePath)
	if infile, err = os.Open(filePath); err != nil {
		return nil, err
	}

	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	if src, _, err = image.Decode(infile); err != nil {
		return nil, err
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

	return colorGrid, nil
}

func addError(errors *[]string, summary string, err error) {
	*errors = append(
		*errors,
		strings.Join([]string{summary, err.Error()}, " "),
	)
}

/*
 * In order to see the SVG as an image,
 *   open the *.html files in a browser
 * Also, you can copy *.png files into this folder
 *   and include their names as command line arguments
 *     e.g. go run main.go my-logo.png my-art.png
 *   That will create corresponding *.html files.
 */
func main() {
	errors := []string{}
	var s pixels2svg.ShapeExtractor
	var colorGrid [][][4]uint8
	var err error

	s.Init(sailboat())
	s.WriteSVGToFile("example_sailboat.html")

	if colorGrid, err = ReadPNGPixels("test1.png"); err == nil {
		s.Init(colorGrid)
		s.WriteSVGToFile("example_test1.html")
	}

	if colorGrid, err = ReadPNGPixels("test2.png"); err == nil {
		s.Init(colorGrid)
		s.WriteSVGToFile("example_test2.html")
	}

	args := os.Args[1:]
	if len(args) <= 0 {
		println("\n To create from other png files, just add their names as command line arguments.")
	}

	for _, nextInput := range args {
		if strings.HasSuffix(nextInput, ".png") {
			colorGrid, err := ReadPNGPixels(nextInput)
			if err != nil {
				addError(&errors, strings.Join([]string{" Error: ", nextInput, "  ... "}, ""), err)
				continue
			}
			s.Init(colorGrid)
			newName := strings.TrimSuffix(nextInput, ".png") + ".html"
			s.WriteSVGToFile(newName)
		}
	}

	if len(errors) > 0 {
		println("\nRan into error(s) ...")
		for _, nextErr := range errors {
			println(nextErr)
		}
	}
}
