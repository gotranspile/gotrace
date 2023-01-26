package gotrace

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const POTRACE_TURNPOLICY_BLACK = 0
const POTRACE_TURNPOLICY_WHITE = 1
const POTRACE_TURNPOLICY_LEFT = 2
const POTRACE_TURNPOLICY_RIGHT = 3
const POTRACE_TURNPOLICY_MINORITY = 4
const POTRACE_TURNPOLICY_MAJORITY = 5
const POTRACE_TURNPOLICY_RANDOM = 6
const POTRACE_CURVETO = 1
const POTRACE_CORNER = 2
const POTRACE_STATUS_OK = 0
const POTRACE_STATUS_INCOMPLETE = 1

type potrace_progress_s struct {
	Callback func(progress float64, privdata unsafe.Pointer)
	Data     unsafe.Pointer
	Min      float64
	Max      float64
	Epsilon  float64
}
type potrace_progress_t potrace_progress_s
type Param struct {
	Turdsize     int
	Turnpolicy   int
	Alphamax     float64
	Opticurve    bool
	Opttolerance float64
	Progress     potrace_progress_t
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
type State struct {
	Status int
	Plist  *Path
	Priv   *potrace_privstate_s
}

var param_default Param = Param{Turdsize: 2, Turnpolicy: POTRACE_TURNPOLICY_MINORITY, Alphamax: 1.0, Opticurve: true, Opttolerance: 0.2, Progress: potrace_progress_t{Callback: nil, Data: nil, Min: 0.0, Max: 1.0, Epsilon: 0.0}}

func ParamDefault() *Param {
	var p *Param
	p = new(Param)
	if p == nil {
		return nil
	}
	libc.MemCpy(unsafe.Pointer(p), unsafe.Pointer(&param_default), int(unsafe.Sizeof(Param{})))
	return p
}
func Trace(param *Param, bm *Bitmap) *State {
	var (
		r       int
		plist   *Path = nil
		st      *State
		prog    progress_t
		subprog progress_t
	)
	prog.Callback = param.Progress.Callback
	prog.Data = param.Progress.Data
	prog.Min = param.Progress.Min
	prog.Max = param.Progress.Max
	prog.Epsilon = param.Progress.Epsilon
	prog.D_prev = param.Progress.Min
	st = new(State)
	if st == nil {
		return nil
	}
	progress_subrange_start(0.0, 0.1, &prog, &subprog)
	r = bm_to_pathlist(bm, &plist, param, &subprog)
	if r != 0 {

		return nil
	}
	st.Status = POTRACE_STATUS_OK
	st.Plist = plist
	st.Priv = nil
	progress_subrange_end(&prog, &subprog)
	progress_subrange_start(0.1, 1.0, &prog, &subprog)
	r = process_path(plist, param, &subprog)
	if r != 0 {
		st.Status = POTRACE_STATUS_INCOMPLETE
	}
	progress_subrange_end(&prog, &subprog)
	return st
}
func potrace_state_free(st *State) {
	pathlist_free(st.Plist)

}
func potrace_param_free(p *Param) {

}
func potrace_version() *byte {
	return libc.CString("potracelib dev")
}
