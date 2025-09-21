package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"strconv"

)

func toRGBA(img image.Image) *image.RGBA {
	newImage := image.NewRGBA(img.Bounds())

	draw.Draw(newImage, img.Bounds(), img, img.Bounds().Min, draw.Src)
	
	return newImage
}


func toGray(img image.Image) *image.Gray {
	newImage := image.NewGray(img.Bounds())

	draw.Draw(newImage, img.Bounds(), img, img.Bounds().Min, draw.Src)
	
	return newImage
}

func drawTriangle(img *image.RGBA) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x <= y { // Simple lower-left triangle
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}
}

// Draw a diagonal pattern
func drawDiagonal(img *image.RGBA) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x == y || x == w-y-1 { // Both diagonals
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}
}

func createTestImage() {
	sizes := []int{64, 128}
	for _, size := range sizes {
		// Triangular image
		triImg := image.NewRGBA(image.Rect(0, 0, size, size))
		drawTriangle(triImg)
		OutputImageForDebugResult(triImg, "./img/triangle_"+strconv.Itoa(size)+".png")

		// Diagonal image
		diagImg := image.NewRGBA(image.Rect(0, 0, size, size))
		drawDiagonal(diagImg)
		OutputImageForDebugResult(diagImg, "./img/diagonal_"+strconv.Itoa(size)+".png")
	}
}

func main() {

	// Test Image
	var filePath string

	flag.StringVar(&filePath, "i", "", "Image To Process")

	flag.Parse()

	if filePath == "" {
		log.Fatal("No Image To Play With")
	}

	imgPath := filePath

	reader, err := os.Open(imgPath)

	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()
	
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(reader)

	if err != nil {
		log.Fatal("failed to parse image")
	}

	rgbaImg := toRGBA(img)
	grayImg := toGray(img)

	newImageSet, horizontalEdge, verticalEdge := EdgeDetection(*grayImg)

	OutputImageForDebugResult(newImageSet, "./resImg/edge.jpg")

	// colorImg := AntiAliasing(rgbaImg, newImageSet)
	colorImg := EdgeAwareAntiAliasing(rgbaImg, newImageSet, horizontalEdge, verticalEdge)

	OutputImageForDebugResult(colorImg, "./resImg/anti-aliasing.jpg")

	fmt.Print("Done\n")

}
