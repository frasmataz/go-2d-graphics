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
	image  image.Image
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
	targetFrameTime     = 20 * time.Millisecond
	targetFPS           = 60
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
	rl.InitWindow(screenWidth, screenHeight, "mandelbrot raylib")
	rl.SetTargetFPS(60)

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
			img := rl.NewImageFromImage(seg.image)
			tex := rl.LoadTextureFromImage(img)
			rl.DrawTexture(tex, int32(seg.x), int32(seg.y), rl.White)
		}
	}

	rl.EndDrawing()

	// p5.TextSize(50)
	// // p5.Text(fmt.Sprintf("%v fps", displayedFrameCount), 50, 50)
	// frameCount++

	// if time.Now().After(nextFpsUpdate) {
	// 	displayedFrameCount = frameCount
	// 	frameCount = 0
	// 	nextFpsUpdate = time.Now().Add(fpsUpdateInterval)
	// }

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
		time.Sleep(time.Second * 1)
	} else if rl.IsKeyDown(rl.KeyPageDown) {
		mandelbrot.zoom += (zoomSpeed * mandelbrot.zoom)
		time.Sleep(time.Second * 1)
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
