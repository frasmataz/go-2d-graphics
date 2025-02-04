package main

import (
	"image/color"
	"math"

	"github.com/go-p5/p5"
	"github.com/viterin/vek"
)

type ball struct {
	pos      []float64
	velocity []float64
	r        float64
}

var (
	screenWidth  = 500
	screenHeight = 500
	ballCount    = 50
	balls        []*ball
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(screenHeight, screenWidth)
	p5.Background(color.Gray{Y: 220})

	for range ballCount {
		balls = append(balls, &ball{
			pos:      []float64{p5.Random(0, float64(screenWidth)), p5.Random(0, float64(screenWidth))},
			velocity: []float64{p5.Random(-2, 2), p5.Random(-2, 2)},
			r:        10,
		})
	}
}

func update() {
	for _, ball := range balls {
		// Bounce off walls
		// Left wall
		if ball.pos[0]-ball.r < 0 {
			ball.pos[0] = ball.r
			ball.velocity[0] = -ball.velocity[0]
		}

		// Right wall
		if ball.pos[0]+ball.r > float64(screenWidth) {
			ball.pos[0] = float64(screenWidth) - ball.r
			ball.velocity[0] = -ball.velocity[0]
		}

		// Top wall
		if ball.pos[1]-ball.r < 0 {
			ball.pos[1] = ball.r
			ball.velocity[1] = -ball.velocity[1]
		}

		// Bottom wall
		if ball.pos[1]+ball.r > float64(screenHeight) {
			ball.pos[1] = float64(screenHeight) - ball.r
			ball.velocity[1] = -ball.velocity[1]
		}

		// Naively iterate over all balls for collision
		for j, ball2 := range balls {
			// don't collide with self
			if i != j {
				deltaVec := vek.Sub(ball.pos, ball2.pos)
				deltaMag := vek.Abs()			
				minBump := vek.MulNumber(((ball.r + ball2.r)-deltaMag)/deltaMag, deltaVec)
				
				invMass1 = 1 / (math.Pi * math.Pow(ball.r, 2))
				invMass2 = 1 / (math.Pi * math.Pow(ball2.r, 2))

				vek.Add_Into(ball.pos, vek.MulNumber(invMass1/(invMass1 + invMass2), minBump))
				
			}


		// 		if distance(ball.pos[0], ball.pos[1], otherBall.pos[0], otherBall.pos[1]) < (ball.r + otherBall.r) {
		// 			// Calculate difference between centres
		// 			dAngle, dMag := vectorXYtoAngleMag(otherBall.pos[0]-ball.pos[0], otherBall.pos[1]-ball.pos[1])
		//
		// 			// Bump balls away from each other
		// 			targetBallDistance := ball.r + otherBall.r
		//
		// 			overlap := targetBallDistance - dMag
		//
		// 			thisBallBumpx, thisBallBumpy := vectorAngleMagtoXY(dAngle, overlap/2)
		// 			otherBallBumpx, otherBallBumpy := vectorAngleMagtoXY(dAngle, overlap/2)
		//
		// 			ball.pos[0] += thisBallBumpx
		// 			ball.pos[1] += thisBallBumpy
		//
		// 			otherBall.pos[0] += otherBallBumpx
		// 			otherBall.pos[1] += otherBallBumpy
		//
		// 			ball.va = ball.va + (dAngle + math.Pi)
		// 			otherBall.va = otherBall.va + dAngle
		// 		}
		// 	}
		// }

		vek.Add_Inplace(ball.pos, ball.velocity)
	}
}

func draw() {
	update()
	p5.StrokeWidth(2)
	p5.Stroke(color.Black)
	p5.Fill(color.RGBA{R: 255, G: 50, B: 180, A: 208})

	for _, ball := range balls {
		p5.Circle(ball.pos[0], ball.pos[1], ball.r*2)
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
