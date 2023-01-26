package gotrace

import (
	"unsafe"
)

const TurnBlack = 0
const TurnWhite = 1
const TurnLeft = 2
const TurnRight = 3
const TurnMinority = 4
const TurnMajority = 5
const TurnRandom = 6
const POTRACE_CURVETO = 1
const POTRACE_CORNER = 2
const statusOK = 0
const statusIncomplete = 1

type Progress struct {
	Callback func(progress float64, privdata unsafe.Pointer)
	Data     unsafe.Pointer
	Min      float64
	Max      float64
	Epsilon  float64
}
type Config struct {
	TurdSize     int
	TurnPolicy   int
	AlphaMax     float64
	OptiCurve    bool
	OptTolerance float64
	Progress     Progress
}
type Word uint
type Bitmap struct {
	W   int
	H   int
	Dy  int
	Map []Word
}
type DPoint struct {
	X float64
	Y float64
}
type Curve struct {
	N   int
	Tag []int
	C   [][3]DPoint
}
type Path struct {
	Area      int
	Sign      int
	Curve     Curve
	Next      *Path
	Childlist *Path
	Sibling   *Path
	Priv      *potrace_privpath_s
}
type traceState struct {
	Status int
	Plist  *Path
	Priv   *potrace_privstate_s
}

var param_default Config = Config{TurdSize: 2, TurnPolicy: TurnMinority, AlphaMax: 1.0, OptiCurve: true, OptTolerance: 0.2, Progress: Progress{Callback: nil, Data: nil, Min: 0.0, Max: 1.0, Epsilon: 0.0}}

func DefaultConfig() *Config {
	var p *Config
	p = new(Config)
	if p == nil {
		return nil
	}
	*p = param_default
	return p
}
func traceBitmap(param *Config, bm *Bitmap) *traceState {
	var (
		r       int
		plist   *Path = nil
		st      *traceState
		prog    progress
		subprog progress
	)
	prog.Callback = param.Progress.Callback
	prog.Data = param.Progress.Data
	prog.Min = param.Progress.Min
	prog.Max = param.Progress.Max
	prog.Epsilon = param.Progress.Epsilon
	prog.D_prev = param.Progress.Min
	st = new(traceState)
	if st == nil {
		return nil
	}
	progress_subrange_start(0.0, 0.1, &prog, &subprog)
	r = bm_to_pathlist(bm, &plist, param, &subprog)
	if r != 0 {

		return nil
	}
	st.Status = statusOK
	st.Plist = plist
	st.Priv = nil
	progress_subrange_end(&prog, &subprog)
	progress_subrange_start(0.1, 1.0, &prog, &subprog)
	r = process_path(plist, param, &subprog)
	if r != 0 {
		st.Status = statusIncomplete
	}
	progress_subrange_end(&prog, &subprog)
	return st
}
