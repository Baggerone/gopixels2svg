# gopixels2svg

Create simple svg xml from a 2-dimensional grid of pixel colors.

## Overview ##
This app takes a 2-dimensional "grid" of pixel colors (4-length arrays of uint8's for the RGBA values of pixels). 
The "grid" will need to be made up of a nested slice ([column][row]) of those color values.

The app reads through the grid and determines which svg polygons and lines are needed to approximately reproduce 
that image in svg format. (Note that at this time, only the Red, Green and Blue values of the original colors are used.)
It then writes the corresponding svg xml to a file.

## Example ##
The **examples/main.go** file has a simple example of how to use the package to convert a grid of colors that represents 
a sailboat.
  
When you run it, it creates an svg file named example_sailboat.xml.  If you copy that file and rename it as a *.html file,
you can open it in a browser to see the image.

