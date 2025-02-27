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
	screenWidth         = 1800
	screenHeight        = 1000
	threads             = 24
	segmentsx           = 0 //Calculated at setup
	segmentsy           = 0
	targetFrameTime     = 20 * time.Millisecond
	fpsUpdateInterval   = 1000 * time.Millisecond
	nextFpsUpdate       = time.Now().Add(fpsUpdateInterval)
	frameCount          = 0
	displayedFrameCount = 0

	panSpeed     = float64(0.2)
	zoomSpeed    = float64(0.05)
	imageChannel = make(chan imageSegment)
	imageBuffer  = make([]*imageSegment, threads)

	mandelbrot = mandelbrotState{
		iterations: 256,
		zoom:       0.5,
		pos: struct {
			x float64
			y float64
		}{-0.5, -0.0},
	}
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	// Determine most square-y canvas segmenting for number of threads
	// Get factors of number of threads
	var factors []int
	for i := 1; i <= threads; i++ {
		if threads%i == 0 {
			factors = append(factors, i)
		}
	}

	// Find midpoint of factors list
	if len(factors)%2 == 0 {
		// If even, take middle two factors as x/y segment count
		segmentsy = factors[len(factors)/2-1]
		segmentsx = factors[(len(factors) / 2)]
	} else {
		// If odd, threads is a square number - take middle factor for both x and y
		segmentsy = factors[(len(factors) / 2)]
		segmentsx = factors[(len(factors) / 2)]
	}

	fmt.Printf("threads: %v, seg x: %v, seg y: %v, factors: %v", threads, segmentsx, segmentsy, factors)

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
		if time.Now().After(nextFrameTime) {
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
		mandelbrot.zoom -= (zoomSpeed * mandelbrot.zoom)
	} else if p5.KeyIsDown(key.NamePageDown) {
		mandelbrot.zoom += (zoomSpeed * mandelbrot.zoom)
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
	xpos := thread % segmentsx
	ypos := thread / segmentsx

	segmentHeight := screenHeight / segmentsy
	segmentWidth := screenWidth / segmentsx

	segmentTop := segmentHeight * ypos
	segmentLeft := segmentWidth * xpos

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
				segx := (segmentWidth * xpos) + x - (screenWidth / 2)
				segy := (segmentHeight * ypos) + y - (screenHeight / 2)

				mbx := ((float64(segx) / float64(segmentWidth)) * mandelbrot.zoom) + mandelbrot.pos.x
				mby := ((float64(segy) / float64(segmentHeight)) * mandelbrot.zoom) + mandelbrot.pos.y

				mbResult := renderMandelbrot(mbx, mby)

				image.SetNRGBA(x, y, color.NRGBA{
					uint8(mbResult),
					uint8(mbResult * 2),
					uint8(mbResult * 8),
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
