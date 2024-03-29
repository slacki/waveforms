package spectogram

import (
	"bufio"
	"image"
	"image/color"
	"image/png"
	"os"
)

// Image128 implements image.Image interface and represents spectogram image
// TODO: rename to Image (then outside this package it'd be spectogram.Image HOW COOL IS THAT?!?!!)
type Image128 struct {
	set    int
	at     int
	pix    []color.Color
	bounds image.Rectangle
}

// ColorModel returns the Image's color model.
func (img *Image128) ColorModel() color.Model {
	return color.RGBA64Model
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (img *Image128) Bounds() image.Rectangle {
	return img.bounds
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (img *Image128) At(x int, y int) color.Color {
	o := img.offset(x, y)
	if o < 0 {
		return color.RGBA{}
	}
	img.at++
	return img.pix[o]
}

func (img *Image128) offset(x, y int) int {
	p := image.Point{x, y}
	if !p.In(img.bounds) {
		return -1
	}
	stride := img.bounds.Dx()
	my := img.bounds.Min.Y
	mx := img.bounds.Min.X
	ny := y - my
	nx := x - mx
	return ny*stride + nx
}

func (img *Image128) Set(x int, y int, c color.Color) {
	o := img.offset(x, y)
	if o < 0 {
		return
	}
	img.set++
	img.pix[o] = c
}

func (img *Image128) Stats() (int, int) {
	return img.at, img.set
}

// ToPNG writes image to specified path in .png format.
func (img *Image128) ToPNG(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	wr := bufio.NewWriter(f)

	err = png.Encode(wr, img)
	if err != nil {
		return err
	}

	return wr.Flush()
}

func NewImage128(bounds image.Rectangle) *Image128 {
	dx := bounds.Dx()
	dy := bounds.Dy()
	sz := dx * dy
	return &Image128{
		pix:    make([]color.Color, sz),
		bounds: bounds,
	}
}

type SubImage128 struct {
	img    *Image128
	bounds image.Rectangle
}

func (sub *SubImage128) ColorModel() color.Model {
	return sub.ColorModel()
}

func (sub *SubImage128) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Pt(0, 0),
		Max: image.Pt(sub.bounds.Dx(), sub.bounds.Dy()),
	}
}

func (sub *SubImage128) At(x int, y int) color.Color {
	mx := sub.bounds.Min.X
	my := sub.bounds.Min.Y
	return sub.img.At(x+mx, y+my)
}

func (sub *SubImage128) Set(x int, y int, c color.Color) {
	mx := sub.bounds.Min.X
	my := sub.bounds.Min.Y
	sub.img.Set(x+mx, y+my, c)
}

func (img *Image128) Sub(bounds image.Rectangle) *SubImage128 {
	return &SubImage128{
		img:    img,
		bounds: bounds,
	}
}
