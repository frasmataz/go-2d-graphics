package main

import (
	"image/color"
	"math"

	"github.com/go-p5/p5"
)

type ball struct {
	x  float64
	y  float64
	r  float64
	vx float64
	vy float64
}

var (
	screenWidth  = 500
	screenHeight = 500
	ballCount    = 50
	balls        [50]*ball
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(screenHeight, screenWidth)
	p5.Background(color.Gray{Y: 220})

	for i := range ballCount {
		balls[i] = &ball{
			x:  p5.Random(0, float64(screenWidth)),
			y:  p5.Random(0, float64(screenWidth)),
			r:  10,
			vx: p5.Random(-2, 2),
			vy: p5.Random(-2, 2),
		}
	}
}

func update() {
	for i, ball := range balls {
		// Bounce of walls
		if ball.x < 0 || ball.x > float64(screenWidth) {
			ball.vx = -ball.vx
		}
		if ball.y < 0 || ball.y > float64(screenWidth) {
			ball.vy = -ball.vy
		}

		ball.x += ball.vx
		ball.y += ball.vy

		// Naively iterate over all balls for collision
		for j, otherBall := range balls {
			// don't collide with self
			if i != j {
				if distance(ball.x, ball.y, otherBall.x, otherBall.y) < (ball.r + otherBall.r) {
					vAngle, vMag := vectorXYtoAngleMag(ball.vx, ball.vy)
					dAngle, _ := vectorXYtoAngleMag(otherBall.x-ball.x, otherBall.y-ball.y)

					newVAngle := vAngle + (dAngle + math.Pi)

					newVx, newVy := vectorAngleMagtoXY(newVAngle, vMag)
					ball.vx = newVx
					ball.vy = newVy
					ball.x += ball.vx
					ball.y += ball.vy
				}
			}
		}
	}
}

func draw() {
	update()
	p5.StrokeWidth(2)
	p5.Stroke(color.Black)
	p5.Fill(color.RGBA{R: 255, G: 50, B: 180, A: 208})

	for _, ball := range balls {
		p5.Circle(ball.x, ball.y, ball.r*2)
	}
}

func distance(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	return math.Sqrt(math.Pow((x2-x1), 2) + math.Pow((y2-y1), 2))
}

func vectorXYtoAngleMag(x float64, y float64) (a float64, m float64) {
	a = math.Atan2(y, x)
	if a < 0 {
		a += math.Pi * 2
	}

	m = math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))

	return a, m
}

func vectorAngleMagtoXY(a float64, m float64) (x float64, y float64) {
	x = math.Sin(a) * m
	y = math.Cos(a) * m
	return x, y
}
