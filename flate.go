package gotrace

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

const OUTSIZE = 1000

func dummy_xship(f *stdio.File, filter int, s *byte, len_ int) int {
	f.WriteN(s, 1, len_)
	return len_
}
func pdf_xship(f *stdio.File, filter int, s *byte, len_ int) int {
	return dummy_xship(f, filter, s, len_)
}
func flate_xship(f *stdio.File, filter int, s *byte, len_ int) int {
	return dummy_xship(f, filter, s, len_)
}
func a85_xship(f *stdio.File, filter int, s *byte, len_ int) int {
	var (
		fstate int = 0
		n      int = 0
	)
	if filter != 0 && fstate == 0 {
		if filter == 1 {
			n += stdio.Fprintf(f, "currentfile /ASCII85Decode filter cvx exec\n")
		}
		n += a85init(f)
		fstate = 1
	} else if filter == 0 && fstate != 0 {
		n += a85finish(f)
		fstate = 0
	}
	if fstate == 0 {
		f.WriteN(s, 1, len_)
		return n + len_
	}
	n += a85write(f, s, len_)
	return n
}

var a85buf [4]uint
var a85n int
var a85col int

func a85init(f *stdio.File) int {
	a85n = 0
	a85col = 0
	return 0
}
func a85finish(f *stdio.File) int {
	var r int = 0
	if a85n != 0 {
		r += a85out(f, a85n)
	}
	f.PutS(libc.CString("~>\n"))
	return r + 2
}
func a85write(f *stdio.File, buf *byte, n int) int {
	var (
		i int
		r int = 0
	)
	for i = 0; i < n; i++ {
		a85buf[a85n] = uint(uint8(*(*byte)(unsafe.Add(unsafe.Pointer(buf), i))))
		a85n++
		if a85n == 4 {
			r += a85out(f, 4)
			a85n = 0
		}
	}
	return r
}
func a85out(f *stdio.File, n int) int {
	var (
		out [5]byte
		s   uint
		r   int = 0
		i   int
	)
	for i = n; i < 4; i++ {
		a85buf[i] = 0
	}
	s = (a85buf[0] << 24) + (a85buf[1] << 16) + (a85buf[2] << 8) + (a85buf[3] << 0)
	if s == 0 {
		r += a85spool(f, 'z')
	} else {
		for i = 4; i >= 0; i-- {
			out[i] = byte(s % 85)
			s /= 85
		}
		for i = 0; i < n+1; i++ {
			r += a85spool(f, int8(out[i]+33))
		}
	}
	return r
}
func a85spool(f *stdio.File, c int8) int {
	f.PutC(int(c))
	a85col++
	if a85col > 70 {
		f.PutC('\n')
		a85col = 0
		return 2
	}
	return 1
}
