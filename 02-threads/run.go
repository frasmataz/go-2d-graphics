package main

import (
	"image/color"
	"time"

	"github.com/go-p5/p5"
)

var (
	screenWidth  = 1000
	screenHeight = 800
	threads      = 12
	framerate    = 1000 / 60
	pixelpipe    chan int
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(screenWidth, screenHeight)
	p5.Background(color.Gray{Y: 80})
}

func draw() {
	startTime := time.Now().UnixMilli()

	for {
		imageSegment, ok := <-pixelpipe

		if ok == false {
			break
		}

		p5.DrawImage(imageSegment)

		endTime := time.Now().UnixMilli()
		if endTime-startTime > framerate {
			continue
		}
	}
}
