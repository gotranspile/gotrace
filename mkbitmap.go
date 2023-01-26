package gotrace

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"os"
	"unsafe"
)

const _POSIX_C_SOURCE = 200809
const _NETBSD_SOURCE = 1

type Config struct {
	Outfile     *byte
	Infiles     **byte
	Infilecount int
	Invert      bool
	Highpass    bool
	Lambda      float64
	Lowpass     bool
	Lambda1     float64
	Scale       int
	Linear      bool
	Bilevel     bool
	Level       float64
	Outext      *byte
}

func lowpass(gm *Greymap, lambda float64) {
	var (
		f float64
		g float64
		c float64
		d float64
		B float64
		x int
		y int
	)
	if gm.H == 0 || gm.W == 0 {
		return
	}
	B = 2/(lambda*lambda) + 1
	c = B - math.Sqrt(B*B-1)
	d = 1 - c
	for y = 0; y < gm.H; y++ {
		f = func() float64 {
			g = 0
			return g
		}()
		for x = 0; x < gm.W; x++ {
			f = f*c + float64(gm.Map[int64(y)*int64(gm.Dy)+int64(x)])*d
			g = g*c + f*d
			gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(g)
		}
		for x = gm.W - 1; x >= 0; x-- {
			f = f*c + float64(gm.Map[int64(y)*int64(gm.Dy)+int64(x)])*d
			g = g*c + f*d
			gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(g)
		}
		for x = 0; x < gm.W; x++ {
			f = f * c
			g = g*c + f*d
			if f+g < 1/255.0 {
				break
			}
			gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(float64(gm.Map[int64(y)*int64(gm.Dy)+int64(x)]) + g)
		}
	}
	for x = 0; x < gm.W; x++ {
		f = func() float64 {
			g = 0
			return g
		}()
		for y = 0; y < gm.H; y++ {
			f = f*c + float64(gm.Map[int64(y)*int64(gm.Dy)+int64(x)])*d
			g = g*c + f*d
			gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(g)
		}
		for y = gm.H - 1; y >= 0; y-- {
			f = f*c + float64(gm.Map[int64(y)*int64(gm.Dy)+int64(x)])*d
			g = g*c + f*d
			gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(g)
		}
		for y = 0; y < gm.H; y++ {
			f = f * c
			g = g*c + f*d
			if f+g < 1/255.0 {
				break
			}
			gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(float64(gm.Map[int64(y)*int64(gm.Dy)+int64(x)]) + g)
		}
	}
}
func highpass(gm *Greymap, lambda float64) int {
	var (
		gm1 *Greymap
		f   float64
		x   int
		y   int
	)
	if gm.H == 0 || gm.W == 0 {
		return 0
	}
	gm1 = gm_dup(gm)
	if gm1 == nil {
		return 1
	}
	lowpass(gm1, lambda)
	for y = 0; y < gm.H; y++ {
		for x = 0; x < gm.W; x++ {
			f = float64(gm.Map[int64(y)*int64(gm.Dy)+int64(x)])
			f -= float64(gm1.Map[int64(y)*int64(gm1.Dy)+int64(x)])
			f += 128
			gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(f)
		}
	}
	gm_free(gm1)
	return 0
}
func threshold(gm *Greymap, c float64) *Bitmap {
	var (
		w      int
		h      int
		bm_out *Bitmap = nil
		c1     float64
		x      int
		y      int
		p      float64
	)
	w = gm.W
	h = gm.H
	bm_out = NewBitmap(w, h)
	if bm_out == nil {
		return nil
	}
	c1 = c * math.MaxUint8
	for y = 0; y < h; y++ {
		for x = 0; x < w; x++ {
			p = float64(gm.Map[int64(y)*int64(gm.Dy)+int64(x)])
			if p < c1 {
				bm_out.Map[int64(y)*int64(bm_out.Dy)+int64(x/(8*(int(unsafe.Sizeof(Word(0))))))] |= Word((1 << ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)) >> (x & ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)))
			} else {
				bm_out.Map[int64(y)*int64(bm_out.Dy)+int64(x/(8*(int(unsafe.Sizeof(Word(0))))))] &= Word(^((1 << ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)) >> (x & ((8 * (int(unsafe.Sizeof(Word(0))))) - 1))))
			}
		}
	}
	return bm_out
}
func interpolate_linear(gm *Greymap, s int, bilevel int, c float64) unsafe.Pointer {
	var (
		p00    int
		p01    int
		p10    int
		p11    int
		i      int
		j      int
		x      int
		y      int
		xx     float64
		yy     float64
		av     float64
		c1     float64 = 0
		w      int
		h      int
		p0     float64
		p1     float64
		gm_out *Greymap = nil
		bm_out *Bitmap  = nil
	)
	w = gm.W
	h = gm.H
	if bilevel != 0 {
		bm_out = NewBitmap(w*s, h*s)
		if bm_out == nil {
			return nil
		}
		c1 = c * math.MaxUint8
	} else {
		gm_out = NewGreymap(w*s, h*s)
		if gm_out == nil {
			return nil
		}
	}
	for i = 0; i < w; i++ {
		for j = 0; j < h; j++ {
			if gm.W == 0 || gm.H == 0 {
				p00 = 0
			} else {
				p00 = int(gm.Map[int64(func() int {
					if j < 0 {
						return 0
					}
					if j >= gm.H {
						return gm.H - 1
					}
					return j
				}())*int64(gm.Dy)+int64(func() int {
					if i < 0 {
						return 0
					}
					if i >= gm.W {
						return gm.W - 1
					}
					return i
				}())])
			}
			if gm.W == 0 || gm.H == 0 {
				p01 = 0
			} else {
				p01 = int(gm.Map[int64(func() int {
					if (j + 1) < 0 {
						return 0
					}
					if (j + 1) >= gm.H {
						return gm.H - 1
					}
					return j + 1
				}())*int64(gm.Dy)+int64(func() int {
					if i < 0 {
						return 0
					}
					if i >= gm.W {
						return gm.W - 1
					}
					return i
				}())])
			}
			if gm.W == 0 || gm.H == 0 {
				p10 = 0
			} else {
				p10 = int(gm.Map[int64(func() int {
					if j < 0 {
						return 0
					}
					if j >= gm.H {
						return gm.H - 1
					}
					return j
				}())*int64(gm.Dy)+int64(func() int {
					if (i + 1) < 0 {
						return 0
					}
					if (i + 1) >= gm.W {
						return gm.W - 1
					}
					return i + 1
				}())])
			}
			if gm.W == 0 || gm.H == 0 {
				p11 = 0
			} else {
				p11 = int(gm.Map[int64(func() int {
					if (j + 1) < 0 {
						return 0
					}
					if (j + 1) >= gm.H {
						return gm.H - 1
					}
					return j + 1
				}())*int64(gm.Dy)+int64(func() int {
					if (i + 1) < 0 {
						return 0
					}
					if (i + 1) >= gm.W {
						return gm.W - 1
					}
					return i + 1
				}())])
			}
			if bilevel != 0 {
				if float64(p00) < c1 && float64(p01) < c1 && float64(p10) < c1 && float64(p11) < c1 {
					for x = 0; x < s; x++ {
						for y = 0; y < s; y++ {
							if true {
								bm_out.Map[int64(j*s+y)*int64(bm_out.Dy)+int64((i*s+x)/(8*(int(unsafe.Sizeof(Word(0))))))] |= Word((1 << ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)) >> ((i*s + x) & ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)))
							} else {
								bm_out.Map[int64(j*s+y)*int64(bm_out.Dy)+int64((i*s+x)/(8*(int(unsafe.Sizeof(Word(0))))))] &= Word(^((1 << ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)) >> ((i*s + x) & ((8 * (int(unsafe.Sizeof(Word(0))))) - 1))))
							}
						}
					}
					continue
				}
				if float64(p00) >= c1 && float64(p01) >= c1 && float64(p10) >= c1 && float64(p11) >= c1 {
					continue
				}
			}
			for x = 0; x < s; x++ {
				xx = float64(x) / float64(s)
				p0 = float64(p00)*(1-xx) + float64(p10)*xx
				p1 = float64(p01)*(1-xx) + float64(p11)*xx
				for y = 0; y < s; y++ {
					yy = float64(y) / float64(s)
					av = p0*(1-yy) + p1*yy
					if bilevel != 0 {
						if av < c1 {
							bm_out.Map[int64(j*s+y)*int64(bm_out.Dy)+int64((i*s+x)/(8*(int(unsafe.Sizeof(Word(0))))))] |= Word((1 << ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)) >> ((i*s + x) & ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)))
						} else {
							bm_out.Map[int64(j*s+y)*int64(bm_out.Dy)+int64((i*s+x)/(8*(int(unsafe.Sizeof(Word(0))))))] &= Word(^((1 << ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)) >> ((i*s + x) & ((8 * (int(unsafe.Sizeof(Word(0))))) - 1))))
						}
					} else {
						gm_out.Map[int64(j*s+y)*int64(gm_out.Dy)+int64(i*s+x)] = gm_sample_t(av)
					}
				}
			}
		}
	}
	if bilevel != 0 {
		return unsafe.Pointer(bm_out)
	} else {
		return unsafe.Pointer(gm_out)
	}
}

type double4 [4]float64

func interpolate_cubic(gm *Greymap, s int, bilevel int, c float64) unsafe.Pointer {
	var (
		w      int
		h      int
		poly   *double4 = nil
		p      [4]float64
		window *double4 = nil
		t      float64
		v      float64
		k      int
		l      int
		i      int
		j      int
		x      int
		y      int
		c1     float64  = 0
		gm_out *Greymap = nil
		bm_out *Bitmap  = nil
	)
	if (func() *double4 {
		poly = (*double4)(unsafe.Pointer(&make([]double4, s)[0][0]))
		return poly
	}()) == nil {
		goto calloc_error
	}
	if (func() *double4 {
		window = (*double4)(unsafe.Pointer(&make([]double4, s)[0][0]))
		return window
	}()) == nil {
		goto calloc_error
	}
	w = gm.W
	h = gm.H
	if bilevel != 0 {
		bm_out = NewBitmap(w*s, h*s)
		if bm_out == nil {
			goto calloc_error
		}
		c1 = c * math.MaxUint8
	} else {
		gm_out = NewGreymap(w*s, h*s)
		if gm_out == nil {
			goto calloc_error
		}
	}
	for k = 0; k < s; k++ {
		t = float64(k) / float64(s)
		(*(*double4)(unsafe.Add(unsafe.Pointer(poly), unsafe.Sizeof(double4{})*uintptr(k))))[0] = t * 0.5 * (t - 1) * (1 - t)
		(*(*double4)(unsafe.Add(unsafe.Pointer(poly), unsafe.Sizeof(double4{})*uintptr(k))))[1] = -(t+1)*(t-1)*(1-t) + (t-1)*0.5*(t-2)*t
		(*(*double4)(unsafe.Add(unsafe.Pointer(poly), unsafe.Sizeof(double4{})*uintptr(k))))[2] = (t+1)*0.5*t*(1-t) - t*(t-2)*t
		(*(*double4)(unsafe.Add(unsafe.Pointer(poly), unsafe.Sizeof(double4{})*uintptr(k))))[3] = t * 0.5 * (t - 1) * t
	}
	for y = 0; y < h; y++ {
		x = 0
		for i = 0; i < 4; i++ {
			for j = 0; j < 4; j++ {
				if gm.W == 0 || gm.H == 0 {
					p[j] = 0
				} else {
					p[j] = float64(gm.Map[int64(func() int {
						if (y + j - 1) < 0 {
							return 0
						}
						if (y + j - 1) >= gm.H {
							return gm.H - 1
						}
						return y + j - 1
					}())*int64(gm.Dy)+int64(func() int {
						if (x + i - 1) < 0 {
							return 0
						}
						if (x + i - 1) >= gm.W {
							return gm.W - 1
						}
						return x + i - 1
					}())])
				}
			}
			for k = 0; k < s; k++ {
				(*(*double4)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(double4{})*uintptr(k))))[i] = 0.0
				for j = 0; j < 4; j++ {
					(*(*double4)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(double4{})*uintptr(k))))[i] += (*(*double4)(unsafe.Add(unsafe.Pointer(poly), unsafe.Sizeof(double4{})*uintptr(k))))[j] * p[j]
				}
			}
		}
		for {
			for l = 0; l < s; l++ {
				for k = 0; k < s; k++ {
					v = 0.0
					for i = 0; i < 4; i++ {
						v += (*(*double4)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(double4{})*uintptr(k))))[i] * (*(*double4)(unsafe.Add(unsafe.Pointer(poly), unsafe.Sizeof(double4{})*uintptr(l))))[i]
					}
					if bilevel != 0 {
						if (x*s+l) >= 0 && (x*s+l) < bm_out.W && ((y*s+k) >= 0 && (y*s+k) < bm_out.H) {
							if v < c1 {
								bm_out.Map[int64(y*s+k)*int64(bm_out.Dy)+int64((x*s+l)/(8*(int(unsafe.Sizeof(Word(0))))))] |= Word((1 << ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)) >> ((x*s + l) & ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)))
							} else {
								bm_out.Map[int64(y*s+k)*int64(bm_out.Dy)+int64((x*s+l)/(8*(int(unsafe.Sizeof(Word(0))))))] &= Word(^((1 << ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)) >> ((x*s + l) & ((8 * (int(unsafe.Sizeof(Word(0))))) - 1))))
							}
						} else {
						}
					} else {
						if (x*s+l) >= 0 && (x*s+l) < gm_out.W && (y*s+k) >= 0 && (y*s+k) < gm_out.H {
							gm_out.Map[int64(y*s+k)*int64(gm_out.Dy)+int64(x*s+l)] = gm_sample_t(v)
						} else {
						}
					}
				}
			}
			x++
			if x >= w {
				break
			}
			for i = 0; i < 3; i++ {
				for k = 0; k < s; k++ {
					(*(*double4)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(double4{})*uintptr(k))))[i] = (*(*double4)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(double4{})*uintptr(k))))[i+1]
				}
			}
			i = 3
			for j = 0; j < 4; j++ {
				if gm.W == 0 || gm.H == 0 {
					p[j] = 0
				} else {
					p[j] = float64(gm.Map[int64(func() int {
						if (y + j - 1) < 0 {
							return 0
						}
						if (y + j - 1) >= gm.H {
							return gm.H - 1
						}
						return y + j - 1
					}())*int64(gm.Dy)+int64(func() int {
						if (x + i - 1) < 0 {
							return 0
						}
						if (x + i - 1) >= gm.W {
							return gm.W - 1
						}
						return x + i - 1
					}())])
				}
			}
			for k = 0; k < s; k++ {
				(*(*double4)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(double4{})*uintptr(k))))[i] = 0.0
				for j = 0; j < 4; j++ {
					(*(*double4)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(double4{})*uintptr(k))))[i] += (*(*double4)(unsafe.Add(unsafe.Pointer(poly), unsafe.Sizeof(double4{})*uintptr(k))))[j] * p[j]
				}
			}
		}
	}

	if bilevel != 0 {
		return unsafe.Pointer(bm_out)
	} else {
		return unsafe.Pointer(gm_out)
	}
calloc_error:

	return nil
}
func ProcessFile(info *Config, fin *stdio.File, fout *stdio.File, infile *byte, outfile *byte) {
	var (
		r     int
		gm    *Greymap
		bm    *Bitmap
		sm    unsafe.Pointer
		x     int
		y     int
		count int
	)
	for count = 0; ; count++ {
		r = GreymapRead(fin, &gm)
		switch r {
		case -1:
			stdio.Fprintf(stdio.Stderr(), "potrace: %s: %s\n", infile, libc.StrError(libc.Errno))
			os.Exit(2)
			fallthrough
		case -2:
			stdio.Fprintf(stdio.Stderr(), "potrace: %s: file format error: %s\n", infile, gm_read_error)
			os.Exit(2)
			fallthrough
		case -3:
			if count > 0 {
				return
			}
			stdio.Fprintf(stdio.Stderr(), "potrace: %s: empty file\n", infile)
			os.Exit(2)
			fallthrough
		case -4:
			if count > 0 {
				stdio.Fprintf(stdio.Stderr(), "potrace: %s: warning: junk at end of file\n", infile)
				return
			}
			stdio.Fprintf(stdio.Stderr(), "potrace: %s: file format not recognized\n", infile)
			stdio.Fprintf(stdio.Stderr(), "Possible input file formats are: pnm (pbm, pgm, ppm), bmp.\n")
			os.Exit(2)
			fallthrough
		case 1:
			stdio.Fprintf(stdio.Stderr(), "potrace: %s: warning: premature end of file\n", infile)
		}
		if info.Invert {
			for y = 0; y < gm.H; y++ {
				for x = 0; x < gm.W; x++ {
					gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = math.MaxUint8 - (gm.Map[int64(y)*int64(gm.Dy)+int64(x)])
				}
			}
		}
		if info.Highpass {
			r = highpass(gm, info.Lambda)
			if r != 0 {
				stdio.Fprintf(stdio.Stderr(), "potrace: %s: %s\n", infile, libc.StrError(libc.Errno))
				os.Exit(2)
			}
		}
		if info.Lowpass {
			lowpass(gm, info.Lambda1)
		}
		if info.Scale == 1 && info.Bilevel {
			sm = unsafe.Pointer(threshold(gm, info.Level))
			gm_free(gm)
		} else if info.Scale == 1 {
			sm = unsafe.Pointer(gm)
		} else if info.Linear {
			sm = interpolate_linear(gm, info.Scale, int(libc.BoolToInt(info.Bilevel)), info.Level)
			gm_free(gm)
		} else {
			sm = interpolate_cubic(gm, info.Scale, int(libc.BoolToInt(info.Bilevel)), info.Level)
			gm_free(gm)
		}
		if sm == nil {
			stdio.Fprintf(stdio.Stderr(), "potrace: %s: %s\n", infile, libc.StrError(libc.Errno))
			os.Exit(2)
		}
		if info.Bilevel {
			bm = (*Bitmap)(sm)
			bm_writepbm(fout, bm)
			bm_free(bm)
		} else {
			gm = (*Greymap)(sm)
			gm_writepgm(fout, gm, nil, 1, GM_MODE_POSITIVE, 1.0)
			gm_free(gm)
		}
	}
}
