package main

import (
	"image"
	"image/color"
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

			red := sumPixelRed/sumWeight
			redFinalPixel := (1 - edgeStrength) * float64(original.R) + float64(edgeStrength) * red

			green := sumPixelGreen/sumWeight
			greenFinalPixel := (1 - edgeStrength) * float64(original.G) + float64(edgeStrength) * green

			blue := sumPIxelBlue/sumWeight
			blueFinalPixel := (1 - edgeStrength) * float64(original.B) + float64(edgeStrength) * blue

			alpha := sumPixelAlpha/sumWeight
			alphaFinalPixel := (1 - edgeStrength) * float64(original.A) + float64(edgeStrength) * alpha

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
