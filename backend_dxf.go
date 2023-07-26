package gotrace

import (
	"math"
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
)

func sub(v DPoint, w DPoint) DPoint {
	var r DPoint
	r.X = v.X - w.X
	r.Y = v.Y - w.Y
	return r
}
func xprodf(v DPoint, w DPoint) float64 {
	return v.X*w.Y - v.Y*w.X
}
func bulge(v DPoint, w DPoint) float64 {
	var (
		v2  float64
		w2  float64
		vw  float64
		vxw float64
		nvw float64
	)
	v2 = iprod(v, v)
	w2 = iprod(w, w)
	vw = iprod(v, w)
	vxw = xprodf(v, w)
	nvw = math.Sqrt(v2 * w2)
	if vxw == 0.0 {
		return 0.0
	}
	return (nvw - vw) / vxw
}
func dxf_ship(fout *stdio.File, gc int, fmt *byte, _rest ...interface{}) int {
	var (
		args libc.ArgList
		r    int
		c    int
	)
	r = stdio.Fprintf(fout, "%3d\n", gc)
	if r < 0 {
		return r
	}
	c = r
	args.Start(fmt, _rest)
	r = stdio.Vfprintf(fout, libc.GoString(fmt), args)
	args.End()
	if r < 0 {
		return r
	}
	c += r
	r = stdio.Fprintf(fout, "\n")
	if r < 0 {
		return r
	}
	c += r
	return c
}
func ship_polyline(fout *stdio.File, layer *byte, closed int) {
	dxf_ship(fout, 0, libc.CString("POLYLINE"))
	dxf_ship(fout, 8, libc.CString("%s"), layer)
	dxf_ship(fout, 66, libc.CString("%d"), 1)
	dxf_ship(fout, 70, libc.CString("%d"), func() int {
		if closed != 0 {
			return 1
		}
		return 0
	}())
}
func ship_vertex(fout *stdio.File, layer *byte, v DPoint, bulge float64) {
	dxf_ship(fout, 0, libc.CString("VERTEX"))
	dxf_ship(fout, 8, libc.CString("%s"), layer)
	dxf_ship(fout, 10, libc.CString("%f"), v.X)
	dxf_ship(fout, 20, libc.CString("%f"), v.Y)
	dxf_ship(fout, 42, libc.CString("%f"), bulge)
}
func ship_seqend(fout *stdio.File) {
	dxf_ship(fout, 0, libc.CString("SEQEND"))
}
func ship_comment(fout *stdio.File, comment *byte) {
	dxf_ship(fout, 999, libc.CString("%s"), comment)
}
func ship_section(fout *stdio.File, name *byte) {
	dxf_ship(fout, 0, libc.CString("SECTION"))
	dxf_ship(fout, 2, libc.CString("%s"), name)
}
func ship_endsec(fout *stdio.File) {
	dxf_ship(fout, 0, libc.CString("ENDSEC"))
}
func ship_eof(fout *stdio.File) {
	dxf_ship(fout, 0, libc.CString("EOF"))
}
func pseudo_quad(fout *stdio.File, layer *byte, A DPoint, C DPoint, B DPoint) {
	var (
		v      DPoint
		w      DPoint
		v2     float64
		w2     float64
		vw     float64
		vxw    float64
		nvw    float64
		a      float64
		b      float64
		c      float64
		y      float64
		G      DPoint
		bulge1 float64
		bulge2 float64
	)
	v = sub(A, C)
	w = sub(B, C)
	v2 = iprod(v, v)
	w2 = iprod(w, w)
	vw = iprod(v, w)
	vxw = xprodf(v, w)
	nvw = math.Sqrt(v2 * w2)
	a = v2 + vw*2 + w2
	b = v2 + nvw*2 + w2
	c = nvw * 4
	if vxw == 0 || a == 0 {
		goto degenerate
	}
	y = (b - math.Sqrt(b*b-a*c)) / a
	G = aux_interval(y, C, aux_interval(0.5, A, B))
	bulge1 = bulge(sub(A, G), v)
	bulge2 = bulge(w, sub(B, G))
	ship_vertex(fout, layer, A, -bulge1)
	ship_vertex(fout, layer, G, -bulge2)
	return
degenerate:
	ship_vertex(fout, layer, A, 0)
	return
}
func pseudo_bezier(fout *stdio.File, layer *byte, A DPoint, B DPoint, C DPoint, D DPoint) {
	var (
		E DPoint = aux_interval(0.75, A, B)
		G DPoint = aux_interval(0.75, D, C)
		F DPoint = aux_interval(0.5, E, G)
	)
	pseudo_quad(fout, layer, A, E, F)
	pseudo_quad(fout, layer, F, G, D)
	return
}
func dxf_path(fout *stdio.File, layer *byte, curve *Curve, t transT) int {
	var (
		i  int
		c  *DPoint
		c1 *DPoint
		n  int = curve.N
	)
	ship_polyline(fout, layer, 1)
	for i = 0; i < n; i++ {
		c = &curve.C[i][0]
		c1 = &curve.C[mod(i-1, n)][0]
		switch curve.Tag[i] {
		case POTRACE_CORNER:
			ship_vertex(fout, layer, trans(*(*DPoint)(unsafe.Add(unsafe.Pointer(c1), unsafe.Sizeof(DPoint{})*2)), t), 0)
			ship_vertex(fout, layer, trans(*(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*1)), t), 0)
		case POTRACE_CURVETO:
			pseudo_bezier(fout, layer, trans(*(*DPoint)(unsafe.Add(unsafe.Pointer(c1), unsafe.Sizeof(DPoint{})*2)), t), trans(*(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*0)), t), trans(*(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*1)), t), trans(*(*DPoint)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(DPoint{})*2)), t))
		}
	}
	ship_seqend(fout)
	return 0
}
func page_dxf(info *RenderConf, fout *stdio.File, plist *Path, imginfo *imgInfo) int {
	var (
		p     *Path
		t     transT
		layer *byte = libc.CString("0")
	)
	t.Bb[0] = imginfo.Trans.Bb[0] + imginfo.Lmar + imginfo.Rmar
	t.Bb[1] = imginfo.Trans.Bb[1] + imginfo.Tmar + imginfo.Bmar
	t.Orig[0] = imginfo.Trans.Orig[0] + imginfo.Lmar
	t.Orig[1] = imginfo.Trans.Orig[1] + imginfo.Bmar
	t.X[0] = imginfo.Trans.X[0]
	t.X[1] = imginfo.Trans.X[1]
	t.Y[0] = imginfo.Trans.Y[0]
	t.Y[1] = imginfo.Trans.Y[1]
	ship_comment(fout, libc.CString("DXF data, created by potrace "+Version+", written by Peter Selinger 2001-2019"))
	ship_section(fout, libc.CString("HEADER"))
	dxf_ship(fout, 9, libc.CString("$ACADVER"))
	dxf_ship(fout, 1, libc.CString("AC1006"))
	dxf_ship(fout, 9, libc.CString("$EXTMIN"))
	dxf_ship(fout, 10, libc.CString("%f"), 0.0)
	dxf_ship(fout, 20, libc.CString("%f"), 0.0)
	dxf_ship(fout, 30, libc.CString("%f"), 0.0)
	dxf_ship(fout, 9, libc.CString("$EXTMAX"))
	dxf_ship(fout, 10, libc.CString("%f"), t.Bb[0])
	dxf_ship(fout, 20, libc.CString("%f"), t.Bb[1])
	dxf_ship(fout, 30, libc.CString("%f"), 0.0)
	ship_endsec(fout)
	ship_section(fout, libc.CString("ENTITIES"))
	for p = plist; p != nil; p = p.Next {
		dxf_path(fout, layer, &p.Curve, t)
	}
	ship_endsec(fout)
	ship_eof(fout)
	fout.Flush()
	return 0
}
