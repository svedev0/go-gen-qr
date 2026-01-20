package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

const quietZone = 4

func renderPNG(mx *matrix, scale int) *image.RGBA {
	var modules int = mx.size + 2*quietZone
	var imgSize int = modules * scale
	var imgRect image.Rectangle = image.Rect(0, 0, imgSize, imgSize)
	var img *image.RGBA = image.NewRGBA(imgRect)

	var white color.RGBA = color.RGBA{255, 255, 255, 255}
	var black color.RGBA = color.RGBA{0, 0, 0, 255}

	for y := range imgSize {
		for x := range imgSize {
			img.Set(x, y, white)
		}
	}

	for row := 0; row < mx.size; row++ {
		for col := 0; col < mx.size; col++ {
			if !mx.m[row][col].val {
				continue
			}

			x0 := (col + quietZone) * scale
			y0 := (row + quietZone) * scale
			for y := y0; y < y0+scale; y++ {
				for x := x0; x < x0+scale; x++ {
					img.Set(x, y, black)
				}
			}
		}
	}

	return img
}

func exportPNG(img *image.RGBA, outPath string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
