package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/go-p5/p5"
)

type imageSegment struct {
	image image.Image
	x     int
	y     int
}

var (
	screenWidth         = 1200
	screenHeight        = 800
	threads             = 24
	fpsUpdateInterval   = 1000 * time.Millisecond
	nextFpsUpdate       = time.Now().Add(fpsUpdateInterval)
	frameCount          = 0
	displayedFrameCount = 0

	mandelbrotScaleX     = float64(1.7)
	mandelbrotScaleY     = float64(1.5)
	mandelbrotPositionX  = float64(-2.0)
	mandelbrotPositionY  = float64(-0.5)
	mandelbrotIterations = 200
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(screenWidth, screenHeight)
}

func draw() {
	imageChannels := make([]chan imageSegment, threads)

	for thread := range threads {
		imageChannels[thread] = make(chan imageSegment)

		go processSegment(thread, &imageChannels[thread])
	}

	for thread := range threads {
		imageSeg := <-imageChannels[thread]
		p5.DrawImage(imageSeg.image, float64(imageSeg.x), float64(imageSeg.y))
	}

	p5.TextSize(50)
	p5.Fill(color.NRGBA{
		200,
		200,
		200,
		255,
	})
	p5.Text(fmt.Sprintf("%v fps", displayedFrameCount), 50, 50)
	frameCount++

	if time.Now().After(nextFpsUpdate) {
		displayedFrameCount = frameCount
		frameCount = 0
		nextFpsUpdate = time.Now().Add(fpsUpdateInterval)
	}
}

func mandelbrot(x float64, y float64) int {
	cursorx := float64(0.0)
	cursory := float64(0.0)
	iteration := 0

	for {
		if math.Pow(cursorx, 2)+math.Pow(cursory, 2) >= 4.0 || iteration >= mandelbrotIterations {
			return iteration
		}

		xtemp := math.Pow(cursorx, 2) - math.Pow(cursory, 2) + x
		cursory = 2*cursorx*cursory + y
		cursorx = xtemp
		iteration++
	}
}

func processSegment(thread int, imageChannel *chan imageSegment) {
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
			segx := x
			segy := (segmentHeight * thread) + y

			mbx := ((float64(segx) / float64(segmentWidth)) * mandelbrotScaleX) + mandelbrotPositionX
			mby := ((float64(segy) / float64(screenHeight)) * mandelbrotScaleY) + mandelbrotPositionY

			mbResult := mandelbrot(mbx, mby)

			image.SetNRGBA(x, y, color.NRGBA{
				uint8(mbResult * 50),
				uint8(mbResult * 10),
				uint8(mbResult * 70),
				255,
			})
		}
	}

	*imageChannel <- imageSegment{
		image: image,
		x:     segmentLeft,
		y:     segmentTop,
	}

	close(*imageChannel)
}
