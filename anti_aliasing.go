package main

import (
	"image"
	"image/color"
	"math"
)

func getMaxEdgeValue(img *image.Gray) uint8 {
	highest := uint8(0)
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			clr := img.GrayAt(x, y).Y

			if clr > highest {
				highest = clr
			}

		}
	}

	return highest
}

// AntiAliasing Implement Anti Aliasing In The Image
// The First Version It's Too Blury
func AntiAliasing(realImg *image.RGBA, img *image.Gray) *image.RGBA {
	maxEdgeValue := float64(getMaxEdgeValue(img))

	newImg := image.NewRGBA(img.Bounds())

	convolutionKernel := [][]int{
		{1, 4, 6, 4, 1},
		{4, 16, 24, 16, 4},
		{6, 24, 36, 24, 4},
		{4, 16, 24, 16, 4},
		{1, 4, 6, 4, 1},
	}

	for y := 0; y < realImg.Bounds().Dy(); y++ {
		for x := 0; x < realImg.Bounds().Dx(); x++ {
			clr := img.GrayAt(x, y)

			edgeStrength := float64(clr.Y) / float64(maxEdgeValue)

			if edgeStrength > 1 {
				edgeStrength = 1
			}

			sumPixelRed := 0.0
			sumPixelGreen := 0.0
			sumPIxelBlue := 0.0
			sumPixelAlpha := 0.0
			sumWeight := 0.0

			for kernelX := 0; kernelX < len(convolutionKernel); kernelX++ {
				for kernelY := 0; kernelY < len(convolutionKernel[kernelX]); kernelY++ {

					nx := x - 2 + kernelX
					ny := y - 2 + kernelY

					if nx < 0 || ny < 0 || nx >= realImg.Bounds().Dx() || ny >= realImg.Bounds().Dy() {
						continue
					}
					clr := realImg.RGBAAt(nx, ny)

					appliedKernelRed := float64(clr.R) * float64(convolutionKernel[kernelX][kernelY])
					appliedKernelGreen := float64(clr.G) * float64(convolutionKernel[kernelX][kernelY])
					appliedKernelBlue := float64(clr.B) * float64(convolutionKernel[kernelX][kernelY])
					appliedKernelAlpha := float64(clr.A) * float64(convolutionKernel[kernelX][kernelY])

					// sumPixelGray += appliedKernelGray
					sumPixelRed += appliedKernelRed
					sumPixelGreen += appliedKernelGreen
					sumPIxelBlue += appliedKernelBlue
					sumPixelAlpha += appliedKernelAlpha

					sumWeight += float64(convolutionKernel[kernelX][kernelY])
				}
			}
			original := realImg.RGBAAt(x, y)

			red := sumPixelRed / sumWeight
			redFinalPixel := (1-edgeStrength)*float64(original.R) + float64(edgeStrength)*red

			green := sumPixelGreen / sumWeight
			greenFinalPixel := (1-edgeStrength)*float64(original.G) + float64(edgeStrength)*green

			blue := sumPIxelBlue / sumWeight
			blueFinalPixel := (1-edgeStrength)*float64(original.B) + float64(edgeStrength)*blue

			alpha := sumPixelAlpha / sumWeight
			alphaFinalPixel := (1-edgeStrength)*float64(original.A) + float64(edgeStrength)*alpha

			newClr := color.RGBA{
				R: ClampUint8(redFinalPixel, 0.0, 255.0),
				G: ClampUint8(greenFinalPixel, 0.0, 255.0),
				B: ClampUint8(blueFinalPixel, 0.0, 255.0),
				A: ClampUint8(alphaFinalPixel, 0.0, 255.0),
			}

			newImg.Set(x, y, newClr)
		}
	}

	return newImg
}

func ClampUint8(x float64, min float64, max float64) uint8 {
	if x > max {
		return 255
	}

	if x < min {
		return 0
	}

	return uint8(x)
}

// EdgeAwareAntiAliasing Second Version Anti Aliasing
// This Not Cause Too Much Blur 
func EdgeAwareAntiAliasing(realImg *image.RGBA, edgeImg *image.Gray, gx, gy [][]float64) *image.RGBA {
	w, h := realImg.Bounds().Dx(), realImg.Bounds().Dy()
	maxMagnitude := getMaxEdgeValue(edgeImg)

	kernelRadius := 2
	// Weight Distribution Blur
	// Why is 1D instead of 2D. I dont want to make the edge too much blur 
	weights := []int{1, 4, 6, 4, 1} // And This too I can implement Gaussian Function to look similar as FXAA
	newImg := image.NewRGBA(realImg.Bounds())

	for y := 0; y < realImg.Bounds().Dy(); y++ {
		for x := 0; x < realImg.Bounds().Dx(); x++ {
			magClr := edgeImg.GrayAt(x, y).Y
			edgeStrength := float64(magClr) / float64(maxMagnitude)

			if edgeStrength > 1 {
				edgeStrength = 1
			}

			// Compute Edge Direction
			dx := gx[y][x]
			dy := gy[y][x]
			distanceXAndY := math.Hypot(dx, dy)

			if distanceXAndY == 0 {
				newImg.SetRGBA(x, y, realImg.RGBAAt(x, y))
				continue
			}

			nx := -dy / distanceXAndY
			ny := dx / distanceXAndY

			sumR, sumG, sumB, sumA := 0.0, 0.0, 0.0, 0.0
			weight := 0.0

			// To Store Bend Value. Instead Find the Surronding Pixel Like Normal Blur
			// This Just find the the Edge Pixel in 1D
			// K = 0 is the center pixel value
			for k := -kernelRadius; k <= kernelRadius; k++ {
				sx := float64(x) + nx*float64(k)
				sy := float64(y) + ny*float64(k)

				// Some Suggest I can do Bilinear Interpolation By finding middle value

				ix := int(math.Round(sx))
				iy := int(math.Round(sy))

				if ix < 0 || ix >= w || iy < 0 || iy >= h {
					continue
				}

				weightIndex := k + kernelRadius

				clr := realImg.RGBAAt(ix, iy)
				sumR += float64(clr.R) * float64(weights[weightIndex])
				sumG += float64(clr.G) * float64(weights[weightIndex])
				sumB += float64(clr.B) * float64(weights[weightIndex])
				sumA += float64(clr.A) * float64(weights[weightIndex])

				weight += float64(weights[weightIndex])
			}
			original := realImg.RGBAAt(x, y)

			red := sumR / weight
			redFinalPixel := (1-edgeStrength)*float64(original.R) + float64(edgeStrength)*red

			green := sumG / weight
			greenFinalPixel := (1-edgeStrength)*float64(original.G) + float64(edgeStrength)*green

			blue := sumB / weight
			blueFinalPixel := (1-edgeStrength)*float64(original.B) + float64(edgeStrength)*blue

			alpha := sumA / weight
			alphaFinalPixel := (1-edgeStrength)*float64(original.A) + float64(edgeStrength)*alpha

			newClr := color.RGBA{
				R: ClampUint8(redFinalPixel, 0.0, 255.0),
				G: ClampUint8(greenFinalPixel, 0.0, 255.0),
				B: ClampUint8(blueFinalPixel, 0.0, 255.0),
				A: ClampUint8(alphaFinalPixel, 0.0, 255.0),
			}
			newImg.SetRGBA(x, y, newClr)
		}
	}

	return newImg
}
