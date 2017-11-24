package wo

import (
	"github.com/faiface/pixel"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestLimits(t *testing.T) {
	as := assert.New(t)

	r := pixel.R(1, 2, 3, 4)
	a, c, b, d := Limits(r)

	as.Equal(float64(1), a)
	as.Equal(float64(2), b)
	as.Equal(float64(3), c)
	as.Equal(float64(4), d)
}

func TestShape(t *testing.T) {
	as := assert.New(t)

	r := pixel.R(1, 2, 300, 400)
	x, y, w, h := Shape(r)

	as.Equal(float64(1), x)
	as.Equal(float64(2), y)
	as.Equal(float64(299), w)
	as.Equal(float64(398), h)
}

func TestFit(t *testing.T) {

	cases := []struct {
		name      string
		src, dest pixel.Rect
	}{
		{"untranslated", pixel.R(0, 0, 10, 11), pixel.R(0, 0, 100, 110)},
		{"pretranslated", pixel.R(5, 6, 10, 11), pixel.R(0, 0, 100, 110)},
		{"posttranslated", pixel.R(0, 0, 10, 11), pixel.R(10, 10, 100, 110)},
		{"translated", pixel.R(10, 11, 300, 450), pixel.R(23, 67, 899, 9202)},
		{"large_to_small", pixel.R(0, 0, 1000, 1000), pixel.R(10, 10, 20, 20)},
		{"small_to_large", pixel.R(10, 10, 20, 20), pixel.R(0, 0, 1000, 1000)},
		{"realistic", pixel.R(0, 0, 2037, 768), pixel.R(0, 0, 800, 600)},
	}

	compare := func(a, b pixel.Vec) func() bool {
		return func() bool {
			return math.Abs(a.X-b.X) < 1 && math.Abs(a.Y-b.Y) < 1
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			as := assert.New(t)

			mat := fit(c.src, c.dest)

			projectBottomLeft := mat.Project(c.src.Min)
			projectTopRight := mat.Project(c.src.Max)

			as.Condition(compare(c.dest.Min, projectBottomLeft), "MIN name=%s mat=%v from=%v to=%v", c.name, mat, c.src.Min, c.dest.Min)
			as.Condition(compare(c.dest.Max, projectTopRight), "MAX name=%s mat=%v from=%v to=%v", c.name, mat, c.src.Max, c.dest.Max)
		})
	}

}

func TestDegToRad(t *testing.T) {
	as := assert.New(t)

	as.Equal(float64(0), DegToRad(0))
	as.Equal(math.Pi, DegToRad(180))
	as.Equal(2*math.Pi, DegToRad(360))
}

func TestRadToDeg(t *testing.T) {
	as := assert.New(t)

	as.Equal(float64(0), RadToDeg(0))
	as.Equal(float64(180), RadToDeg(math.Pi))
	as.Equal(float64(360), RadToDeg(2*math.Pi))
}
