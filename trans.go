package gotrace

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

type transT struct {
	Bb     [2]float64
	Orig   [2]float64
	X      [2]float64
	Y      [2]float64
	Scalex float64
	Scaley float64
}

func trans(p DPoint, t transT) DPoint {
	var res DPoint
	res.X = t.Orig[0] + p.X*t.X[0] + p.Y*t.Y[0]
	res.Y = t.Orig[1] + p.X*t.X[1] + p.Y*t.Y[1]
	return res
}
func trans_rotate(r *transT, alpha float64) {
	var (
		s        float64
		c        float64
		x0       float64
		x1       float64
		y0       float64
		y1       float64
		o0       float64
		o1       float64
		t_struct transT
		t        *transT = &t_struct
	)
	libc.MemCpy(unsafe.Pointer(t), unsafe.Pointer(r), int(unsafe.Sizeof(transT{})))
	s = math.Sin(alpha / 180 * math.Pi)
	c = math.Cos(alpha / 180 * math.Pi)
	x0 = c * t.Bb[0]
	x1 = s * t.Bb[0]
	y0 = -s * t.Bb[1]
	y1 = c * t.Bb[1]
	r.Bb[0] = math.Abs(x0) + math.Abs(y0)
	r.Bb[1] = math.Abs(x1) + math.Abs(y1)
	o0 = -(func() float64 {
		if x0 < 0 {
			return x0
		}
		return 0
	}()) - (func() float64 {
		if y0 < 0 {
			return y0
		}
		return 0
	}())
	o1 = -(func() float64 {
		if x1 < 0 {
			return x1
		}
		return 0
	}()) - (func() float64 {
		if y1 < 0 {
			return y1
		}
		return 0
	}())
	r.Orig[0] = o0 + c*t.Orig[0] - s*t.Orig[1]
	r.Orig[1] = o1 + s*t.Orig[0] + c*t.Orig[1]
	r.X[0] = c*t.X[0] - s*t.X[1]
	r.X[1] = s*t.X[0] + c*t.X[1]
	r.Y[0] = c*t.Y[0] - s*t.Y[1]
	r.Y[1] = s*t.Y[0] + c*t.Y[1]
}
func trans_from_rect(r *transT, w float64, h float64) {
	r.Bb[0] = w
	r.Bb[1] = h
	r.Orig[0] = 0.0
	r.Orig[1] = 0.0
	r.X[0] = 1.0
	r.X[1] = 0.0
	r.Y[0] = 0.0
	r.Y[1] = 1.0
	r.Scalex = 1.0
	r.Scaley = 1.0
}
func trans_rescale(r *transT, sc float64) {
	r.Bb[0] *= sc
	r.Bb[1] *= sc
	r.Orig[0] *= sc
	r.Orig[1] *= sc
	r.X[0] *= sc
	r.X[1] *= sc
	r.Y[0] *= sc
	r.Y[1] *= sc
	r.Scalex *= sc
	r.Scaley *= sc
}
func trans_scale_to_size(r *transT, w float64, h float64) {
	var (
		xsc float64 = w / r.Bb[0]
		ysc float64 = h / r.Bb[1]
	)
	r.Bb[0] = w
	r.Bb[1] = h
	r.Orig[0] *= xsc
	r.Orig[1] *= ysc
	r.X[0] *= xsc
	r.X[1] *= ysc
	r.Y[0] *= xsc
	r.Y[1] *= ysc
	r.Scalex *= xsc
	r.Scaley *= ysc
	if w < 0 {
		r.Orig[0] -= w
		r.Bb[0] = -w
	}
	if h < 0 {
		r.Orig[1] -= h
		r.Bb[1] = -h
	}
}
func trans_tighten(r *transT, plist *Path) {
	var (
		i   interval_t
		dir DPoint
		j   int
	)
	if plist == nil {
		return
	}
	for j = 0; j < 2; j++ {
		dir.X = r.X[j]
		dir.Y = r.Y[j]
		path_limits(plist, dir, &i)
		if i.Min == i.Max {
			i.Max = i.Min + 0.5
			i.Min = i.Min - 0.5
		}
		r.Bb[j] = i.Max - i.Min
		r.Orig[j] = -i.Min
	}
}
