package main

import (
	"image"
	"image/color"
	"math"
)

func allocNewArray[V int | float64 | uint8](x0, y0 int) [][]V {
	arr := make([][]V, y0)
	for y := 0; y < y0; y++ {
		arr[y] = make([]V, x0)
	}

	return arr
}

func EdgeDetection(img image.Gray) (*image.Gray, [][]float64, [][]float64) {
	Gx := allocNewArray[float64](img.Bounds().Dx(), img.Bounds().Dy())
	Gy := allocNewArray[float64](img.Bounds().Dx(), img.Bounds().Dy())

	newPaddedImage := AddEdgePaddingExtenstion(&img, 1, 1)

	horintalKernelConvluation := [3][3]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	verticalKernelConvluation := [3][3]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	for y := newPaddedImage.Bounds().Min.Y + 1; y < newPaddedImage.Bounds().Max.Y-1; y++ {
		for x := newPaddedImage.Bounds().Min.X + 1; x < newPaddedImage.Bounds().Max.X-1; x++ {
			Gx[y-1][x-1] = EdgeCalculation(newPaddedImage, x, y, horintalKernelConvluation)
			Gy[y-1][x-1] = EdgeCalculation(newPaddedImage, x, y, verticalKernelConvluation)
		}
	}

	imgGray := image.NewGray(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	maxMagnitude := 0.0

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			gx := Gx[y][x]
			gy := Gy[y][x]

			magnitude := math.Hypot(gx, gy)

			if magnitude > maxMagnitude {
				maxMagnitude = magnitude
			}
		}
	}

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			gx := Gx[y][x]
			gy := Gy[y][x]

			magnitude := math.Hypot(gx, gy)
			norm := magnitude

			// fmt.Printf("%0.f ", magnitude)

			if magnitude > 0 {
				norm = (magnitude / maxMagnitude) * 255.0
			}

			gray := color.Gray{Y: uint8(norm)}
			// gray := color.Gray{Y: uint8(magnitude)}

			imgGray.SetGray(x, y, gray)
		}
	}

	return imgGray, Gx, Gy
}

func EdgeCalculation(img *image.Gray, x, y int, kernelConvolution [3][3]int) float64 {
	pixelTopLeft := img.GrayAt(x-1, y-1).Y
	pixelTopCenter := img.GrayAt(x, y-1).Y
	pixelTopRight := img.GrayAt(x+1, y-1).Y

	pixelMiddleLeft := img.GrayAt(x-1, y).Y
	pixelMiddleCenter := img.GrayAt(x, y).Y
	pixelMiddleRight := img.GrayAt(x+1, y).Y

	pixelBottomLeft := img.GrayAt(x-1, y+1).Y
	pixelBottomCenter := img.GrayAt(x, y+1).Y
	pixelBottomRight := img.GrayAt(x+1, y+1).Y

	topLeft := int(pixelTopLeft) * kernelConvolution[0][0]
	topCenter := int(pixelTopCenter) * kernelConvolution[0][1]
	topRight := int(pixelTopRight) * kernelConvolution[0][2]

	middleLeft := int(pixelMiddleLeft) * kernelConvolution[1][0]
	middleCenter := int(pixelMiddleCenter) * kernelConvolution[1][1]
	middleRight := int(pixelMiddleRight) * kernelConvolution[1][2]

	bottomLeft := int(pixelBottomLeft) * kernelConvolution[2][0]
	bottomCenter := int(pixelBottomCenter) * kernelConvolution[2][1]
	bottomRight := int(pixelBottomRight) * kernelConvolution[2][2]

	sumPix := float64((topLeft + topCenter + topRight + int(middleLeft) + int(middleCenter) + int(middleRight) + int(bottomLeft) + int(bottomCenter) + int(bottomRight)))
	return sumPix
}
