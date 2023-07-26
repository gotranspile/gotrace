package gotrace

import (
	"math"
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
)

func unit(info *RenderConf, p DPoint) Point {
	var q Point
	q.X = int(math.Floor(p.X*info.Unit + 0.5))
	q.Y = int(math.Floor(p.Y*info.Unit + 0.5))
	return q
}

var cur Point
var lastop int8 = 0
var column int = 0
var newline int = 1

func shiptoken(fout *stdio.File, token *byte) {
	var c int = libc.StrLen(token)
	if newline == 0 && column+c+1 > 75 {
		stdio.Fprintf(fout, "\n")
		column = 0
		newline = 1
	} else if newline == 0 {
		stdio.Fprintf(fout, " ")
		column++
	}
	stdio.Fprintf(fout, "%s", token)
	column += c
	newline = 0
}
func ship(fout *stdio.File, fmt *byte, _rest ...interface{}) {
	var (
		args libc.ArgList
		buf  [4096]byte
		p    *byte
		q    *byte
	)
	args.Start(fmt, _rest)
	stdio.Vsprintf(&buf[0], libc.GoString(fmt), args)
	buf[4095] = 0
	args.End()
	p = &buf[0]
	for (func() *byte {
		q = libc.StrChr(p, ' ')
		return q
	}()) != nil {
		*q = 0
		shiptoken(fout, p)
		p = (*byte)(unsafe.Add(unsafe.Pointer(q), 1))
	}
	shiptoken(fout, p)
}
func svg_moveto(info *RenderConf, fout *stdio.File, p DPoint) {
	cur = unit(info, p)
	ship(fout, libc.CString("M%ld %ld"), cur.X, cur.Y)
	lastop = 'M'
}
func svg_rmoveto(info *RenderConf, fout *stdio.File, p DPoint) {
	var q Point
	q = unit(info, p)
	ship(fout, libc.CString("m%ld %ld"), q.X-cur.X, q.Y-cur.Y)
	cur = q
	lastop = 'm'
}
func svg_lineto(info *RenderConf, fout *stdio.File, p DPoint) {
	var q Point
	q = unit(info, p)
	if int(lastop) != 'l' {
		ship(fout, libc.CString("l%ld %ld"), q.X-cur.X, q.Y-cur.Y)
	} else {
		ship(fout, libc.CString("%ld %ld"), q.X-cur.X, q.Y-cur.Y)
	}
	cur = q
	lastop = 'l'
}
func svg_curveto(info *RenderConf, fout *stdio.File, p1 DPoint, p2 DPoint, p3 DPoint) {
	var (
		q1 Point
		q2 Point
		q3 Point
	)
	q1 = unit(info, p1)
	q2 = unit(info, p2)
	q3 = unit(info, p3)
	if int(lastop) != 'c' {
		ship(fout, libc.CString("c%ld %ld %ld %ld %ld %ld"), q1.X-cur.X, q1.Y-cur.Y, q2.X-cur.X, q2.Y-cur.Y, q3.X-cur.X, q3.Y-cur.Y)
	} else {
		ship(fout, libc.CString("%ld %ld %ld %ld %ld %ld"), q1.X-cur.X, q1.Y-cur.Y, q2.X-cur.X, q2.Y-cur.Y, q3.X-cur.X, q3.Y-cur.Y)
	}
	cur = q3
	lastop = 'c'
}
func svg_path(info *RenderConf, fout *stdio.File, curve *Curve, abs int) int {
	var (
		i int
		c *DPoint
		m int = curve.N
	)
	c = &curve.C[m-1][0]
	if abs != 0 {
		svg_moveto(info, fout, *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*2)))
	} else {
		svg_rmoveto(info, fout, *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*2)))
	}
	for i = 0; i < m; i++ {
		c = &curve.C[i][0]
		switch curve.Tag[i] {
		case POTRACE_CORNER:
			svg_lineto(info, fout, *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*1)))
			svg_lineto(info, fout, *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*2)))
		case POTRACE_CURVETO:
			svg_curveto(info, fout, *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*0)), *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*1)), *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*2)))
		}
	}
	newline = 1
	shiptoken(fout, libc.CString("z"))
	return 0
}
func svg_jaggy_path(info *RenderConf, fout *stdio.File, pt *Point, n int, abs int) int {
	var (
		i    int
		cur  Point
		prev Point
	)
	if abs != 0 {
		cur = func() Point {
			prev = *(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(n-1)))
			return prev
		}()
		svg_moveto(info, fout, dpoint(cur))
		for i = 0; i < n; i++ {
			if (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).X != cur.X && (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).Y != cur.Y {
				cur = prev
				svg_lineto(info, fout, dpoint(cur))
			}
			prev = *(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))
		}
		svg_lineto(info, fout, dpoint(*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(n-1)))))
	} else {
		cur = func() Point {
			prev = *(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*0))
			return prev
		}()
		svg_rmoveto(info, fout, dpoint(cur))
		for i = n - 1; i >= 0; i-- {
			if (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).X != cur.X && (*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))).Y != cur.Y {
				cur = prev
				svg_lineto(info, fout, dpoint(cur))
			}
			prev = *(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*uintptr(i)))
		}
		svg_lineto(info, fout, dpoint(*(*Point)(unsafe.Add(unsafe.Pointer(pt), unsafe.Sizeof(Point{})*0))))
	}
	newline = 1
	shiptoken(fout, libc.CString("z"))
	return 0
}
func write_paths_opaque(info *RenderConf, fout *stdio.File, tree *Path) {
	var (
		p *Path
		q *Path
	)
	for p = tree; p != nil; p = p.Sibling {
		if info.Grouping == 2 {
			stdio.Fprintf(fout, "<g>\n")
			stdio.Fprintf(fout, "<g>\n")
		}
		column = stdio.Fprintf(fout, "<path fill=\"#%06x\" stroke=\"none\" d=\"", info.Color)
		newline = 1
		lastop = 0
		if info.Debug {
			svg_jaggy_path(info, fout, &p.Priv.Pt[0], p.Priv.Len, 1)
		} else {
			svg_path(info, fout, &p.Curve, 1)
		}
		stdio.Fprintf(fout, "\"/>\n")
		for q = p.Childlist; q != nil; q = q.Sibling {
			column = stdio.Fprintf(fout, "<path fill=\"#%06x\" stroke=\"none\" d=\"", info.Fillcolor)
			newline = 1
			lastop = 0
			if info.Debug {
				svg_jaggy_path(info, fout, &q.Priv.Pt[0], q.Priv.Len, 1)
			} else {
				svg_path(info, fout, &q.Curve, 1)
			}
			stdio.Fprintf(fout, "\"/>\n")
		}
		if info.Grouping == 2 {
			stdio.Fprintf(fout, "</g>\n")
		}
		for q = p.Childlist; q != nil; q = q.Sibling {
			write_paths_opaque(info, fout, q.Childlist)
		}
		if info.Grouping == 2 {
			stdio.Fprintf(fout, "</g>\n")
		}
	}
}
func write_paths_transparent_rec(info *RenderConf, fout *stdio.File, tree *Path) {
	var (
		p *Path
		q *Path
	)
	for p = tree; p != nil; p = p.Sibling {
		if info.Grouping == 2 {
			stdio.Fprintf(fout, "<g>\n")
		}
		if info.Grouping != 0 {
			column = stdio.Fprintf(fout, "<path d=\"")
			newline = 1
			lastop = 0
		}
		if info.Debug {
			svg_jaggy_path(info, fout, &p.Priv.Pt[0], p.Priv.Len, 1)
		} else {
			svg_path(info, fout, &p.Curve, 1)
		}
		for q = p.Childlist; q != nil; q = q.Sibling {
			if info.Debug {
				svg_jaggy_path(info, fout, &q.Priv.Pt[0], q.Priv.Len, 0)
			} else {
				svg_path(info, fout, &q.Curve, 0)
			}
		}
		if info.Grouping != 0 {
			stdio.Fprintf(fout, "\"/>\n")
		}
		for q = p.Childlist; q != nil; q = q.Sibling {
			write_paths_transparent_rec(info, fout, q.Childlist)
		}
		if info.Grouping == 2 {
			stdio.Fprintf(fout, "</g>\n")
		}
	}
}
func write_paths_transparent(info *RenderConf, fout *stdio.File, tree *Path) {
	if info.Grouping == 0 {
		column = stdio.Fprintf(fout, "<path d=\"")
		newline = 1
		lastop = 0
	}
	write_paths_transparent_rec(info, fout, tree)
	if info.Grouping == 0 {
		stdio.Fprintf(fout, "\"/>\n")
	}
}
func page_svg(info *RenderConf, fout *stdio.File, plist *Path, imginfo *imgInfo) int {
	var (
		bboxx  float64 = imginfo.Trans.Bb[0] + imginfo.Lmar + imginfo.Rmar
		bboxy  float64 = imginfo.Trans.Bb[1] + imginfo.Tmar + imginfo.Bmar
		origx  float64 = imginfo.Trans.Orig[0] + imginfo.Lmar
		origy  float64 = bboxy - imginfo.Trans.Orig[1] - imginfo.Bmar
		scalex float64 = imginfo.Trans.Scalex / info.Unit
		scaley float64 = -imginfo.Trans.Scaley / info.Unit
	)
	stdio.Fprintf(fout, "<?xml version=\"1.0\" standalone=\"no\"?>\n")
	stdio.Fprintf(fout, "<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 20010904//EN\"\n")
	stdio.Fprintf(fout, " \"http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd\">\n")
	stdio.Fprintf(fout, "<svg version=\"1.0\" xmlns=\"http://www.w3.org/2000/svg\"\n")
	stdio.Fprintf(fout, " width=\"%fpt\" height=\"%fpt\" viewBox=\"0 0 %f %f\"\n", bboxx, bboxy, bboxx, bboxy)
	stdio.Fprintf(fout, " preserveAspectRatio=\"xMidYMid meet\">\n")
	stdio.Fprintf(fout, "<metadata>\n")
	stdio.Fprintf(fout, "Created by potrace "+Version+", written by Peter Selinger 2001-2019\n")
	stdio.Fprintf(fout, "</metadata>\n")
	stdio.Fprintf(fout, "<g transform=\"")
	if origx != 0 || origy != 0 {
		stdio.Fprintf(fout, "translate(%f,%f) ", origx, origy)
	}
	if info.Angle != 0 {
		stdio.Fprintf(fout, "rotate(%.2f) ", -info.Angle)
	}
	stdio.Fprintf(fout, "scale(%f,%f)", scalex, scaley)
	stdio.Fprintf(fout, "\"\n")
	stdio.Fprintf(fout, "fill=\"#%06x\" stroke=\"none\">\n", info.Color)
	if info.Opaque {
		write_paths_opaque(info, fout, plist)
	} else {
		write_paths_transparent(info, fout, plist)
	}
	stdio.Fprintf(fout, "</g>\n")
	stdio.Fprintf(fout, "</svg>\n")
	fout.Flush()
	return 0
}
func page_gimp(info *RenderConf, fout *stdio.File, plist *Path, imginfo *imgInfo) int {
	info.Opaque = false
	info.Grouping = 0
	return page_svg(info, fout, plist, imginfo)
}
