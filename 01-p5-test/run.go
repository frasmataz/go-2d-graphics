package main

import (
	"image/color"

	"github.com/go-p5/p5"
)

type ball struct {
	x float64
	y float64
	r float64
}

var (
	screenWidth  = 500
	screenHeight = 500
	ballCount    = 50
	balls        []ball
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(screenHeight, screenWidth)
	p5.Background(color.Gray{Y: 220})

	for range ballCount {
		balls = append(balls, ball{
			x: p5.Random(0, float64(screenWidth)),
			y: p5.Random(0, float64(screenWidth)),
			r: 10,
		})
	}
}

func draw() {
	p5.StrokeWidth(2)
	p5.Stroke(color.Black)
	p5.Fill(color.RGBA{R: 255, G: 50, B: 180, A: 208})

	for _, ball := range balls {
		p5.Circle(ball.x, ball.y, ball.r*2)
	}
}
