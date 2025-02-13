package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/go-p5/p5"
)

type imageSegment struct {
	image image.Image
	x     int
	y     int
}

var (
	screenWidth         = 1000
	screenHeight        = 800
	threads             = 12
	fpsUpdateInterval   = 1000 * time.Millisecond
	nextFpsUpdate       = time.Now().Add(fpsUpdateInterval)
	frameCount          = 0
	displayedFrameCount = 0
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(screenWidth, screenHeight)
	p5.Background(color.Gray{Y: 80})
}

func draw() {
	imageChannels := make([]chan imageSegment, threads)

	for thread := range threads {
		imageChannels[thread] = make(chan imageSegment)

		go func(thread int) {
			segmentHeight := screenHeight / threads
			segmentWidth := screenWidth

			segmentTop := segmentHeight * thread
			segmentLeft := 0

			image := image.NewNRGBA(
				image.Rectangle{
					image.Point{
						0,
						0,
					},
					image.Point{
						segmentWidth,
						segmentHeight,
					},
				},
			)

			for x := range image.Rect.Size().X {
				for y := range image.Rect.Size().Y {
					image.SetNRGBA(x, y, color.NRGBA{
						uint8(x + segmentTop + y),
						uint8(x + segmentTop - y),
						uint8(x + segmentTop*y),
						255,
					})
				}
			}

			imageChannels[thread] <- imageSegment{
				image: image,
				x:     segmentLeft,
				y:     segmentTop,
			}

			close(imageChannels[thread])
		}(thread)
	}

	for thread := range threads {
		imageSeg := <-imageChannels[thread]
		p5.DrawImage(imageSeg.image, float64(imageSeg.x), float64(imageSeg.y))
	}

	p5.TextSize(50)
	p5.Text(fmt.Sprintf("%v fps", displayedFrameCount), 50, 50)
	frameCount++

	if time.Now().After(nextFpsUpdate) {
		displayedFrameCount = frameCount
		frameCount = 0
		nextFpsUpdate = time.Now().Add(fpsUpdateInterval)
	}
}
