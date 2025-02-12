package main

import (
	"image"
	"image/color"
	"log"
	"time"

	"github.com/go-p5/p5"
)

type imageSegment struct {
	image image.Image
	x     int
	y     int
}

var (
	screenWidth  = 1000
	screenHeight = 800
	threads      = 12
	framerate    = int64(1000 / 60)
	pixelpipe    chan imageSegment
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	pixelpipe = make(chan imageSegment)
	p5.Canvas(screenWidth, screenHeight)
	p5.Background(color.Gray{Y: 80})
	for thread := range threads {
		log.Printf("thread %v starting", thread)
		go func(thread int) {
			time.Sleep(time.Second)
			log.Printf("thread %v started", thread)
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
						image.SetNRGBA(x, y, color.NRGBA{
							uint8(x + segmentTop + y),
							uint8(x + segmentTop - y),
							uint8(x + segmentTop*y),
							255,
						})
					}
				}

				log.Printf("thread %v posting", thread)
				pixelpipe <- imageSegment{
					image: image,
					x:     segmentLeft,
					y:     segmentTop,
				}
			}
		}(thread)
	}
}

func draw() {
	for {
		log.Println("draw starting")
		startTime := time.Now().UnixMilli()

		imageSegment := <-pixelpipe

		p5.DrawImage(imageSegment.image, float64(imageSegment.x), float64(imageSegment.y))

		endTime := time.Now().UnixMilli()
		log.Printf("start: %v, end: %v, framerate: %v", startTime, endTime, framerate)
		if endTime-startTime > framerate {
			return
		}
	}
}
