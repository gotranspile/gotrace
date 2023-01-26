package gotrace

import "unsafe"

type progress struct {
	Callback func(progress float64, privdata unsafe.Pointer)
	Data     unsafe.Pointer
	Min      float64
	Max      float64
	Epsilon  float64
	B        float64
	D_prev   float64
}

func progress_update(d float64, prog *progress) {
	var d_scaled float64
	if prog != nil && prog.Callback != nil {
		d_scaled = prog.Min*(1-d) + prog.Max*d
		if d == 1.0 || d_scaled >= prog.D_prev+prog.Epsilon {
			prog.Callback(prog.Min*(1-d)+prog.Max*d, prog.Data)
			prog.D_prev = d_scaled
		}
	}
}
func progress_subrange_start(a float64, b float64, prog *progress, sub *progress) {
	var (
		min float64
		max float64
	)
	if prog == nil || prog.Callback == nil {
		sub.Callback = nil
		return
	}
	min = prog.Min*(1-a) + prog.Max*a
	max = prog.Min*(1-b) + prog.Max*b
	if max-min < prog.Epsilon {
		sub.Callback = nil
		sub.B = b
		return
	}
	sub.Callback = prog.Callback
	sub.Data = prog.Data
	sub.Epsilon = prog.Epsilon
	sub.Min = min
	sub.Max = max
	sub.D_prev = prog.D_prev
	return
}
func progress_subrange_end(prog *progress, sub *progress) {
	if prog != nil && prog.Callback != nil {
		if sub.Callback == nil {
			progress_update(sub.B, prog)
		} else {
			prog.D_prev = sub.D_prev
		}
	}
}
