# gopixels2svg

Create simple svg xml from a 2-dimensional grid of pixel colors.

## Overview ##
This app takes a 2-dimensional "grid" of pixel colors (4-length arrays of uint8's for the RGBA values of pixels). 
The "grid" will need to be made up of a nested slice ([column][row]) of those color values.

The app reads through the grid and determines which svg polygons and lines are needed to approximately reproduce 
that image in svg format. (Note that at this time, only the Red, Green and Blue values of the original colors are used.)
It then writes the corresponding svg xml to a file.

## Example ##
The **examples/main.go** file has simple examples of how to use the package to 
 - convert a grid of colors or
 - read in a *.png file and convert it
  
You can also copy *.png files into the `examples` folder and include their names as command line arguments, e.g.
  
   `.../examples> go run main.go my-logo.png my-art.png`
  
That will create corresponding *.html files.

## Binaries ##
The dist folder includes binaries for different operating systems.
They allow you to convert *.png files to svg by including the png file name(s) as command line arguments.