package main

import (
	"image"

	"github.com/disintegration/imaging"
)

type Image struct {
	Image image.Image
	Type string
}

func ImageNew(image image.Image, imgType string) *Image {
	return &Image{
		Image: image,
		Type: imgType,
	}
}

func (image *Image) changeBrightness(intensity float64)  {
	image.Image = imaging.AdjustBrightness(image.Image, intensity)
}

func (image *Image) blur(intensity float64)  {
	image.Image = imaging.Blur(image.Image, intensity)
}

func (image *Image) sharpen(intensity float64)  {
	image.Image = imaging.Sharpen(image.Image, intensity);
}

func (image *Image) changeContrast(intensity float64) {
	image.Image = imaging.AdjustContrast(image.Image, intensity);
}
