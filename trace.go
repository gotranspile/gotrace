package gotrace

import (
	"math"
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/libc"
)

const infty = 10000000

func dorth_infty(p0 DPoint, p2 DPoint) Point {
	var r Point
	if (p2.X - p0.X) > 0 {
		r.Y = 1
	} else if (p2.X - p0.X) < 0 {
		r.Y = -1
	} else {
		r.Y = 0
	}
	if (p2.Y - p0.Y) > 0 {
		r.X = -1
	} else if (p2.Y - p0.Y) < 0 {
		r.X = 1
	} else {
		r.X = 0
	}
	return r
}
func dpara(p0 DPoint, p1 DPoint, p2 DPoint) float64 {
	var (
		x1 float64
		y1 float64
		x2 float64
		y2 float64
	)
	x1 = p1.X - p0.X
	y1 = p1.Y - p0.Y
	x2 = p2.X - p0.X
	y2 = p2.Y - p0.Y
	return x1*y2 - x2*y1
}
func ddenom(p0 DPoint, p2 DPoint) float64 {
	var r Point = dorth_infty(p0, p2)
	return float64(r.Y)*(p2.X-p0.X) - float64(r.X)*(p2.Y-p0.Y)
}
func cyclic(a int, b int, c int) int {
	if a <= c {
		return int(libc.BoolToInt(a <= b && b < c))
	} else {
		return int(libc.BoolToInt(a <= b || b < c))
	}
}
func pointslope(pp *potrace_privpath_s, i int, j int, ctr *DPoint, dir *DPoint) {
	var (
		n       int     = pp.Len
		sums    *sums_s = &pp.Sums[0]
		x       float64
		y       float64
		x2      float64
		xy      float64
		y2      float64
		k       float64
		a       float64
		b       float64
		c       float64
		lambda2 float64
		l       float64
		r       int = 0
	)
	for j >= n {
		j -= n
		r += 1
	}
	for i >= n {
		i -= n
		r -= 1
	}
	for j < 0 {
		j += n
		r -= 1
	}
	for i < 0 {
		i += n
		r += 1
	}
	x = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).X - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).X + float64(r)*(*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(n)))).X
	y = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).Y - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).Y + float64(r)*(*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(n)))).Y
	x2 = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).X2 - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).X2 + float64(r)*(*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(n)))).X2
	xy = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).Xy - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).Xy + float64(r)*(*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(n)))).Xy
	y2 = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).Y2 - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).Y2 + float64(r)*(*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(n)))).Y2
	k = float64(j + 1 - i + r*n)
	ctr.X = x / k
	ctr.Y = y / k
	a = (x2 - x*x/k) / k
	b = (xy - x*y/k) / k
	c = (y2 - y*y/k) / k
	lambda2 = (a + c + math.Sqrt((a-c)*(a-c)+b*4*b)) / 2
	a -= lambda2
	c -= lambda2
	if math.Abs(a) >= math.Abs(c) {
		l = math.Sqrt(a*a + b*b)
		if l != 0 {
			dir.X = -b / l
			dir.Y = a / l
		}
	} else {
		l = math.Sqrt(c*c + b*b)
		if l != 0 {
			dir.X = -c / l
			dir.Y = b / l
		}
	}
	if l == 0 {
		dir.X = func() float64 {
			p := &dir.Y
			dir.Y = 0
			return *p
		}()
	}
}

type quadform_t [3][3]float64

func quadform(Q quadform_t, w DPoint) float64 {
	var (
		v   [3]float64
		i   int
		j   int
		sum float64
	)
	v[0] = w.X
	v[1] = w.Y
	v[2] = 1
	sum = 0.0
	for i = 0; i < 3; i++ {
		for j = 0; j < 3; j++ {
			sum += v[i] * Q[i][j] * v[j]
		}
	}
	return sum
}
func xprod(p1 Point, p2 Point) int {
	return p1.X*p2.Y - p1.Y*p2.X
}
func cprod(p0 DPoint, p1 DPoint, p2 DPoint, p3 DPoint) float64 {
	var (
		x1 float64
		y1 float64
		x2 float64
		y2 float64
	)
	x1 = p1.X - p0.X
	y1 = p1.Y - p0.Y
	x2 = p3.X - p2.X
	y2 = p3.Y - p2.Y
	return x1*y2 - x2*y1
}
func trace_iprod(p0 DPoint, p1 DPoint, p2 DPoint) float64 {
	var (
		x1 float64
		y1 float64
		x2 float64
		y2 float64
	)
	x1 = p1.X - p0.X
	y1 = p1.Y - p0.Y
	x2 = p2.X - p0.X
	y2 = p2.Y - p0.Y
	return x1*x2 + y1*y2
}
func iprod1(p0 DPoint, p1 DPoint, p2 DPoint, p3 DPoint) float64 {
	var (
		x1 float64
		y1 float64
		x2 float64
		y2 float64
	)
	x1 = p1.X - p0.X
	y1 = p1.Y - p0.Y
	x2 = p3.X - p2.X
	y2 = p3.Y - p2.Y
	return x1*x2 + y1*y2
}
func ddist(p DPoint, q DPoint) float64 {
	return math.Sqrt(((p.X - q.X) * (p.X - q.X)) + (p.Y-q.Y)*(p.Y-q.Y))
}
func trace_bezier(t float64, p0 DPoint, p1 DPoint, p2 DPoint, p3 DPoint) DPoint {
	var (
		s   float64 = 1 - t
		res DPoint
	)
	res.X = s*s*s*p0.X + (s*s*t)*3*p1.X + (t*t*s)*3*p2.X + t*t*t*p3.X
	res.Y = s*s*s*p0.Y + (s*s*t)*3*p1.Y + (t*t*s)*3*p2.Y + t*t*t*p3.Y
	return res
}
func tangent(p0 DPoint, p1 DPoint, p2 DPoint, p3 DPoint, q0 DPoint, q1 DPoint) float64 {
	var (
		A  float64
		B  float64
		C  float64
		a  float64
		b  float64
		c  float64
		d  float64
		s  float64
		r1 float64
		r2 float64
	)
	A = cprod(p0, p1, q0, q1)
	B = cprod(p1, p2, q0, q1)
	C = cprod(p2, p3, q0, q1)
	a = A - B*2 + C
	b = A*float64(-2) + B*2
	c = A
	d = b*b - a*4*c
	if a == 0 || d < 0 {
		return -1.0
	}
	s = math.Sqrt(d)
	r1 = (-b + s) / (a * 2)
	r2 = (-b - s) / (a * 2)
	if r1 >= 0 && r1 <= 1 {
		return r1
	} else if r2 >= 0 && r2 <= 1 {
		return r2
	} else {
		return -1.0
	}
}
func calc_sums(pp *potrace_privpath_s) int {
	var (
		i int
		x int
		y int
		n int = pp.Len
	)
	if (func() []sums_s {
		p := &pp.Sums
		pp.Sums = make([]sums_s, pp.Len+1)
		return *p
	}()) == nil {
		goto calloc_error
	}
	pp.X0 = pp.Pt[0].X
	pp.Y0 = pp.Pt[0].Y
	pp.Sums[0].X2 = func() float64 {
		p := &pp.Sums[0].Xy
		pp.Sums[0].Xy = func() float64 {
			p := &pp.Sums[0].Y2
			pp.Sums[0].Y2 = func() float64 {
				p := &pp.Sums[0].X
				pp.Sums[0].X = func() float64 {
					p := &pp.Sums[0].Y
					pp.Sums[0].Y = 0
					return *p
				}()
				return *p
			}()
			return *p
		}()
		return *p
	}()
	for i = 0; i < n; i++ {
		x = pp.Pt[i].X - pp.X0
		y = pp.Pt[i].Y - pp.Y0
		pp.Sums[i+1].X = pp.Sums[i].X + float64(x)
		pp.Sums[i+1].Y = pp.Sums[i].Y + float64(y)
		pp.Sums[i+1].X2 = pp.Sums[i].X2 + float64(x)*float64(x)
		pp.Sums[i+1].Xy = pp.Sums[i].Xy + float64(x)*float64(y)
		pp.Sums[i+1].Y2 = pp.Sums[i].Y2 + float64(y)*float64(y)
	}
	return 0
calloc_error:
	return 1
}
func calc_lon(pp *potrace_privpath_s) int {
	var (
		pt         *Point = &pp.Pt[0]
		n          int    = pp.Len
		i          int
		j          int
		k          int
		k1         int
		ct         [4]int
		dir        int
		constraint [2]Point
		cur        Point
		off        Point
		pivk       *int = nil
		nc         *int = nil
		dk         Point
		a          int
		b          int
		c          int
		d          int
	)
	if (func() *int {
		pivk = &make([]int, n)[0]
		return pivk
	}()) == nil {
		goto calloc_error
	}
	if (func() *int {
		nc = &make([]int, n)[0]
		return nc
	}()) == nil {
		goto calloc_error
	}
	k = 0
	for i = n - 1; i >= 0; i-- {
		if (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).X != (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k)))).X && (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).Y != (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k)))).Y {
			k = i + 1
		}
		*(*int)(unsafe.Add(unsafe.Pointer(nc), unsafe.Sizeof(int(0))*uintptr(i))) = k
	}
	if (func() []int {
		p := &pp.Lon
		pp.Lon = make([]int, n)
		return *p
	}()) == nil {
		goto calloc_error
	}
	for i = n - 1; i >= 0; i-- {
		ct[0] = func() int {
			p := &ct[1]
			ct[1] = func() int {
				p := &ct[2]
				ct[2] = func() int {
					p := &ct[3]
					ct[3] = 0
					return *p
				}()
				return *p
			}()
			return *p
		}()
		dir = (((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(mod(i+1, n))))).X-(*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).X)*3 + 3 + ((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(mod(i+1, n))))).Y - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).Y)) / 2
		ct[dir]++
		constraint[0].X = 0
		constraint[0].Y = 0
		constraint[1].X = 0
		constraint[1].Y = 0
		k = *(*int)(unsafe.Add(unsafe.Pointer(nc), unsafe.Sizeof(int(0))*uintptr(i)))
		k1 = i
		for {
			dir = ((func() int {
				if ((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k)))).X - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k1)))).X) > 0 {
					return 1
				}
				if ((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k)))).X - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k1)))).X) < 0 {
					return -1
				}
				return 0
			}())*3 + 3 + (func() int {
				if ((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k)))).Y - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k1)))).Y) > 0 {
					return 1
				}
				if ((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k)))).Y - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k1)))).Y) < 0 {
					return -1
				}
				return 0
			}())) / 2
			ct[dir]++
			if ct[0] != 0 && ct[1] != 0 && ct[2] != 0 && ct[3] != 0 {
				*(*int)(unsafe.Add(unsafe.Pointer(pivk), unsafe.Sizeof(int(0))*uintptr(i))) = k1
				goto foundk
			}
			cur.X = (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k)))).X - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).X
			cur.Y = (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k)))).Y - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).Y
			if xprod(constraint[0], cur) < 0 || xprod(constraint[1], cur) > 0 {
				goto constraint_viol
			}
			if (func() int {
				if cur.X > 0 {
					return cur.X
				}
				return -cur.X
			}()) <= 1 && (func() int {
				if cur.Y > 0 {
					return cur.Y
				}
				return -cur.Y
			}()) <= 1 {
			} else {
				off.X = cur.X + int(func() int8 {
					if cur.Y >= 0 && (cur.Y > 0 || cur.X < 0) {
						return 1
					}
					return -1
				}())
				off.Y = cur.Y + int(func() int8 {
					if cur.X <= 0 && (cur.X < 0 || cur.Y < 0) {
						return 1
					}
					return -1
				}())
				if xprod(constraint[0], off) >= 0 {
					constraint[0] = off
				}
				off.X = cur.X + int(func() int8 {
					if cur.Y <= 0 && (cur.Y < 0 || cur.X < 0) {
						return 1
					}
					return -1
				}())
				off.Y = cur.Y + int(func() int8 {
					if cur.X >= 0 && (cur.X > 0 || cur.Y < 0) {
						return 1
					}
					return -1
				}())
				if xprod(constraint[1], off) <= 0 {
					constraint[1] = off
				}
			}
			k1 = k
			k = *(*int)(unsafe.Add(unsafe.Pointer(nc), unsafe.Sizeof(int(0))*uintptr(k1)))
			if cyclic(k, i, k1) == 0 {
				break
			}
		}
	constraint_viol:
		if ((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k)))).X - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k1)))).X) > 0 {
			dk.X = 1
		} else if ((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k)))).X - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k1)))).X) < 0 {
			dk.X = -1
		} else {
			dk.X = 0
		}
		if ((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k)))).Y - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k1)))).Y) > 0 {
			dk.Y = 1
		} else if ((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k)))).Y - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k1)))).Y) < 0 {
			dk.Y = -1
		} else {
			dk.Y = 0
		}
		cur.X = (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k1)))).X - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).X
		cur.Y = (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(k1)))).Y - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).Y
		a = xprod(constraint[0], cur)
		b = xprod(constraint[0], dk)
		c = xprod(constraint[1], cur)
		d = xprod(constraint[1], dk)
		j = infty
		if b < 0 {
			j = floordiv(a, -b)
		}
		if d > 0 {
			if j < floordiv(-c, d) {
				j = j
			} else {
				j = floordiv(-c, d)
			}
		}
		*(*int)(unsafe.Add(unsafe.Pointer(pivk), unsafe.Sizeof(int(0))*uintptr(i))) = mod(k1+j, n)
	foundk:
	}
	j = *(*int)(unsafe.Add(unsafe.Pointer(pivk), unsafe.Sizeof(int(0))*uintptr(n-1)))
	pp.Lon[n-1] = j
	for i = n - 2; i >= 0; i-- {
		if cyclic(i+1, *(*int)(unsafe.Add(unsafe.Pointer(pivk), unsafe.Sizeof(int(0))*uintptr(i))), j) != 0 {
			j = *(*int)(unsafe.Add(unsafe.Pointer(pivk), unsafe.Sizeof(int(0))*uintptr(i)))
		}
		pp.Lon[i] = j
	}
	for i = n - 1; cyclic(mod(i+1, n), j, pp.Lon[i]) != 0; i-- {
		pp.Lon[i] = j
	}

	return 0
calloc_error:

	return 1
}
func penalty3(pp *potrace_privpath_s, i int, j int) float64 {
	var (
		n    int     = pp.Len
		pt   *Point  = &pp.Pt[0]
		sums *sums_s = &pp.Sums[0]
		x    float64
		y    float64
		x2   float64
		xy   float64
		y2   float64
		k    float64
		a    float64
		b    float64
		c    float64
		s    float64
		px   float64
		py   float64
		ex   float64
		ey   float64
		r    int = 0
	)
	if j >= n {
		j -= n
		r = 1
	}
	if r == 0 {
		x = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).X - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).X
		y = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).Y - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).Y
		x2 = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).X2 - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).X2
		xy = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).Xy - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).Xy
		y2 = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).Y2 - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).Y2
		k = float64(j + 1 - i)
	} else {
		x = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).X - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).X + (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(n)))).X
		y = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).Y - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).Y + (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(n)))).Y
		x2 = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).X2 - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).X2 + (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(n)))).X2
		xy = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).Xy - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).Xy + (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(n)))).Xy
		y2 = (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(j+1)))).Y2 - (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(i)))).Y2 + (*(*sums_s)(unsafe.Add(unsafe.Pointer(sums), unsafe.Sizeof(sums_s{})*uintptr(n)))).Y2
		k = float64(j + 1 - i + n)
	}
	px = float64((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).X+(*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(j)))).X)/2.0 - float64((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*0))).X)
	py = float64((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).Y+(*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(j)))).Y)/2.0 - float64((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*0))).Y)
	ey = float64((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(j)))).X - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).X)
	ex = float64(-((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(j)))).Y - (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).Y))
	a = (x2-x*2*px)/k + px*px
	b = (xy-x*py-y*px)/k + px*py
	c = (y2-y*2*py)/k + py*py
	s = ex*ex*a + ex*2*ey*b + ey*ey*c
	return math.Sqrt(s)
}
func bestpolygon(pp *potrace_privpath_s) int {
	var (
		i       int
		j       int
		m       int
		k       int
		n       int      = pp.Len
		pen     *float64 = nil
		prev    *int     = nil
		clip0   *int     = nil
		clip1   *int     = nil
		seg0    *int     = nil
		seg1    *int     = nil
		thispen float64
		best    float64
		c       int
	)
	if (func() *float64 {
		pen = &make([]float64, n+1)[0]
		return pen
	}()) == nil {
		goto calloc_error
	}
	if (func() *int {
		prev = &make([]int, n+1)[0]
		return prev
	}()) == nil {
		goto calloc_error
	}
	if (func() *int {
		clip0 = &make([]int, n)[0]
		return clip0
	}()) == nil {
		goto calloc_error
	}
	if (func() *int {
		clip1 = &make([]int, n+1)[0]
		return clip1
	}()) == nil {
		goto calloc_error
	}
	if (func() *int {
		seg0 = &make([]int, n+1)[0]
		return seg0
	}()) == nil {
		goto calloc_error
	}
	if (func() *int {
		seg1 = &make([]int, n+1)[0]
		return seg1
	}()) == nil {
		goto calloc_error
	}
	for i = 0; i < n; i++ {
		c = mod(pp.Lon[mod(i-1, n)]-1, n)
		if c == i {
			c = mod(i+1, n)
		}
		if c < i {
			*(*int)(unsafe.Add(unsafe.Pointer(clip0), unsafe.Sizeof(int(0))*uintptr(i))) = n
		} else {
			*(*int)(unsafe.Add(unsafe.Pointer(clip0), unsafe.Sizeof(int(0))*uintptr(i))) = c
		}
	}
	j = 1
	for i = 0; i < n; i++ {
		for j <= *(*int)(unsafe.Add(unsafe.Pointer(clip0), unsafe.Sizeof(int(0))*uintptr(i))) {
			*(*int)(unsafe.Add(unsafe.Pointer(clip1), unsafe.Sizeof(int(0))*uintptr(j))) = i
			j++
		}
	}
	i = 0
	for j = 0; i < n; j++ {
		*(*int)(unsafe.Add(unsafe.Pointer(seg0), unsafe.Sizeof(int(0))*uintptr(j))) = i
		i = *(*int)(unsafe.Add(unsafe.Pointer(clip0), unsafe.Sizeof(int(0))*uintptr(i)))
	}
	*(*int)(unsafe.Add(unsafe.Pointer(seg0), unsafe.Sizeof(int(0))*uintptr(j))) = n
	m = j
	i = n
	for j = m; j > 0; j-- {
		*(*int)(unsafe.Add(unsafe.Pointer(seg1), unsafe.Sizeof(int(0))*uintptr(j))) = i
		i = *(*int)(unsafe.Add(unsafe.Pointer(clip1), unsafe.Sizeof(int(0))*uintptr(i)))
	}
	*(*int)(unsafe.Add(unsafe.Pointer(seg1), unsafe.Sizeof(int(0))*0)) = 0
	*(*float64)(unsafe.Add(unsafe.Pointer(pen), unsafe.Sizeof(float64(0))*0)) = 0
	for j = 1; j <= m; j++ {
		for i = *(*int)(unsafe.Add(unsafe.Pointer(seg1), unsafe.Sizeof(int(0))*uintptr(j))); i <= *(*int)(unsafe.Add(unsafe.Pointer(seg0), unsafe.Sizeof(int(0))*uintptr(j))); i++ {
			best = float64(-1)
			for k = *(*int)(unsafe.Add(unsafe.Pointer(seg0), unsafe.Sizeof(int(0))*uintptr(j-1))); k >= *(*int)(unsafe.Add(unsafe.Pointer(clip1), unsafe.Sizeof(int(0))*uintptr(i))); k-- {
				thispen = penalty3(pp, k, i) + *(*float64)(unsafe.Add(unsafe.Pointer(pen), unsafe.Sizeof(float64(0))*uintptr(k)))
				if best < 0 || thispen < best {
					*(*int)(unsafe.Add(unsafe.Pointer(prev), unsafe.Sizeof(int(0))*uintptr(i))) = k
					best = thispen
				}
			}
			*(*float64)(unsafe.Add(unsafe.Pointer(pen), unsafe.Sizeof(float64(0))*uintptr(i))) = best
		}
	}
	pp.M = m
	if (func() []int {
		p := &pp.Po
		pp.Po = make([]int, m)
		return *p
	}()) == nil {
		goto calloc_error
	}
	for func() int {
		i = n
		return func() int {
			j = m - 1
			return j
		}()
	}(); i > 0; j-- {
		i = *(*int)(unsafe.Add(unsafe.Pointer(prev), unsafe.Sizeof(int(0))*uintptr(i)))
		pp.Po[j] = i
	}

	return 0
calloc_error:

	return 1
}
func adjust_vertices(pp *potrace_privpath_s) int {
	var (
		m   int         = pp.M
		po  *int        = &pp.Po[0]
		n   int         = pp.Len
		pt  *Point      = &pp.Pt[0]
		x0  int         = pp.X0
		y0  int         = pp.Y0
		ctr *DPoint     = nil
		dir *DPoint     = nil
		q   *quadform_t = nil
		v   [3]float64
		d   float64
		i   int
		j   int
		k   int
		l   int
		s   DPoint
		r   int
	)
	if (func() *DPoint {
		ctr = &make([]DPoint, m)[0]
		return ctr
	}()) == nil {
		goto calloc_error
	}
	if (func() *DPoint {
		dir = &make([]DPoint, m)[0]
		return dir
	}()) == nil {
		goto calloc_error
	}
	if (func() *quadform_t {
		q = (*quadform_t)(unsafe.Pointer(&make([]quadform_t, m)[0][0][0]))
		return q
	}()) == nil {
		goto calloc_error
	}
	r = privcurve_init(&pp.Curve, m)
	if r != 0 {
		goto calloc_error
	}
	for i = 0; i < m; i++ {
		j = *(*int)(unsafe.Add(unsafe.Pointer(po), unsafe.Sizeof(int(0))*uintptr(mod(i+1, m))))
		j = mod(j-*(*int)(unsafe.Add(unsafe.Pointer(po), unsafe.Sizeof(int(0))*uintptr(i))), n) + *(*int)(unsafe.Add(unsafe.Pointer(po), unsafe.Sizeof(int(0))*uintptr(i)))
		pointslope(pp, *(*int)(unsafe.Add(unsafe.Pointer(po), unsafe.Sizeof(int(0))*uintptr(i))), j, (*DPoint)(unsafe.Add(unsafe.Pointer(ctr), unsafe.Sizeof(DPoint{})*uintptr(i))), (*DPoint)(unsafe.Add(unsafe.Pointer(dir), unsafe.Sizeof(DPoint{})*uintptr(i))))
	}
	for i = 0; i < m; i++ {
		d = ((*(*DPoint)(unsafe.Add(unsafe.Pointer(dir), unsafe.Sizeof(DPoint{})*uintptr(i)))).X * (*(*DPoint)(unsafe.Add(unsafe.Pointer(dir), unsafe.Sizeof(DPoint{})*uintptr(i)))).X) + (*(*DPoint)(unsafe.Add(unsafe.Pointer(dir), unsafe.Sizeof(DPoint{})*uintptr(i)))).Y*(*(*DPoint)(unsafe.Add(unsafe.Pointer(dir), unsafe.Sizeof(DPoint{})*uintptr(i)))).Y
		if d == 0.0 {
			for j = 0; j < 3; j++ {
				for k = 0; k < 3; k++ {
					(*(*quadform_t)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(quadform_t{})*uintptr(i))))[j][k] = 0
				}
			}
		} else {
			v[0] = (*(*DPoint)(unsafe.Add(unsafe.Pointer(dir), unsafe.Sizeof(DPoint{})*uintptr(i)))).Y
			v[1] = -(*(*DPoint)(unsafe.Add(unsafe.Pointer(dir), unsafe.Sizeof(DPoint{})*uintptr(i)))).X
			v[2] = -v[1]*(*(*DPoint)(unsafe.Add(unsafe.Pointer(ctr), unsafe.Sizeof(DPoint{})*uintptr(i)))).Y - v[0]*(*(*DPoint)(unsafe.Add(unsafe.Pointer(ctr), unsafe.Sizeof(DPoint{})*uintptr(i)))).X
			for l = 0; l < 3; l++ {
				for k = 0; k < 3; k++ {
					(*(*quadform_t)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(quadform_t{})*uintptr(i))))[l][k] = v[l] * v[k] / d
				}
			}
		}
	}
	for i = 0; i < m; i++ {
		var (
			Q    quadform_t
			w    DPoint
			dx   float64
			dy   float64
			det  float64
			min  float64
			cand float64
			xmin float64
			ymin float64
			z    int
		)
		s.X = float64((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(*(*int)(unsafe.Add(unsafe.Pointer(po), unsafe.Sizeof(int(0))*uintptr(i))))))).X - x0)
		s.Y = float64((*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(*(*int)(unsafe.Add(unsafe.Pointer(po), unsafe.Sizeof(int(0))*uintptr(i))))))).Y - y0)
		j = mod(i-1, m)
		for l = 0; l < 3; l++ {
			for k = 0; k < 3; k++ {
				Q[l][k] = (*(*quadform_t)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(quadform_t{})*uintptr(j))))[l][k] + (*(*quadform_t)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(quadform_t{})*uintptr(i))))[l][k]
			}
		}
		for {
			det = Q[0][0]*Q[1][1] - Q[0][1]*Q[1][0]
			if det != 0.0 {
				w.X = (-Q[0][2]*Q[1][1] + Q[1][2]*Q[0][1]) / det
				w.Y = (Q[0][2]*Q[1][0] - Q[1][2]*Q[0][0]) / det
				break
			}
			if Q[0][0] > Q[1][1] {
				v[0] = -Q[0][1]
				v[1] = Q[0][0]
			} else if Q[1][1] != 0 {
				v[0] = -Q[1][1]
				v[1] = Q[1][0]
			} else {
				v[0] = 1
				v[1] = 0
			}
			d = ((v[0]) * (v[0])) + (v[1])*(v[1])
			v[2] = -v[1]*s.Y - v[0]*s.X
			for l = 0; l < 3; l++ {
				for k = 0; k < 3; k++ {
					Q[l][k] += v[l] * v[k] / d
				}
			}
		}
		dx = math.Abs(w.X - s.X)
		dy = math.Abs(w.Y - s.Y)
		if dx <= 0.5 && dy <= 0.5 {
			pp.Curve.Vertex[i].X = w.X + float64(x0)
			pp.Curve.Vertex[i].Y = w.Y + float64(y0)
			continue
		}
		min = quadform(Q, s)
		xmin = s.X
		ymin = s.Y
		if Q[0][0] == 0.0 {
			goto fixx
		}
		for z = 0; z < 2; z++ {
			w.Y = s.Y - 0.5 + float64(z)
			w.X = -(Q[0][1]*w.Y + Q[0][2]) / Q[0][0]
			dx = math.Abs(w.X - s.X)
			cand = quadform(Q, w)
			if dx <= 0.5 && cand < min {
				min = cand
				xmin = w.X
				ymin = w.Y
			}
		}
	fixx:
		if Q[1][1] == 0.0 {
			goto corners
		}
		for z = 0; z < 2; z++ {
			w.X = s.X - 0.5 + float64(z)
			w.Y = -(Q[1][0]*w.X + Q[1][2]) / Q[1][1]
			dy = math.Abs(w.Y - s.Y)
			cand = quadform(Q, w)
			if dy <= 0.5 && cand < min {
				min = cand
				xmin = w.X
				ymin = w.Y
			}
		}
	corners:
		for l = 0; l < 2; l++ {
			for k = 0; k < 2; k++ {
				w.X = s.X - 0.5 + float64(l)
				w.Y = s.Y - 0.5 + float64(k)
				cand = quadform(Q, w)
				if cand < min {
					min = cand
					xmin = w.X
					ymin = w.Y
				}
			}
		}
		pp.Curve.Vertex[i].X = xmin + float64(x0)
		pp.Curve.Vertex[i].Y = ymin + float64(y0)
		continue
	}

	return 0
calloc_error:

	return 1
}
func reverse(curve *privCurve) {
	var (
		m   int = curve.N
		i   int
		j   int
		tmp DPoint
	)
	for func() int {
		i = 0
		return func() int {
			j = m - 1
			return j
		}()
	}(); i < j; func() int {
		i++
		return func() int {
			p := &j
			x := *p
			*p--
			return x
		}()
	}() {
		tmp = curve.Vertex[i]
		curve.Vertex[i] = curve.Vertex[j]
		curve.Vertex[j] = tmp
	}
}
func smooth(curve *privCurve, alphamax float64) {
	var (
		m     int = curve.N
		i     int
		j     int
		k     int
		dd    float64
		denom float64
		alpha float64
		p2    DPoint
		p3    DPoint
		p4    DPoint
	)
	for i = 0; i < m; i++ {
		j = mod(i+1, m)
		k = mod(i+2, m)
		p4 = aux_interval(1/2.0, curve.Vertex[k], curve.Vertex[j])
		denom = ddenom(curve.Vertex[i], curve.Vertex[k])
		if denom != 0.0 {
			dd = dpara(curve.Vertex[i], curve.Vertex[j], curve.Vertex[k]) / denom
			dd = math.Abs(dd)
			if dd > 1 {
				alpha = 1 - 1.0/dd
			} else {
				alpha = 0
			}
			alpha = alpha / 0.75
		} else {
			alpha = 4 / 3.0
		}
		curve.Alpha0[j] = alpha
		if alpha >= alphamax {
			curve.Tag[j] = POTRACE_CORNER
			curve.C[j][1] = curve.Vertex[j]
			curve.C[j][2] = p4
		} else {
			if alpha < 0.55 {
				alpha = 0.55
			} else if alpha > 1 {
				alpha = 1
			}
			p2 = aux_interval(alpha*0.5+0.5, curve.Vertex[i], curve.Vertex[j])
			p3 = aux_interval(alpha*0.5+0.5, curve.Vertex[k], curve.Vertex[j])
			curve.Tag[j] = POTRACE_CURVETO
			curve.C[j][0] = p2
			curve.C[j][1] = p3
			curve.C[j][2] = p4
		}
		curve.Alpha[j] = alpha
		curve.Beta[j] = 0.5
	}
	curve.Alphacurve = 1
	return
}

type opti_s struct {
	Pen   float64
	C     [2]DPoint
	T     float64
	S     float64
	Alpha float64
}
type opti_t opti_s

func opti_penalty(pp *potrace_privpath_s, i int, j int, res *opti_t, opttolerance float64, convc *int, areac *float64) int {
	var (
		m     int = pp.Curve.N
		k     int
		k1    int
		k2    int
		conv  int
		i1    int
		area  float64
		alpha float64
		d     float64
		d1    float64
		d2    float64
		p0    DPoint
		p1    DPoint
		p2    DPoint
		p3    DPoint
		pt    DPoint
		A     float64
		R     float64
		A1    float64
		A2    float64
		A3    float64
		A4    float64
		s     float64
		t     float64
	)
	if i == j {
		return 1
	}
	k = i
	i1 = mod(i+1, m)
	k1 = mod(k+1, m)
	conv = *(*int)(unsafe.Add(unsafe.Pointer(convc), unsafe.Sizeof(int(0))*uintptr(k1)))
	if conv == 0 {
		return 1
	}
	d = ddist(pp.Curve.Vertex[i], pp.Curve.Vertex[i1])
	for k = k1; k != j; k = k1 {
		k1 = mod(k+1, m)
		k2 = mod(k+2, m)
		if *(*int)(unsafe.Add(unsafe.Pointer(convc), unsafe.Sizeof(int(0))*uintptr(k1))) != conv {
			return 1
		}
		if (func() int {
			if cprod(pp.Curve.Vertex[i], pp.Curve.Vertex[i1], pp.Curve.Vertex[k1], pp.Curve.Vertex[k2]) > 0 {
				return 1
			}
			if cprod(pp.Curve.Vertex[i], pp.Curve.Vertex[i1], pp.Curve.Vertex[k1], pp.Curve.Vertex[k2]) < 0 {
				return -1
			}
			return 0
		}()) != conv {
			return 1
		}
		if iprod1(pp.Curve.Vertex[i], pp.Curve.Vertex[i1], pp.Curve.Vertex[k1], pp.Curve.Vertex[k2]) < d*ddist(pp.Curve.Vertex[k1], pp.Curve.Vertex[k2])*(-0.999847695156) {
			return 1
		}
	}
	p0 = pp.Curve.C[mod(i, m)][2]
	p1 = pp.Curve.Vertex[mod(i+1, m)]
	p2 = pp.Curve.Vertex[mod(j, m)]
	p3 = pp.Curve.C[mod(j, m)][2]
	area = *(*float64)(unsafe.Add(unsafe.Pointer(areac), unsafe.Sizeof(float64(0))*uintptr(j))) - *(*float64)(unsafe.Add(unsafe.Pointer(areac), unsafe.Sizeof(float64(0))*uintptr(i)))
	area -= dpara(pp.Curve.Vertex[0], pp.Curve.C[i][2], pp.Curve.C[j][2]) / 2
	if i >= j {
		area += *(*float64)(unsafe.Add(unsafe.Pointer(areac), unsafe.Sizeof(float64(0))*uintptr(m)))
	}
	A1 = dpara(p0, p1, p2)
	A2 = dpara(p0, p1, p3)
	A3 = dpara(p0, p2, p3)
	A4 = A1 + A3 - A2
	if A2 == A1 {
		return 1
	}
	t = A3 / (A3 - A4)
	s = A2 / (A2 - A1)
	A = A2 * t / 2.0
	if A == 0.0 {
		return 1
	}
	R = area / A
	alpha = 2 - math.Sqrt(4-R/0.3)
	res.C[0] = aux_interval(t*alpha, p0, p1)
	res.C[1] = aux_interval(s*alpha, p3, p2)
	res.Alpha = alpha
	res.T = t
	res.S = s
	p1 = res.C[0]
	p2 = res.C[1]
	res.Pen = 0
	for k = mod(i+1, m); k != j; k = k1 {
		k1 = mod(k+1, m)
		t = tangent(p0, p1, p2, p3, pp.Curve.Vertex[k], pp.Curve.Vertex[k1])
		if t < -0.5 {
			return 1
		}
		pt = trace_bezier(t, p0, p1, p2, p3)
		d = ddist(pp.Curve.Vertex[k], pp.Curve.Vertex[k1])
		if d == 0.0 {
			return 1
		}
		d1 = dpara(pp.Curve.Vertex[k], pp.Curve.Vertex[k1], pt) / d
		if math.Abs(d1) > opttolerance {
			return 1
		}
		if trace_iprod(pp.Curve.Vertex[k], pp.Curve.Vertex[k1], pt) < 0 || trace_iprod(pp.Curve.Vertex[k1], pp.Curve.Vertex[k], pt) < 0 {
			return 1
		}
		res.Pen += d1 * d1
	}
	for k = i; k != j; k = k1 {
		k1 = mod(k+1, m)
		t = tangent(p0, p1, p2, p3, pp.Curve.C[k][2], pp.Curve.C[k1][2])
		if t < -0.5 {
			return 1
		}
		pt = trace_bezier(t, p0, p1, p2, p3)
		d = ddist(pp.Curve.C[k][2], pp.Curve.C[k1][2])
		if d == 0.0 {
			return 1
		}
		d1 = dpara(pp.Curve.C[k][2], pp.Curve.C[k1][2], pt) / d
		d2 = dpara(pp.Curve.C[k][2], pp.Curve.C[k1][2], pp.Curve.Vertex[k1]) / d
		d2 *= pp.Curve.Alpha[k1] * 0.75
		if d2 < 0 {
			d1 = -d1
			d2 = -d2
		}
		if d1 < d2-opttolerance {
			return 1
		}
		if d1 < d2 {
			res.Pen += (d1 - d2) * (d1 - d2)
		}
	}
	return 0
}
func opticurve(pp *potrace_privpath_s, opttolerance float64) int {
	var (
		m     int      = pp.Curve.N
		pt    *int     = nil
		pen   *float64 = nil
		len_  *int     = nil
		opt   *opti_t  = nil
		om    int
		i     int
		j     int
		r     int
		o     opti_t
		p0    DPoint
		i1    int
		area  float64
		alpha float64
		s     *float64 = nil
		t     *float64 = nil
		convc *int     = nil
		areac *float64 = nil
	)
	if (func() *int {
		pt = &make([]int, m+1)[0]
		return pt
	}()) == nil {
		goto calloc_error
	}
	if (func() *float64 {
		pen = &make([]float64, m+1)[0]
		return pen
	}()) == nil {
		goto calloc_error
	}
	if (func() *int {
		len_ = &make([]int, m+1)[0]
		return len_
	}()) == nil {
		goto calloc_error
	}
	if (func() *opti_t {
		opt = &make([]opti_t, m+1)[0]
		return opt
	}()) == nil {
		goto calloc_error
	}
	if (func() *int {
		convc = &make([]int, m)[0]
		return convc
	}()) == nil {
		goto calloc_error
	}
	if (func() *float64 {
		areac = &make([]float64, m+1)[0]
		return areac
	}()) == nil {
		goto calloc_error
	}
	for i = 0; i < m; i++ {
		if pp.Curve.Tag[i] == POTRACE_CURVETO {
			if dpara(pp.Curve.Vertex[mod(i-1, m)], pp.Curve.Vertex[i], pp.Curve.Vertex[mod(i+1, m)]) > 0 {
				*(*int)(unsafe.Add(unsafe.Pointer(convc), unsafe.Sizeof(int(0))*uintptr(i))) = 1
			} else if dpara(pp.Curve.Vertex[mod(i-1, m)], pp.Curve.Vertex[i], pp.Curve.Vertex[mod(i+1, m)]) < 0 {
				*(*int)(unsafe.Add(unsafe.Pointer(convc), unsafe.Sizeof(int(0))*uintptr(i))) = -1
			} else {
				*(*int)(unsafe.Add(unsafe.Pointer(convc), unsafe.Sizeof(int(0))*uintptr(i))) = 0
			}
		} else {
			*(*int)(unsafe.Add(unsafe.Pointer(convc), unsafe.Sizeof(int(0))*uintptr(i))) = 0
		}
	}
	area = 0.0
	*(*float64)(unsafe.Add(unsafe.Pointer(areac), unsafe.Sizeof(float64(0))*0)) = 0.0
	p0 = pp.Curve.Vertex[0]
	for i = 0; i < m; i++ {
		i1 = mod(i+1, m)
		if pp.Curve.Tag[i1] == POTRACE_CURVETO {
			alpha = pp.Curve.Alpha[i1]
			area += alpha * 0.3 * (4 - alpha) * dpara(pp.Curve.C[i][2], pp.Curve.Vertex[i1], pp.Curve.C[i1][2]) / 2
			area += dpara(p0, pp.Curve.C[i][2], pp.Curve.C[i1][2]) / 2
		}
		*(*float64)(unsafe.Add(unsafe.Pointer(areac), unsafe.Sizeof(float64(0))*uintptr(i+1))) = area
	}
	*(*int)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(int(0))*0)) = -1
	*(*float64)(unsafe.Add(unsafe.Pointer(pen), unsafe.Sizeof(float64(0))*0)) = 0
	*(*int)(unsafe.Add(unsafe.Pointer(len_), unsafe.Sizeof(int(0))*0)) = 0
	for j = 1; j <= m; j++ {
		*(*int)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(int(0))*uintptr(j))) = j - 1
		*(*float64)(unsafe.Add(unsafe.Pointer(pen), unsafe.Sizeof(float64(0))*uintptr(j))) = *(*float64)(unsafe.Add(unsafe.Pointer(pen), unsafe.Sizeof(float64(0))*uintptr(j-1)))
		*(*int)(unsafe.Add(unsafe.Pointer(len_), unsafe.Sizeof(int(0))*uintptr(j))) = *(*int)(unsafe.Add(unsafe.Pointer(len_), unsafe.Sizeof(int(0))*uintptr(j-1))) + 1
		for i = j - 2; i >= 0; i-- {
			r = opti_penalty(pp, i, mod(j, m), &o, opttolerance, convc, areac)
			if r != 0 {
				break
			}
			if *(*int)(unsafe.Add(unsafe.Pointer(len_), unsafe.Sizeof(int(0))*uintptr(j))) > *(*int)(unsafe.Add(unsafe.Pointer(len_), unsafe.Sizeof(int(0))*uintptr(i)))+1 || *(*int)(unsafe.Add(unsafe.Pointer(len_), unsafe.Sizeof(int(0))*uintptr(j))) == *(*int)(unsafe.Add(unsafe.Pointer(len_), unsafe.Sizeof(int(0))*uintptr(i)))+1 && *(*float64)(unsafe.Add(unsafe.Pointer(pen), unsafe.Sizeof(float64(0))*uintptr(j))) > *(*float64)(unsafe.Add(unsafe.Pointer(pen), unsafe.Sizeof(float64(0))*uintptr(i)))+o.Pen {
				*(*int)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(int(0))*uintptr(j))) = i
				*(*float64)(unsafe.Add(unsafe.Pointer(pen), unsafe.Sizeof(float64(0))*uintptr(j))) = *(*float64)(unsafe.Add(unsafe.Pointer(pen), unsafe.Sizeof(float64(0))*uintptr(i))) + o.Pen
				*(*int)(unsafe.Add(unsafe.Pointer(len_), unsafe.Sizeof(int(0))*uintptr(j))) = *(*int)(unsafe.Add(unsafe.Pointer(len_), unsafe.Sizeof(int(0))*uintptr(i))) + 1
				*(*opti_t)(unsafe.Add(unsafe.Pointer(opt), unsafe.Sizeof(opti_t{})*uintptr(j))) = o
			}
		}
	}
	om = *(*int)(unsafe.Add(unsafe.Pointer(len_), unsafe.Sizeof(int(0))*uintptr(m)))
	r = privcurve_init(&pp.Ocurve, om)
	if r != 0 {
		goto calloc_error
	}
	if (func() *float64 {
		s = &make([]float64, om)[0]
		return s
	}()) == nil {
		goto calloc_error
	}
	if (func() *float64 {
		t = &make([]float64, om)[0]
		return t
	}()) == nil {
		goto calloc_error
	}
	j = m
	for i = om - 1; i >= 0; i-- {
		if *(*int)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(int(0))*uintptr(j))) == j-1 {
			pp.Ocurve.Tag[i] = pp.Curve.Tag[mod(j, m)]
			pp.Ocurve.C[i][0] = pp.Curve.C[mod(j, m)][0]
			pp.Ocurve.C[i][1] = pp.Curve.C[mod(j, m)][1]
			pp.Ocurve.C[i][2] = pp.Curve.C[mod(j, m)][2]
			pp.Ocurve.Vertex[i] = pp.Curve.Vertex[mod(j, m)]
			pp.Ocurve.Alpha[i] = pp.Curve.Alpha[mod(j, m)]
			pp.Ocurve.Alpha0[i] = pp.Curve.Alpha0[mod(j, m)]
			pp.Ocurve.Beta[i] = pp.Curve.Beta[mod(j, m)]
			*(*float64)(unsafe.Add(unsafe.Pointer(s), unsafe.Sizeof(float64(0))*uintptr(i))) = func() float64 {
				p := (*float64)(unsafe.Add(unsafe.Pointer(t), unsafe.Sizeof(float64(0))*uintptr(i)))
				*(*float64)(unsafe.Add(unsafe.Pointer(t), unsafe.Sizeof(float64(0))*uintptr(i))) = 1.0
				return *p
			}()
		} else {
			pp.Ocurve.Tag[i] = POTRACE_CURVETO
			pp.Ocurve.C[i][0] = (*(*opti_t)(unsafe.Add(unsafe.Pointer(opt), unsafe.Sizeof(opti_t{})*uintptr(j)))).C[0]
			pp.Ocurve.C[i][1] = (*(*opti_t)(unsafe.Add(unsafe.Pointer(opt), unsafe.Sizeof(opti_t{})*uintptr(j)))).C[1]
			pp.Ocurve.C[i][2] = pp.Curve.C[mod(j, m)][2]
			pp.Ocurve.Vertex[i] = aux_interval((*(*opti_t)(unsafe.Add(unsafe.Pointer(opt), unsafe.Sizeof(opti_t{})*uintptr(j)))).S, pp.Curve.C[mod(j, m)][2], pp.Curve.Vertex[mod(j, m)])
			pp.Ocurve.Alpha[i] = (*(*opti_t)(unsafe.Add(unsafe.Pointer(opt), unsafe.Sizeof(opti_t{})*uintptr(j)))).Alpha
			pp.Ocurve.Alpha0[i] = (*(*opti_t)(unsafe.Add(unsafe.Pointer(opt), unsafe.Sizeof(opti_t{})*uintptr(j)))).Alpha
			*(*float64)(unsafe.Add(unsafe.Pointer(s), unsafe.Sizeof(float64(0))*uintptr(i))) = (*(*opti_t)(unsafe.Add(unsafe.Pointer(opt), unsafe.Sizeof(opti_t{})*uintptr(j)))).S
			*(*float64)(unsafe.Add(unsafe.Pointer(t), unsafe.Sizeof(float64(0))*uintptr(i))) = (*(*opti_t)(unsafe.Add(unsafe.Pointer(opt), unsafe.Sizeof(opti_t{})*uintptr(j)))).T
		}
		j = *(*int)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(int(0))*uintptr(j)))
	}
	for i = 0; i < om; i++ {
		i1 = mod(i+1, om)
		pp.Ocurve.Beta[i] = *(*float64)(unsafe.Add(unsafe.Pointer(s), unsafe.Sizeof(float64(0))*uintptr(i))) / (*(*float64)(unsafe.Add(unsafe.Pointer(s), unsafe.Sizeof(float64(0))*uintptr(i))) + *(*float64)(unsafe.Add(unsafe.Pointer(t), unsafe.Sizeof(float64(0))*uintptr(i1))))
	}
	pp.Ocurve.Alphacurve = 1

	return 0
calloc_error:

	return 1
}
func process_path(plist *Path, param *Config, progress *progress) int {
	var (
		p  *Path
		nn float64 = 0
		cn float64 = 0
	)
	if progress.Callback != nil {
		nn = 0
		for p = plist; p != nil; p = p.Next {
			nn += float64(p.Priv.Len)
		}
		cn = 0
	}
	for p = plist; p != nil; p = p.Next {
		if calc_sums(p.Priv) != 0 {
			goto try_error
		}
		if calc_lon(p.Priv) != 0 {
			goto try_error
		}
		if bestpolygon(p.Priv) != 0 {
			goto try_error
		}
		if adjust_vertices(p.Priv) != 0 {
			goto try_error
		}
		if p.Sign == '-' {
			reverse(&p.Priv.Curve)
		}
		smooth(&p.Priv.Curve, param.AlphaMax)
		if param.OptiCurve {
			if opticurve(p.Priv, param.OptTolerance) != 0 {
				goto try_error
			}
			p.Priv.Fcurve = &p.Priv.Ocurve
		} else {
			p.Priv.Fcurve = &p.Priv.Curve
		}
		privcurve_to_curve(p.Priv.Fcurve, &p.Curve)
		if progress.Callback != nil {
			cn += float64(p.Priv.Len)
			progress_update(cn/nn, progress)
		}
	}
	progress_update(1.0, progress)
	return 0
try_error:
	return 1
}
