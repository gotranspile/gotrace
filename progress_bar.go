package gotrace

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"unsafe"
)

const COL0 = "\x1b[G"

type progress_bar_s struct {
	Init func(prog *potrace_progress_t, filename *byte, count int) int
	Term func(prog *potrace_progress_t)
}
type progress_bar_t progress_bar_s
type vt100_progress_s struct {
	Name  [22]byte
	Dnext float64
}
type vt100_progress_t vt100_progress_s

func vt100_progress(d float64, data unsafe.Pointer) {
	var (
		p *vt100_progress_t = (*vt100_progress_t)(data)
		b [41]byte          = func() [41]byte {
			var t [41]byte
			copy(t[:], []byte("========================================"))
			return t
		}()
		tick int
		perc int
	)
	if d >= p.Dnext {
		tick = int(math.Floor(d*40 + 0.01))
		perc = int(math.Floor(d*100 + 0.025))
		stdio.Fprintf(stdio.Stderr(), "%-21s |%-40s| %d%% \x1b[G", &p.Name[0], (*byte)(unsafe.Add(unsafe.Pointer(&b[40]), -tick)), perc)
		stdio.Stderr().Flush()
		p.Dnext = (float64(tick) + 0.995) / 40.0
	}
}
func init_vt100_progress(prog *potrace_progress_t, filename *byte, count int) int {
	var (
		p    *vt100_progress_t
		q    *byte
		s    *byte
		len_ int
	)
	p = new(vt100_progress_t)
	if p == nil {
		return 1
	}
	p.Dnext = 0
	if count != 0 {
		stdio.Sprintf(&p.Name[0], " (p.%d):", count+1)
	} else {
		s = filename
		if (func() *byte {
			q = libc.StrRChr(s, '/')
			return q
		}()) != nil {
			s = (*byte)(unsafe.Add(unsafe.Pointer(q), 1))
		}
		len_ = libc.StrLen(s)
		libc.StrNCpy(&p.Name[0], s, 21)
		p.Name[20] = 0
		if len_ > 20 {
			p.Name[17] = '.'
			p.Name[18] = '.'
			p.Name[19] = '.'
		}
		libc.StrCat(&p.Name[0], libc.CString(":"))
	}
	prog.Callback = vt100_progress
	prog.Data = unsafe.Pointer(p)
	prog.Min = 0.0
	prog.Max = 1.0
	prog.Epsilon = 0.0
	vt100_progress(0.0, prog.Data)
	return 0
}
func term_vt100_progress(prog *potrace_progress_t) {
	stdio.Fprintf(stdio.Stderr(), "\n")
	stdio.Stderr().Flush()
	prog.Data = nil
	return
}

var progress_bar_vt100_struct progress_bar_t = progress_bar_t{Init: init_vt100_progress, Term: term_vt100_progress}
var progress_bar_vt100 *progress_bar_t = &progress_bar_vt100_struct

type simplified_progress_s struct {
	N     int
	Dnext float64
}
type simplified_progress_t simplified_progress_s

func simplified_progress(d float64, data unsafe.Pointer) {
	var (
		p    *simplified_progress_t = (*simplified_progress_t)(data)
		tick int
	)
	if d >= p.Dnext {
		tick = int(math.Floor(d*40 + 0.01))
		for p.N < tick {
			stdio.Stderr().PutC('=')
			p.N++
		}
		stdio.Stderr().Flush()
		p.Dnext = (float64(tick) + 0.995) / 40.0
	}
}
func init_simplified_progress(prog *potrace_progress_t, filename *byte, count int) int {
	var (
		p    *simplified_progress_t
		q    *byte
		s    *byte
		len_ int
		buf  [22]byte
	)
	p = new(simplified_progress_t)
	if p == nil {
		return 1
	}
	p.N = 0
	p.Dnext = 0
	if count != 0 {
		stdio.Sprintf(&buf[0], " (p.%d):", count+1)
	} else {
		s = filename
		if (func() *byte {
			q = libc.StrRChr(s, '/')
			return q
		}()) != nil {
			s = (*byte)(unsafe.Add(unsafe.Pointer(q), 1))
		}
		len_ = libc.StrLen(s)
		libc.StrNCpy(&buf[0], s, 21)
		buf[20] = 0
		if len_ > 20 {
			buf[17] = '.'
			buf[18] = '.'
			buf[19] = '.'
		}
		libc.StrCat(&buf[0], libc.CString(":"))
	}
	stdio.Fprintf(stdio.Stderr(), "%-21s |", &buf[0])
	prog.Callback = simplified_progress
	prog.Data = unsafe.Pointer(p)
	prog.Min = 0.0
	prog.Max = 1.0
	prog.Epsilon = 0.0
	simplified_progress(0.0, prog.Data)
	return 0
}
func term_simplified_progress(prog *potrace_progress_t) {
	var p *simplified_progress_t = (*simplified_progress_t)(prog.Data)
	simplified_progress(1.0, unsafe.Pointer(p))
	stdio.Fprintf(stdio.Stderr(), "| 100%%\n")
	stdio.Stderr().Flush()

	return
}

var progress_bar_simplified_struct progress_bar_t = progress_bar_t{Init: init_simplified_progress, Term: term_simplified_progress}
var progress_bar_simplified *progress_bar_t = &progress_bar_simplified_struct
