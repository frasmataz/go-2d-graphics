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
	screenWidth  = 1000
	screenHeight = 1000
	ballCount    = 80
	balls        []*ball
	bounciness   = 0.9
	gravity      = 0.3
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(screenHeight, screenWidth)
	p5.Background(color.Gray{Y: 80})

	for range ballCount {
		size := p5.Random(20, 40)
		balls = append(balls, &ball{
			pos:      []float64{p5.Random(0, float64(screenWidth)), p5.Random(0, float64(screenWidth))},
			velocity: vek.MulNumber([]float64{p5.Random(-200, 200), p5.Random(-200, 200)}, 1/size),
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

				mass1 := (math.Pi * math.Pow(ball.r, 2) * 10)
				mass2 := (math.Pi * math.Pow(ball2.r, 2) * 10)

				invMass1 := 1 / (math.Pi * math.Pow(ball.r, 2) * 10)
				invMass2 := 1 / (math.Pi * math.Pow(ball2.r, 2) * 10)

				vek.Add_Inplace(ball.pos, vek.MulNumber(nudge, invMass1/(invMass1+invMass2)))
				vek.Sub_Inplace(ball2.pos, vek.MulNumber(nudge, invMass1/(invMass1+invMass2)))

				// Elastic collisions - ported from https://www.plasmaphysics.org.uk/programs/coll2d_cpp.htm
				// Precalc some values
				m21 := mass2 / mass1
				x21 := ball2.pos[0] - ball.pos[0]
				y21 := ball2.pos[1] - ball.pos[1]
				vx21 := ball2.velocity[0] - ball.velocity[0]
				vy21 := ball2.velocity[1] - ball.velocity[1]

				// If balls are approaching:
				if (vx21*x21 + vy21*y21) < 0 {
					a := y21 / x21
					dvx2 := -2 * (vx21 + a*vy21) / ((1 + a*a) * (1 + m21))
					vx2 := ball2.velocity[0] + dvx2
					vy2 := ball2.velocity[1] + a*dvx2
					vx1 := ball.velocity[0] - m21*dvx2
					vy1 := ball.velocity[1] - a*m21*dvx2

					ball.velocity = []float64{vx1, vy1}
					ball2.velocity = []float64{vx2, vy2}

					ball1DampingVec := vek.DivNumber(vek.MulNumber(ball.velocity, mass1), (mass1 + mass2))
					ball2DampingVec := vek.DivNumber(vek.MulNumber(ball2.velocity, mass2), (mass1 + mass2))

					ball.velocity = vek.Add(vek.MulNumber(vek.Sub(ball.velocity, ball1DampingVec), bounciness), ball1DampingVec)
					ball2.velocity = vek.Add(vek.MulNumber(vek.Sub(ball2.velocity, ball2DampingVec), bounciness), ball2DampingVec)
				}
			}
		}
		vek.Add_Inplace(ball.pos, ball.velocity)

		// Apply gravity
		vek.Add_Inplace(ball.velocity, []float64{0.0, gravity})
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
