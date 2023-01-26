package gotrace

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"unsafe"
)

type pdfColor int
type intarray_s struct {
	Size int
	Data []int
}

func intarray_init(ar *intarray_s) {
	ar.Size = 0
	ar.Data = nil
}
func intarray_term(ar *intarray_s) {

	ar.Size = 0
	ar.Data = nil
}
func intarray_set(ar *intarray_s, n int, val int) int {
	var (
		p []int
		s int
	)
	if n >= ar.Size {
		s = n + 1024
		old := p
		p = make([]int, s)
		copy(p, old)
		if p == nil {
			return -1
		}
		ar.Data = []int(p)
		ar.Size = s
	}
	ar.Data[n] = val
	return 0
}

var xref intarray_s
var nxref int = 0
var pages intarray_s
var npages int
var streamofs int
var outcount uint64
var xship func(f *stdio.File, filter int, s *byte, len_ int) int
var xship_file *stdio.File

func pdf_ship(fmt *byte, _rest ...interface{}) int {
	var (
		args libc.ArgList
		buf  [4096]byte
	)
	args.Start(fmt, _rest)
	stdio.Vsprintf(&buf[0], libc.GoString(fmt), args)
	buf[4095] = 0
	args.End()
	outcount += uint64(xship(xship_file, 1, &buf[0], libc.StrLen(&buf[0])))
	return 0
}
func shipclear(fmt *byte, _rest ...interface{}) int {
	var (
		buf  [4096]byte
		args libc.ArgList
	)
	args.Start(fmt, _rest)
	stdio.Vsprintf(&buf[0], libc.GoString(fmt), args)
	buf[4095] = 0
	args.End()
	outcount += uint64(xship(xship_file, 0, &buf[0], libc.StrLen(&buf[0])))
	return 0
}
func pdf_callbacks(info *RenderConf, fout *stdio.File) {
	if info.Compress {
		xship = pdf_xship
	} else {
		xship = dummy_xship
	}
	xship_file = fout
}
func pdf_unit(info *RenderConf, p DPoint) Point {
	var q Point
	q.X = int(math.Floor(p.X*info.Unit + 0.5))
	q.Y = int(math.Floor(p.Y*info.Unit + 0.5))
	return q
}
func pdf_coords(info *RenderConf, p DPoint) {
	var cur Point = pdf_unit(info, p)
	pdf_ship(libc.CString("%ld %ld "), cur.X, cur.Y)
}
func pdf_moveto(info *RenderConf, p DPoint) {
	pdf_coords(info, p)
	pdf_ship(libc.CString("m\n"))
}
func pdf_lineto(info *RenderConf, p DPoint) {
	pdf_coords(info, p)
	pdf_ship(libc.CString("l\n"))
}
func pdf_curveto(info *RenderConf, p1 DPoint, p2 DPoint, p3 DPoint) {
	var (
		q1 Point
		q2 Point
		q3 Point
	)
	q1 = pdf_unit(info, p1)
	q2 = pdf_unit(info, p2)
	q3 = pdf_unit(info, p3)
	pdf_ship(libc.CString("%ld %ld %ld %ld %ld %ld c\n"), q1.X, q1.Y, q2.X, q2.Y, q3.X, q3.Y)
}
func pdf_colorstring(col pdfColor) *byte {
	var (
		r   float64
		g   float64
		b   float64
		buf [100]byte
	)
	r = float64((col & 0xFF0000) >> 16)
	g = float64((col & 0xFF00) >> 8)
	b = float64((col & math.MaxUint8) >> 0)
	if r == 0 && g == 0 && b == 0 {
		return libc.CString("0 g")
	} else if r == math.MaxUint8 && g == math.MaxUint8 && b == math.MaxUint8 {
		return libc.CString("1 g")
	} else if r == g && g == b {
		stdio.Sprintf(&buf[0], "%.3f g", r/255.0)
		return &buf[0]
	} else {
		stdio.Sprintf(&buf[0], "%.3f %.3f %.3f rg", r/255.0, g/255.0, b/255.0)
		return &buf[0]
	}
}

var pdf_color pdfColor = -1

func pdf_setcolor(col pdfColor) {
	if col == pdf_color {
		return
	}
	pdf_color = col
	pdf_ship(libc.CString("%s\n"), pdf_colorstring(col))
}
func pdf_path(info *RenderConf, curve *Curve) int {
	var (
		i int
		c *DPoint
		m int = curve.N
	)
	c = &curve.C[m-1][0]
	pdf_moveto(info, *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*2)))
	for i = 0; i < m; i++ {
		c = &curve.C[i][0]
		switch curve.Tag[i] {
		case POTRACE_CORNER:
			pdf_lineto(info, *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*1)))
			pdf_lineto(info, *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*2)))
		case POTRACE_CURVETO:
			pdf_curveto(info, *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*0)), *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*1)), *(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*2)))
		}
	}
	return 0
}
func render0(info *RenderConf, plist *Path) int {
	var p *Path
	pdf_setcolor(pdfColor(info.Color))
	for p = plist; p != nil; p = p.Next {
		pdf_path(info, &p.Curve)
		pdf_ship(libc.CString("h\n"))
		if p.Next == nil || p.Next.Sign == '+' {
			pdf_ship(libc.CString("f\n"))
		}
	}
	return 0
}
func render0_opaque(info *RenderConf, plist *Path) int {
	var p *Path
	for p = plist; p != nil; p = p.Next {
		pdf_path(info, &p.Curve)
		pdf_ship(libc.CString("h\n"))
		pdf_setcolor(pdfColor(func() int {
			if p.Sign == '+' {
				return info.Color
			}
			return info.Fillcolor
		}()))
		pdf_ship(libc.CString("f\n"))
	}
	return 0
}
func pdf_render(info *RenderConf, plist *Path) int {
	if info.Opaque {
		return render0_opaque(info, plist)
	}
	return render0(info, plist)
}
func init_pdf(info *RenderConf, fout *stdio.File) int {
	intarray_init(&xref)
	intarray_init(&pages)
	nxref = 0
	npages = 0
	pdf_callbacks(info, fout)
	outcount = 0
	shipclear(libc.CString("%%PDF-1.3\n"))
	if intarray_set(&xref, func() int {
		p := &nxref
		x := *p
		*p++
		return x
	}(), int(outcount)) != 0 {
		goto try_error
	}
	shipclear(libc.CString("1 0 obj\n<</Type/Catalog/Pages 3 0 R>>\nendobj\n"))
	if intarray_set(&xref, func() int {
		p := &nxref
		x := *p
		*p++
		return x
	}(), int(outcount)) != 0 {
		goto try_error
	}
	shipclear(libc.CString("2 0 obj\n<</Creator(potrace " + Version + ", written by Peter Selinger 2001-2019)>>\nendobj\n"))
	nxref++
	fout.Flush()
	return 0
try_error:
	return 1
}
func term_pdf(info *RenderConf, fout *stdio.File) int {
	var (
		startxref int
		i         int
	)
	pdf_callbacks(info, fout)
	if intarray_set(&xref, 2, int(outcount)) != 0 {
		goto try_error
	}
	shipclear(libc.CString("3 0 obj\n<</Type/Pages/Count %d/Kids[\n"), npages)
	for i = 0; i < npages; i++ {
		shipclear(libc.CString("%d 0 R\n"), pages.Data[i])
	}
	shipclear(libc.CString("]>>\nendobj\n"))
	startxref = int(outcount)
	shipclear(libc.CString("xref\n0 %d\n"), nxref+1)
	shipclear(libc.CString("0000000000 65535 f \n"))
	for i = 0; i < nxref; i++ {
		shipclear(libc.CString("%0.10d 00000 n \n"), xref.Data[i])
	}
	shipclear(libc.CString("trailer\n<</Size %d/Root 1 0 R/Info 2 0 R>>\n"), nxref+1)
	shipclear(libc.CString("startxref\n%d\n%%%%EOF\n"), startxref)
	fout.Flush()
	intarray_term(&xref)
	intarray_term(&pages)
	return 0
try_error:
	return 1
}
func pdf_pageinit(info *RenderConf, imginfo *imgInfo, largebbox int) int {
	var (
		origx float64 = imginfo.Trans.Orig[0] + imginfo.Lmar
		origy float64 = imginfo.Trans.Orig[1] + imginfo.Bmar
		dxx   float64 = imginfo.Trans.X[0] / info.Unit
		dxy   float64 = imginfo.Trans.X[1] / info.Unit
		dyx   float64 = imginfo.Trans.Y[0] / info.Unit
		dyy   float64 = imginfo.Trans.Y[1] / info.Unit
		pagew float64 = imginfo.Trans.Bb[0] + imginfo.Lmar + imginfo.Rmar
		pageh float64 = imginfo.Trans.Bb[1] + imginfo.Tmar + imginfo.Bmar
	)
	pdf_color = -1
	if intarray_set(&xref, func() int {
		p := &nxref
		x := *p
		*p++
		return x
	}(), int(outcount)) != 0 {
		goto try_error
	}
	shipclear(libc.CString("%d 0 obj\n"), nxref)
	shipclear(libc.CString("<</Type/Page/Parent 3 0 R/Resources<</ProcSet[/PDF]>>"))
	if largebbox != 0 {
		shipclear(libc.CString("/MediaBox[0 0 %d %d]"), info.Paperwidth, info.Paperheight)
	} else {
		shipclear(libc.CString("/MediaBox[0 0 %f %f]"), pagew, pageh)
	}
	shipclear(libc.CString("/Contents %d 0 R>>\n"), nxref+1)
	shipclear(libc.CString("endobj\n"))
	if intarray_set(&pages, func() int {
		p := &npages
		x := *p
		*p++
		return x
	}(), nxref) != 0 {
		goto try_error
	}
	if intarray_set(&xref, func() int {
		p := &nxref
		x := *p
		*p++
		return x
	}(), int(outcount)) != 0 {
		goto try_error
	}
	shipclear(libc.CString("%d 0 obj\n"), nxref)
	if info.Compress {
		shipclear(libc.CString("<</Filter/FlateDecode/Length %d 0 R>>\n"), nxref+1)
	} else {
		shipclear(libc.CString("<</Length %d 0 R>>\n"), nxref+1)
	}
	shipclear(libc.CString("stream\n"))
	streamofs = int(outcount)
	pdf_ship(libc.CString("%f %f %f %f %f %f cm\n"), dxx, dxy, dyx, dyy, origx, origy)
	return 0
try_error:
	return 1
}
func pdf_pageterm() int {
	var streamlen int
	shipclear(libc.CString(""))
	streamlen = int(outcount - uint64(streamofs))
	shipclear(libc.CString("endstream\nendobj\n"))
	if intarray_set(&xref, func() int {
		p := &nxref
		x := *p
		*p++
		return x
	}(), int(outcount)) != 0 {
		goto try_error
	}
	shipclear(libc.CString("%d 0 obj\n%d\nendobj\n"), nxref, streamlen)
	return 0
try_error:
	return 1
}
func page_pdf(info *RenderConf, fout *stdio.File, plist *Path, imginfo *imgInfo) int {
	var r int
	pdf_callbacks(info, fout)
	if pdf_pageinit(info, imginfo, 0) != 0 {
		goto try_error
	}
	r = pdf_render(info, plist)
	if r != 0 {
		return r
	}
	if pdf_pageterm() != 0 {
		goto try_error
	}
	fout.Flush()
	return 0
try_error:
	return 1
}
func page_pdfpage(info *RenderConf, fout *stdio.File, plist *Path, imginfo *imgInfo) int {
	var r int
	pdf_callbacks(info, fout)
	if pdf_pageinit(info, imginfo, 1) != 0 {
		goto try_error
	}
	r = pdf_render(info, plist)
	if r != 0 {
		return r
	}
	if pdf_pageterm() != 0 {
		goto try_error
	}
	fout.Flush()
	return 0
try_error:
	return 1
}
