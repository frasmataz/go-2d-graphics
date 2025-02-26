package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"gioui.org/io/key"
	"github.com/frasmataz/p5"
)

type imageSegment struct {
	thread int
	image  image.Image
	x      int
	y      int
}

type mandelbrotState struct {
	iterations uint
	zoom       float64
	pos        struct {
		x float64
		y float64
	}

	panSpeed  float64
	zoomSpeed float64
}

var (
	screenWidth         = 1200
	screenHeight        = 800
	threads             = 24
	targetFrameTime     = 20 * time.Millisecond
	fpsUpdateInterval   = 1000 * time.Millisecond
	nextFpsUpdate       = time.Now().Add(fpsUpdateInterval)
	frameCount          = 0
	displayedFrameCount = 0

	panSpeed     = float64(0.1)
	zoomSpeed    = float64(0.05)
	imageChannel = make(chan imageSegment)
	imageBuffer  = make([]*imageSegment, threads)

	mandelbrot = mandelbrotState{
		iterations: 200,
		zoom:       1.5,
		pos: struct {
			x float64
			y float64
		}{-2.0, -0.5},
	}
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(screenWidth, screenHeight)
	for thread := range threads {
		go processSegment(thread, &imageChannel)
	}
}

func draw() {
	processInput()
	now := time.Now()
	nextFrameTime := now.Add(targetFrameTime)

	for {
		now := time.Now()
		if nextFrameTime.Before(now) {
			break
		}

		imageSeg := <-imageChannel
		imageBuffer[imageSeg.thread] = &imageSeg
	}

	for _, seg := range imageBuffer {
		if seg != nil {
			p5.DrawImage(seg.image, float64(seg.x), float64(seg.y))
		}
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

func processInput() {
	if p5.KeyIsDown(key.NameLeftArrow) {
		mandelbrot.pos.x -= (panSpeed * mandelbrot.zoom)
	} else if p5.KeyIsDown(key.NameRightArrow) {
		mandelbrot.pos.x += (panSpeed * mandelbrot.zoom)
	}

	if p5.KeyIsDown(key.NameUpArrow) {
		mandelbrot.pos.y -= (panSpeed * mandelbrot.zoom)
	} else if p5.KeyIsDown(key.NameDownArrow) {
		mandelbrot.pos.y += (panSpeed * mandelbrot.zoom)
	}

	if p5.KeyIsDown(key.NamePageUp) {
		mandelbrot.zoom += (zoomSpeed * mandelbrot.zoom)
	} else if p5.KeyIsDown(key.NamePageDown) {
		mandelbrot.zoom -= (zoomSpeed * mandelbrot.zoom)
	}

}

func renderMandelbrot(x float64, y float64) uint {
	cursorx := float64(0.0)
	cursory := float64(0.0)
	iteration := uint(0)

	for {
		cursorxsq := cursorx * cursorx
		cursorysq := cursory * cursory
		if cursorxsq+cursorysq >= 4.0 || iteration >= mandelbrot.iterations {
			return iteration
		}

		xtemp := cursorxsq - cursorysq + x
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

	for {
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

				mbx := ((float64(segx) / float64(segmentWidth) * 1.5) * mandelbrot.zoom) + mandelbrot.pos.x
				mby := ((float64(segy) / float64(screenHeight)) * mandelbrot.zoom) + mandelbrot.pos.y

				mbResult := renderMandelbrot(mbx, mby)

				image.SetNRGBA(x, y, color.NRGBA{
					uint8(mbResult * 50),
					uint8(mbResult * 10),
					uint8(mbResult * 70),
					255,
				})
			}
		}

		*imageChannel <- imageSegment{
			thread: thread,
			image:  image,
			x:      segmentLeft,
			y:      segmentTop,
		}
	}
}
