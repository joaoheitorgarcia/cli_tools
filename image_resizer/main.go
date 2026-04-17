package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	xdraw "golang.org/x/image/draw"
)

func main() {
	//flag definition
	filePath := flag.String("f", "", "Path to image file to pixelize")
	finalWidth := flag.Int("w", 0, "Final width in pixels.\nIf Width is specified without Height (ommiting flag or passing 0) the output image will assume original aspect ratio")
	finalHeight := flag.Int("h", 0, "Final height in pixels.\nIf Height is specified without Width (ommiting flag or passing 0) the output image will assume original aspect ratio")
	flag.Parse()

	//flag validation
	if *finalWidth == 0 && *finalHeight == 0 {
		fmt.Print("Error: At least one value needed for Width or Height\n")
		return
	}
	if *finalWidth < 0 || *finalHeight < 0 {
		fmt.Print("Error: Negative values not supported for Width/Height\n")
		return
	}
	_, err := os.Stat(*filePath)
	if err != nil {
		fmt.Print("Error: File not Found\n")
		return
	}

	//load
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Print("Error: File could not be opened\n")
		return
	}
	defer file.Close()

	inputImage, _, err := image.Decode(file)
	if err != nil {
		fmt.Print("Error: File not a supported image format\n")
		return
	}

	ratio := (float64)(inputImage.Bounds().Max.Y) / (float64)(inputImage.Bounds().Max.X)
	outputWidth := *finalWidth
	outputHeight := *finalHeight

	if outputWidth == 0 {
		outputWidth = int(float64(*finalHeight) * ratio)
	} else if outputHeight == 0 {
		outputHeight = int(float64(*finalWidth) * ratio)
	}

	outputImage := image.NewRGBA(image.Rect(0, 0, outputWidth, outputHeight))
	xdraw.NearestNeighbor.Scale(outputImage, outputImage.Bounds(), inputImage, inputImage.Bounds(), draw.Src, nil)

	originalFileName := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
	newFile, err := os.Create(originalFileName + "_resized.png")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	err = png.Encode(newFile, outputImage)
	if err != nil {
		panic(err)
	}
}
