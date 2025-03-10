package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type imageSegment struct {
	thread int32
	pixels []color.RGBA
	x      int32
	y      int32
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
	screenWidth         = int32(1800)
	screenHeight        = int32(1000)
	threads             = int32(24)
	segmentsx           = int32(0) //Calculated at setup
	segmentsy           = int32(0)
	segmentHeight       = int32(0)
	segmentWidth        = int32(0)
	targetFrameTime     = 20 * time.Millisecond
	targetFPS           = 60
	fpsUpdateInterval   = 1000 * time.Millisecond
	nextFpsUpdate       = time.Now().Add(fpsUpdateInterval)
	frameCount          = 0
	displayedFrameCount = 0

	panSpeed      = float64(0.2)
	zoomSpeed     = float64(0.05)
	imageChannel  = make(chan imageSegment)
	imageBuffer   = make([]*imageSegment, threads)
	screenTexture rl.Texture2D

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
	rl.InitWindow(screenWidth, screenHeight, "mandelbrot raylib")
	rl.SetTargetFPS(60)

	screenRect := image.Rect(0, 0, int(screenWidth), int(screenHeight))
	screenImg := image.NewRGBA(screenRect)

	screenTexture = rl.LoadTextureFromImage(
		rl.NewImageFromImage(screenImg),
	)

	// Determine most square-y canvas segmenting for number of threads
	// Get factors of number of threads
	var factors []int32
	for i := int32(1); i <= threads; i++ {
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

	segmentHeight = screenHeight / segmentsy
	segmentWidth = screenWidth / segmentsx

	fmt.Printf("threads: %v, seg x: %v, seg y: %v, factors: %v", threads, segmentsx, segmentsy, factors)

	for thread := range threads {
		go processSegment(thread, &imageChannel)
	}

	for !rl.WindowShouldClose() {
		draw()
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

	rl.BeginDrawing()

	for _, seg := range imageBuffer {
		if seg != nil {
			rl.UpdateTextureRec(
				screenTexture,
				rl.NewRectangle(float32(seg.x), float32(seg.y), float32(segmentWidth), float32(segmentHeight)),
				seg.pixels,
			)
		}
	}
	rl.DrawTexture(screenTexture, 0, 0, rl.White)

	rl.EndDrawing()
}

func processInput() {
	if rl.IsKeyDown(rl.KeyLeft) {
		mandelbrot.pos.x -= (panSpeed * mandelbrot.zoom)
	} else if rl.IsKeyDown(rl.KeyRight) {
		mandelbrot.pos.x += (panSpeed * mandelbrot.zoom)
	}

	if rl.IsKeyDown(rl.KeyUp) {
		mandelbrot.pos.y -= (panSpeed * mandelbrot.zoom)
	} else if rl.IsKeyDown(rl.KeyDown) {
		mandelbrot.pos.y += (panSpeed * mandelbrot.zoom)
	}

	if rl.IsKeyDown(rl.KeyPageUp) {
		mandelbrot.zoom -= (zoomSpeed * mandelbrot.zoom)
	} else if rl.IsKeyDown(rl.KeyPageDown) {
		mandelbrot.zoom += (zoomSpeed * mandelbrot.zoom)
	}

	if rl.IsKeyDown(rl.KeyHome) {
		mandelbrot.iterations *= 2
	} else if rl.IsKeyDown(rl.KeyEnd) {
		mandelbrot.iterations /= 2
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

func processSegment(thread int32, imageChannel *chan imageSegment) {
	pixels := make([]color.RGBA, segmentWidth*segmentHeight)

	xpos := thread % segmentsx
	ypos := thread / segmentsx

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
					int(segmentWidth),
					int(segmentHeight),
				},
			},
		)

		for x := range image.Rect.Size().X {
			for y := range image.Rect.Size().Y {
				segx := (segmentWidth * xpos) + int32(x) - (screenWidth / 2)
				segy := (segmentHeight * ypos) + int32(y) - (screenHeight / 2)

				mbx := ((float64(segx) / float64(segmentWidth)) * mandelbrot.zoom) + mandelbrot.pos.x
				mby := ((float64(segy) / float64(segmentHeight)) * mandelbrot.zoom) + mandelbrot.pos.y

				mbResult := renderMandelbrot(mbx, mby)

				color := color.RGBA{
					uint8(mbResult),
					uint8(mbResult * 2),
					uint8(mbResult * 8),
					255,
				}

				pixels[y*int(segmentWidth)+x] = color
			}
		}

		*imageChannel <- imageSegment{
			thread: thread,
			pixels: pixels,
			x:      segmentLeft,
			y:      segmentTop,
		}
	}
}
