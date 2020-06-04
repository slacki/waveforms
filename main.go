package main

import (
	"fmt"
	"image/color"

	"github.com/slacki/waveforms/spectogram"
)

func main() {
	fmt.Println("start")

	s, err := spectogram.NewSpectogram(&spectogram.Config{
		BG0:    color.RGBA{33, 41, 201, 255},
		Width:  300,
		Height: 50,
	}, "test.wav")
	if err != nil {
		panic(err)
	}

	img, err := s.Generate()
	if err != nil {
		panic(err)
	}

	img.ToPNG("./test.png")
}
