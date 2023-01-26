package gotrace

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func detrand(x int, y int) int {
	var (
		z uint
		t [256]uint8 = [256]uint8{0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0}
	)
	z = uint(((x * 0x4B3E375) ^ y) * 0x5A8EF93)
	z = uint(int(t[z&math.MaxUint8]) ^ int(t[(z>>8)&math.MaxUint8]) ^ int(t[(z>>16)&math.MaxUint8]) ^ int(t[(z>>24)&math.MaxUint8]))
	return int(z)
}
func bm_clearexcess(bm *Bitmap) {
	var (
		mask Word
		y    int
	)
	if bm.W%(8*(int(unsafe.Sizeof(Word(0))))) != 0 {
		mask = (^Word(0)) << Word((8*(int(unsafe.Sizeof(Word(0)))))-bm.W%(8*(int(unsafe.Sizeof(Word(0))))))
		for y = 0; y < bm.H; y++ {
			bm.Map[int64(y)*int64(bm.Dy)+int64(bm.W/(8*(int(unsafe.Sizeof(Word(0))))))] &= mask
		}
	}
}

type bbox_s struct {
	X0 int
	X1 int
	Y0 int
	Y1 int
}
type bbox_t bbox_s

func clear_bm_with_bbox(bm *Bitmap, bbox *bbox_t) {
	var (
		imin int = (bbox.X0 / (8 * (int(unsafe.Sizeof(Word(0))))))
		imax int = ((bbox.X1 + 8*(int(unsafe.Sizeof(Word(0)))) - 1) / (8 * (int(unsafe.Sizeof(Word(0))))))
		i    int
		y    int
	)
	for y = bbox.Y0; y < bbox.Y1; y++ {
		for i = imin; i < imax; i++ {
			bm.Map[int64(y)*int64(bm.Dy)+int64(i)] = 0
		}
	}
}
func majority(bm *Bitmap, x int, y int) int {
	var (
		i  int
		a  int
		ct int
	)
	for i = 2; i < 5; i++ {
		ct = 0
		for a = -i + 1; a <= i-1; a++ {
			if func() bool {
				if (x+a) >= 0 && (x+a) < bm.W && ((y+i-1) >= 0 && (y+i-1) < bm.H) {
					return (bm.Map[int64(y+i-1)*int64(bm.Dy)+int64((x+a)/(8*(int(unsafe.Sizeof(Word(0))))))] & Word((1<<((8*(int(unsafe.Sizeof(Word(0)))))-1))>>((x+a)&((8*(int(unsafe.Sizeof(Word(0)))))-1)))) != 0
				}
				return false
			}() {
				ct += 1
			} else {
				ct += -1
			}
			if func() bool {
				if (x+i-1) >= 0 && (x+i-1) < bm.W && ((y+a-1) >= 0 && (y+a-1) < bm.H) {
					return (bm.Map[int64(y+a-1)*int64(bm.Dy)+int64((x+i-1)/(8*(int(unsafe.Sizeof(Word(0))))))] & Word((1<<((8*(int(unsafe.Sizeof(Word(0)))))-1))>>((x+i-1)&((8*(int(unsafe.Sizeof(Word(0)))))-1)))) != 0
				}
				return false
			}() {
				ct += 1
			} else {
				ct += -1
			}
			if func() bool {
				if (x+a-1) >= 0 && (x+a-1) < bm.W && ((y-i) >= 0 && (y-i) < bm.H) {
					return (bm.Map[int64(y-i)*int64(bm.Dy)+int64((x+a-1)/(8*(int(unsafe.Sizeof(Word(0))))))] & Word((1<<((8*(int(unsafe.Sizeof(Word(0)))))-1))>>((x+a-1)&((8*(int(unsafe.Sizeof(Word(0)))))-1)))) != 0
				}
				return false
			}() {
				ct += 1
			} else {
				ct += -1
			}
			if func() bool {
				if (x-i) >= 0 && (x-i) < bm.W && ((y+a) >= 0 && (y+a) < bm.H) {
					return (bm.Map[int64(y+a)*int64(bm.Dy)+int64((x-i)/(8*(int(unsafe.Sizeof(Word(0))))))] & Word((1<<((8*(int(unsafe.Sizeof(Word(0)))))-1))>>((x-i)&((8*(int(unsafe.Sizeof(Word(0)))))-1)))) != 0
				}
				return false
			}() {
				ct += 1
			} else {
				ct += -1
			}
		}
		if ct > 0 {
			return 1
		} else if ct < 0 {
			return 0
		}
	}
	return 0
}
func xor_to_ref(bm *Bitmap, x int, y int, xa int) {
	var (
		xhi int = x & (-(8 * (int(unsafe.Sizeof(Word(0))))))
		xlo int = x & ((8 * (int(unsafe.Sizeof(Word(0))))) - 1)
		i   int
	)
	if xhi < xa {
		for i = xhi; i < xa; i += 8 * (int(unsafe.Sizeof(Word(0)))) {
			bm.Map[int64(y)*int64(bm.Dy)+int64(i/(8*(int(unsafe.Sizeof(Word(0))))))] ^= ^Word(0)
		}
	} else {
		for i = xa; i < xhi; i += 8 * (int(unsafe.Sizeof(Word(0)))) {
			bm.Map[int64(y)*int64(bm.Dy)+int64(i/(8*(int(unsafe.Sizeof(Word(0))))))] ^= ^Word(0)
		}
	}
	if xlo != 0 {
		bm.Map[int64(y)*int64(bm.Dy)+int64(xhi/(8*(int(unsafe.Sizeof(Word(0))))))] ^= (^Word(0)) << Word((8*(int(unsafe.Sizeof(Word(0)))))-xlo)
	}
}
func xor_path(bm *Bitmap, p *Path) {
	var (
		xa int
		x  int
		y  int
		k  int
		y1 int
	)
	if p.Priv.Len <= 0 {
		return
	}
	y1 = p.Priv.Pt[p.Priv.Len-1].Y
	xa = p.Priv.Pt[0].X & (-(8 * (int(unsafe.Sizeof(Word(0))))))
	for k = 0; k < p.Priv.Len; k++ {
		x = p.Priv.Pt[k].X
		y = p.Priv.Pt[k].Y
		if y != y1 {
			xor_to_ref(bm, x, func() int {
				if y < y1 {
					return y
				}
				return y1
			}(), xa)
			y1 = y
		}
	}
}
func setbbox_path(bbox *bbox_t, p *Path) {
	var (
		x int
		y int
		k int
	)
	bbox.Y0 = math.MaxInt64
	bbox.Y1 = 0
	bbox.X0 = math.MaxInt64
	bbox.X1 = 0
	for k = 0; k < p.Priv.Len; k++ {
		x = p.Priv.Pt[k].X
		y = p.Priv.Pt[k].Y
		if x < bbox.X0 {
			bbox.X0 = x
		}
		if x > bbox.X1 {
			bbox.X1 = x
		}
		if y < bbox.Y0 {
			bbox.Y0 = y
		}
		if y > bbox.Y1 {
			bbox.Y1 = y
		}
	}
}
func findpath(bm *Bitmap, x0 int, y0 int, sign int, turnpolicy int) *Path {
	var (
		x       int
		y       int
		dirx    int
		diry    int
		len_    int
		size    int
		area    uint64
		c       int
		d       int
		tmp     int
		pt, pt1 []Point
		p       *Path = nil
	)
	x = x0
	y = y0
	dirx = 0
	diry = -1
	len_ = func() int {
		size = 0
		return size
	}()
	pt = nil
	area = 0
	for {
		if len_ >= size {
			size += 100
			size = int(float64(size) * 1.3)
			pt1 = make([]Point, size)
			copy(pt1, pt)
			if pt1 == nil {
				goto error
			}
			pt = pt1
		}
		pt[len_].X = x
		pt[len_].Y = y
		len_++
		x += dirx
		y += diry
		area += uint64(x * diry)
		if x == x0 && y == y0 {
			break
		}
		c = int(libc.BoolToInt(func() bool {
			if (x+(dirx+diry-1)/2) >= 0 && (x+(dirx+diry-1)/2) < bm.W && ((y+(diry-dirx-1)/2) >= 0 && (y+(diry-dirx-1)/2) < bm.H) {
				return (bm.Map[int64(y+(diry-dirx-1)/2)*int64(bm.Dy)+int64((x+(dirx+diry-1)/2)/(8*(int(unsafe.Sizeof(Word(0))))))] & Word((1<<((8*(int(unsafe.Sizeof(Word(0)))))-1))>>((x+(dirx+diry-1)/2)&((8*(int(unsafe.Sizeof(Word(0)))))-1)))) != 0
			}
			return false
		}()))
		d = int(libc.BoolToInt(func() bool {
			if (x+(dirx-diry-1)/2) >= 0 && (x+(dirx-diry-1)/2) < bm.W && ((y+(diry+dirx-1)/2) >= 0 && (y+(diry+dirx-1)/2) < bm.H) {
				return (bm.Map[int64(y+(diry+dirx-1)/2)*int64(bm.Dy)+int64((x+(dirx-diry-1)/2)/(8*(int(unsafe.Sizeof(Word(0))))))] & Word((1<<((8*(int(unsafe.Sizeof(Word(0)))))-1))>>((x+(dirx-diry-1)/2)&((8*(int(unsafe.Sizeof(Word(0)))))-1)))) != 0
			}
			return false
		}()))
		if c != 0 && d == 0 {
			if turnpolicy == POTRACE_TURNPOLICY_RIGHT || turnpolicy == POTRACE_TURNPOLICY_BLACK && sign == '+' || turnpolicy == POTRACE_TURNPOLICY_WHITE && sign == '-' || turnpolicy == POTRACE_TURNPOLICY_RANDOM && detrand(x, y) != 0 || turnpolicy == POTRACE_TURNPOLICY_MAJORITY && majority(bm, x, y) != 0 || turnpolicy == POTRACE_TURNPOLICY_MINORITY && majority(bm, x, y) == 0 {
				tmp = dirx
				dirx = diry
				diry = -tmp
			} else {
				tmp = dirx
				dirx = -diry
				diry = tmp
			}
		} else if c != 0 {
			tmp = dirx
			dirx = diry
			diry = -tmp
		} else if d == 0 {
			tmp = dirx
			dirx = -diry
			diry = tmp
		}
	}
	p = path_new()
	if p == nil {
		goto error
	}
	p.Priv.Pt = []Point(pt)
	p.Priv.Len = len_
	if area <= math.MaxInt64 {
		p.Area = int(area)
	} else {
		p.Area = math.MaxInt64
	}
	p.Sign = sign
	return p
error:

	return nil
}
func pathlist_to_tree(plist *Path, bm *Bitmap) {
	var (
		p          *Path
		p1         *Path
		heap       *Path
		heap1      *Path
		cur        *Path
		head       *Path
		plist_hook **Path
		hook_in    **Path
		hook_out   **Path
		bbox       bbox_t
	)
	bm_clear(bm, 0)
	for p = plist; p != nil; p = p.Next {
		p.Sibling = p.Next
		p.Childlist = nil
	}
	heap = plist
	for heap != nil {
		cur = heap
		heap = heap.Childlist
		cur.Childlist = nil
		head = cur
		cur = cur.Next
		head.Next = nil
		xor_path(bm, head)
		setbbox_path(&bbox, head)
		hook_in = &head.Childlist
		hook_out = &head.Next
		for p = cur; func() int {
			if p != nil {
				return func() int {
					cur = p.Next
					p.Next = nil
					return 1
				}()
			}
			return 0
		}() != 0; p = cur {
			if p.Priv.Pt[0].Y <= bbox.Y0 {
				for {
					p.Next = *hook_out
					*hook_out = p
					hook_out = &p.Next
					if true {
						break
					}
				}
				*hook_out = cur
				break
			}
			if func() bool {
				if p.Priv.Pt[0].X >= 0 && p.Priv.Pt[0].X < bm.W && ((p.Priv.Pt[0].Y-1) >= 0 && (p.Priv.Pt[0].Y-1) < bm.H) {
					return (bm.Map[int64(p.Priv.Pt[0].Y-1)*int64(bm.Dy)+int64(p.Priv.Pt[0].X/(8*(int(unsafe.Sizeof(Word(0))))))] & Word((1<<((8*(int(unsafe.Sizeof(Word(0)))))-1))>>(p.Priv.Pt[0].X&((8*(int(unsafe.Sizeof(Word(0)))))-1)))) != 0
				}
				return false
			}() {
				for {
					p.Next = *hook_in
					*hook_in = p
					hook_in = &p.Next
					if true {
						break
					}
				}
			} else {
				for {
					p.Next = *hook_out
					*hook_out = p
					hook_out = &p.Next
					if true {
						break
					}
				}
			}
		}
		clear_bm_with_bbox(bm, &bbox)
		if head.Next != nil {
			head.Next.Childlist = heap
			heap = head.Next
		}
		if head.Childlist != nil {
			head.Childlist.Childlist = heap
			heap = head.Childlist
		}
	}
	p = plist
	for p != nil {
		p1 = p.Sibling
		p.Sibling = p.Next
		p = p1
	}
	heap = plist
	if heap != nil {
		heap.Next = nil
	}
	plist = nil
	plist_hook = &plist
	for heap != nil {
		heap1 = heap.Next
		for p = heap; p != nil; p = p.Sibling {
			for {
				p.Next = *plist_hook
				*plist_hook = p
				plist_hook = &p.Next
				if true {
					break
				}
			}
			for p1 = p.Childlist; p1 != nil; p1 = p1.Sibling {
				for {
					p1.Next = *plist_hook
					*plist_hook = p1
					plist_hook = &p1.Next
					if true {
						break
					}
				}
				if p1.Childlist != nil {
					for {
						{
							var _hook **Path
							for _hook = &heap1; *_hook != nil; _hook = &(*_hook).Next {
							}
							for {
								p1.Childlist.Next = *_hook
								*_hook = p1.Childlist
								if true {
									break
								}
							}
						}
						if true {
							break
						}
					}
				}
			}
		}
		heap = heap1
	}
	return
}
func findnext(bm *Bitmap, xp *int, yp *int) int {
	var (
		x  int
		y  int
		x0 int
	)
	x0 = (*xp) & ^((8 * (int(unsafe.Sizeof(Word(0))))) - 1)
	for y = *yp; y >= 0; y-- {
		for x = x0; x < bm.W && x >= 0; x += int(uint(8 * (int(unsafe.Sizeof(Word(0)))))) {
			if bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(unsafe.Sizeof(Word(0))))))] != 0 {
				for !(func() bool {
					if x >= 0 && x < bm.W && (y >= 0 && y < bm.H) {
						return (bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(unsafe.Sizeof(Word(0))))))] & Word((1<<((8*(int(unsafe.Sizeof(Word(0)))))-1))>>(x&((8*(int(unsafe.Sizeof(Word(0)))))-1)))) != 0
					}
					return false
				}()) {
					x++
				}
				*xp = x
				*yp = y
				return 0
			}
		}
		x0 = 0
	}
	return 1
}
func bm_to_pathlist(bm *Bitmap, plistp **Path, param *Param, progress *progress_t) int {
	var (
		x          int
		y          int
		p          *Path
		plist      *Path   = nil
		plist_hook **Path  = &plist
		bm1        *Bitmap = nil
		sign       int
	)
	bm1 = bm_dup(bm)
	if bm1 == nil {
		goto error
	}
	bm_clearexcess(bm1)
	x = 0
	y = bm1.H - 1
	for findnext(bm1, &x, &y) == 0 {
		if func() bool {
			if x >= 0 && x < bm.W && (y >= 0 && y < bm.H) {
				return (bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(unsafe.Sizeof(Word(0))))))] & Word((1<<((8*(int(unsafe.Sizeof(Word(0)))))-1))>>(x&((8*(int(unsafe.Sizeof(Word(0)))))-1)))) != 0
			}
			return false
		}() {
			sign = '+'
		} else {
			sign = '-'
		}
		p = findpath(bm1, x, y+1, sign, param.Turnpolicy)
		if p == nil {
			goto error
		}
		xor_path(bm1, p)
		if p.Area <= param.Turdsize {
			path_free(p)
		} else {
			for {
				p.Next = *plist_hook
				*plist_hook = p
				plist_hook = &p.Next
				if true {
					break
				}
			}
		}
		if bm1.H > 0 {
			progress_update(1-float64(y)/float64(bm1.H), progress)
		}
	}
	pathlist_to_tree(plist, bm1)
	bm_free(bm1)
	*plistp = plist
	progress_update(1.0, progress)
	return 0
error:
	bm_free(bm1)
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
	return -1
}
