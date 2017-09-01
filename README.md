# gopixels2svg

Create simple svg xml from a 2-dimensional grid of pixel colors.

For a quick idea of how to use this app, look at **TestWriteSVGToFile** in the **pixels2svg_test.go** file.
 
## Overview ##
This app takes a 2-dimensional "grid" of pixel colors (4-length arrays of uint8's for the RGBA values of pixels). 
The "grid" will need to be made up of a nested slice ([column][row]) of those color values.

The app reads through the grid and determines which svg polygons and lines are needed to approximately reproduce 
that image in svg format. (Note that at this time, only the Red, Green and Blue values of the original colors are used.)
It then writes the corresponding svg xml to a file.



