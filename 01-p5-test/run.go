package main

import (
	"image/color"
	"log"
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
	screenWidth  = 1000
	screenHeight = 1000
	ballCount    = 50
	balls        []*ball
	bounciness   = 1.0001
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(screenHeight, screenWidth)
	p5.Background(color.Gray{Y: 80})

	for range ballCount {
		size := p5.Random(20, 80)
		balls = append(balls, &ball{
			pos:      []float64{p5.Random(0, float64(screenWidth)), p5.Random(0, float64(screenWidth))},
			velocity: vek.MulNumber([]float64{p5.Random(-100, 100), p5.Random(-100, 100)}, 1/size),
			r:        size,
		})
	}
}

func update() {
	for i, ball := range balls {
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
			if i != j && vek.Norm(vek.Sub(ball.pos, ball2.pos)) < ball.r+ball2.r {
				// Get vector between circle centres
				deltaVec := vek.Sub(ball.pos, ball2.pos)
				deltaMag := vek.Norm(deltaVec)

				// Bump circles off each other
				nudge := vek.MulNumber(deltaVec, ((ball.r+ball2.r)-deltaMag)/deltaMag)

				invMass1 := 1 / (math.Pi * math.Pow(ball.r, 2) * 10)
				invMass2 := 1 / (math.Pi * math.Pow(ball2.r, 2*10))

				vek.Add_Inplace(ball.pos, vek.MulNumber(nudge, invMass1/(invMass1+invMass2)))
				vek.Sub_Inplace(ball2.pos, vek.MulNumber(nudge, invMass1/(invMass1+invMass2)))

				// Ricochet balls off each other
				deltaVelocity := vek.Sub(ball.velocity, ball2.velocity)
				deltaVelocityMag := vek.Dot(nudge, deltaVelocity)

				log.Printf("minBump: %v, delta: %v", nudge, deltaVelocity)

				// If balls not moving away from each other already
				if deltaVelocityMag < 0.0 {
					log.Println("collide")
					impulseForce := (deltaVelocityMag) / (invMass1 + invMass2)
					impulseVec := vek.MulNumber(nudge, impulseForce)

					log.Printf("impulseForce: %v", impulseForce)
					log.Printf("impulseVec: %v", impulseVec)

					log.Printf("vel before: %v", ball.velocity)
					vek.Add_Inplace(ball.velocity, vek.MulNumber(impulseVec, invMass1))
					vek.Sub_Inplace(ball2.velocity, vek.MulNumber(impulseVec, invMass2))
					log.Printf("vel after: %v", ball.velocity)
				}

				log.Println("done colliding")
			}
		}
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
