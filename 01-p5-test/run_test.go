package main

import (
	"math"
	"testing"
)

func TestDistance(t *testing.T) {
	type coords struct {
		x1 float64
		x2 float64
		y1 float64
		y2 float64
	}

	type test struct {
		input coords
		want  float64
	}

	tests := []test{
		{
			input: coords{
				6.0, 8.0, 0.0, 0.0,
			},
			want: 10.0,
		},
		{
			input: coords{
				-3.0, 0.0, 0.0, -4.0,
			},
			want: 5.0,
		},
		{
			input: coords{
				0.0, 0.0, 0.0, 0.0,
			},
			want: 0.0,
		},
	}

	for _, test := range tests {
		got := distance(test.input.x1, test.input.x2, test.input.y1, test.input.y2)
		if got != test.want {
			t.Errorf("expected %f, got %f", test.want, got)
		}
	}
}

func TestVectorXYtoAngleMag(t *testing.T) {
	type coords struct {
		x float64
		y float64
	}

	type angleMag struct {
		a float64
		m float64
	}

	type test struct {
		input coords
		want  angleMag
	}

	tests := []test{
		{
			input: coords{
				1.0, 1.0,
			},
			want: angleMag{
				math.Pi / 4,
				distance(1.0, 1.0, 0.0, 0.0),
			},
		},
		{
			input: coords{
				-1.0, -1.0,
			},
			want: angleMag{
				(math.Pi / 4) * 5,
				distance(1.0, 1.0, 0.0, 0.0),
			},
		},
	}

	for _, test := range tests {
		gotA, gotM := vectorXYtoAngleMag(test.input.x, test.input.y)
		if gotA != test.want.a || gotM != test.want.m {
			t.Errorf("expected a=%f, m=%f, got a=%f, m=%f", test.want.a, test.want.m, gotA, gotM)
		}
	}
}

func TestVectorAngleMagtoXY(t *testing.T) {
	type angleMag struct {
		a float64
		m float64
	}

	type coords struct {
		x float64
		y float64
	}

	type test struct {
		input angleMag
		want  coords
	}

	tests := []test{
		{
			want: coords{
				1.0, 1.0,
			},
			input: angleMag{
				math.Pi / 4,
				distance(1.0, 1.0, 0.0, 0.0),
			},
		},
		{
			want: coords{
				-1.0, -1.0,
			},
			input: angleMag{
				(math.Pi / 4) * 5,
				distance(1.0, 1.0, 0.0, 0.0),
			},
		},
	}

	// This test failed on float precision in dev, so compare with a small tolerance
	tolerance := 1e-9
	for _, test := range tests {
		gotX, gotY := vectorAngleMagtoXY(test.input.a, test.input.m)
		if math.Abs(gotX-test.want.x) > tolerance || math.Abs(gotY-test.want.y) > tolerance {
			t.Errorf("expected x=%f, y=%f, got x=%f, y=%f", test.want.x, test.want.y, gotX, gotY)
		}
	}
}
