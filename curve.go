package gotrace

type privCurve struct {
	N          int
	Tag        []int
	C          [][3]DPoint
	Alphacurve int
	Vertex     []DPoint
	Alpha      []float64
	Alpha0     []float64
	Beta       []float64
}
type sums_s struct {
	X  float64
	Y  float64
	X2 float64
	Xy float64
	Y2 float64
}
type potrace_privpath_s struct {
	Len    int
	Pt     []Point
	Lon    []int
	X0     int
	Y0     int
	Sums   []sums_s
	M      int
	Po     []int
	Curve  privCurve
	Ocurve privCurve
	Fcurve *privCurve
}

func path_new() *Path {
	var (
		p    *Path               = nil
		priv *potrace_privpath_s = nil
	)
	if (func() *Path {
		p = new(Path)
		return p
	}()) == nil {
		goto calloc_error
	}
	*p = Path{}
	if (func() *potrace_privpath_s {
		priv = new(potrace_privpath_s)
		return priv
	}()) == nil {
		goto calloc_error
	}
	*priv = potrace_privpath_s{}
	p.Priv = priv
	return p
calloc_error:

	return nil
}
func privcurve_free_members(curve *privCurve) {

}
func path_free(p *Path) {
	if p != nil {
		if p.Priv != nil {

			privcurve_free_members(&p.Priv.Curve)
			privcurve_free_members(&p.Priv.Ocurve)
		}

	}

}
func pathlist_free(plist *Path) {
	var p *Path
	for p = plist; func() int {
		if p != nil {
			return func() int {
				plist = p.Next
				p.Next = nil
				return 1
			}()
		}
		return 0
	}() != 0; p = plist {
		path_free(p)
	}
}
func privcurve_init(curve *privCurve, n int) int {
	*curve = privCurve{}
	curve.N = n
	if (func() []int {
		p := &curve.Tag
		curve.Tag = make([]int, n)
		return *p
	}()) == nil {
		goto calloc_error
	}
	if (func() [][3]DPoint {
		p := &curve.C
		curve.C = make([][3]DPoint, n)
		return *p
	}()) == nil {
		goto calloc_error
	}
	if (func() []DPoint {
		p := &curve.Vertex
		curve.Vertex = make([]DPoint, n)
		return *p
	}()) == nil {
		goto calloc_error
	}
	if (func() []float64 {
		p := &curve.Alpha
		curve.Alpha = make([]float64, n)
		return *p
	}()) == nil {
		goto calloc_error
	}
	if (func() []float64 {
		p := &curve.Alpha0
		curve.Alpha0 = make([]float64, n)
		return *p
	}()) == nil {
		goto calloc_error
	}
	if (func() []float64 {
		p := &curve.Beta
		curve.Beta = make([]float64, n)
		return *p
	}()) == nil {
		goto calloc_error
	}
	return 0
calloc_error:

	return 1
}
func privcurve_to_curve(pc *privCurve, c *Curve) {
	c.N = pc.N
	c.Tag = pc.Tag
	c.C = pc.C
}
