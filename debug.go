package main

import (
	"image"
	"image/jpeg"
	"log"
	"os"
)

func OutputImageForDebugResult(img image.Image, filePathName string) {
	outFile, err := os.Create(filePathName)
	if err != nil {
		log.Fatal(err)
	}

	defer outFile.Close()

	err = jpeg.Encode(outFile, img, nil)

	if err != nil {
		log.Fatal(err)
	}
}

