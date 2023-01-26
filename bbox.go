package gotrace

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
)

type interval_s struct {
	Min float64
	Max float64
}
type interval_t interval_s

func interval(i *interval_t, min float64, max float64) {
	i.Min = min
	i.Max = max
}
func singleton(i *interval_t, x float64) {
	interval(i, x, x)
}
func extend(i *interval_t, x float64) {
	if x < i.Min {
		i.Min = x
	} else if x > i.Max {
		i.Max = x
	}
}
func in_interval(i *interval_t, x float64) int {
	return int(libc.BoolToInt(i.Min <= x && x <= i.Max))
}
func iprod(a DPoint, b DPoint) float64 {
	return a.X*b.X + a.Y*b.Y
}
func bezier(t float64, x0 float64, x1 float64, x2 float64, x3 float64) float64 {
	var s float64 = 1 - t
	return s*s*s*x0 + (s*s*t)*3*x1 + (t*t*s)*3*x2 + t*t*t*x3
}
func bezier_limits(x0 float64, x1 float64, x2 float64, x3 float64, i *interval_t) {
	var (
		a float64
		b float64
		c float64
		d float64
		r float64
		t float64
		x float64
	)
	extend(i, x3)
	if in_interval(i, x1) != 0 && in_interval(i, x2) != 0 {
		return
	}
	a = x0*float64(-3) + x1*9 - x2*9 + x3*3
	b = x0*6 - x1*12 + x2*6
	c = x0*float64(-3) + x1*3
	d = b*b - a*4*c
	if d > 0 {
		r = math.Sqrt(d)
		t = (-b - r) / (a * 2)
		if t > 0 && t < 1 {
			x = bezier(t, x0, x1, x2, x3)
			extend(i, x)
		}
		t = (-b + r) / (a * 2)
		if t > 0 && t < 1 {
			x = bezier(t, x0, x1, x2, x3)
			extend(i, x)
		}
	}
	return
}
func segment_limits(tag int, a DPoint, c [3]DPoint, dir DPoint, i *interval_t) {
	switch tag {
	case POTRACE_CORNER:
		extend(i, iprod(c[1], dir))
		extend(i, iprod(c[2], dir))
	case POTRACE_CURVETO:
		bezier_limits(iprod(a, dir), iprod(c[0], dir), iprod(c[1], dir), iprod(c[2], dir), i)
	}
}
func curve_limits(curve *Curve, dir DPoint, i *interval_t) {
	var (
		k int
		n int = curve.N
	)
	segment_limits(curve.Tag[0], curve.C[n-1][2], curve.C[0], dir, i)
	for k = 1; k < n; k++ {
		segment_limits(curve.Tag[k], curve.C[k-1][2], curve.C[k], dir, i)
	}
}
func path_limits(path *Path, dir DPoint, i *interval_t) {
	var p *Path
	if path == nil {
		interval(i, 0, 0)
		return
	}
	singleton(i, iprod(path.Curve.C[0][2], dir))
	for p = path; p != nil; p = p.Next {
		curve_limits(&p.Curve, dir, i)
	}
}
