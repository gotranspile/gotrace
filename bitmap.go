package gotrace

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func getsize(dy int, h int) int64 {
	var size int64
	if dy < 0 {
		dy = -dy
	}
	size = int64(dy) * int64(h) * int64(int(unsafe.Sizeof(Word(0))))
	if size < 0 || h != 0 && dy != 0 && size/int64(h)/int64(dy) != int64(int(unsafe.Sizeof(Word(0)))) {
		return -1
	}
	return size
}
func bm_size(bm *Bitmap) int64 {
	return getsize(bm.Dy, bm.H)
}
func bm_base(bm *Bitmap) *Word {
	var dy int = bm.Dy
	if dy >= 0 || bm.H == 0 {
		return &bm.Map[0]
	} else {
		return &bm.Map[int64(bm.H-1)*int64(bm.Dy)]
	}
}
func bm_free(bm *Bitmap) {
	if bm != nil && bm.Map != nil {
		libc.Free(unsafe.Pointer(bm_base(bm)))
	}

}
func NewBitmap(w int, h int) *Bitmap {
	var (
		bm *Bitmap
		dy int
	)
	if w == 0 {
		dy = 0
	} else {
		dy = (w-1)/(8*(int(unsafe.Sizeof(Word(0))))) + 1
	}
	var size int64
	size = getsize(dy, h)
	if size < 0 {
		panic("out of memory")
		return nil
	}
	if size == 0 {
		size = 1
	}
	bm = new(Bitmap)
	if bm == nil {
		return nil
	}
	bm.W = w
	bm.H = h
	bm.Dy = dy
	bm.Map = make([]Word, uintptr(size)/unsafe.Sizeof(Word(0)))
	if bm.Map == nil {

		return nil
	}
	return bm
}
func bm_clear(bm *Bitmap, c int) {
	var size int64 = bm_size(bm)
	libc.MemSet(unsafe.Pointer(bm_base(bm)), byte(int8(func() int {
		if c != 0 {
			return -1
		}
		return 0
	}())), int(size))
}
func bm_dup(bm *Bitmap) *Bitmap {
	var (
		bm1 *Bitmap = NewBitmap(bm.W, bm.H)
		y   int
	)
	if bm1 == nil {
		return nil
	}
	for y = 0; y < bm.H; y++ {
		libc.MemCpy(unsafe.Pointer(&bm1.Map[int64(y)*int64(bm1.Dy)]), unsafe.Pointer(&bm.Map[int64(y)*int64(bm.Dy)]), int(uint64(bm1.Dy)*uint64(int(unsafe.Sizeof(Word(0))))))
	}
	return bm1
}
func bm_invert(bm *Bitmap) {
	var (
		dy int = bm.Dy
		y  int
		i  int
		p  *Word
	)
	if dy < 0 {
		dy = -dy
	}
	for y = 0; y < bm.H; y++ {
		p = &bm.Map[int64(y)*int64(bm.Dy)]
		for i = 0; i < dy; i++ {
			*(*Word)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(Word(0))*uintptr(i))) ^= ^Word(0)
		}
	}
}
func bm_resize(bm *Bitmap, h int) int {
	var (
		dy      int = bm.Dy
		newsize int64
		newmap  []Word
	)
	if dy < 0 {
		bm_flip(bm)
	}
	newsize = getsize(dy, h)
	if newsize < 0 {
		panic("out of memory")
		goto error
	}
	if newsize == 0 {
		newsize = 1
	}
	newmap = make([]Word, uintptr(newsize)/unsafe.Sizeof(Word(0)))
	copy(newmap, bm.Map)
	if newmap == nil {
		goto error
	}
	bm.Map = []Word(newmap)
	bm.H = h
	if dy < 0 {
		bm_flip(bm)
	}
	return 0
error:
	if dy < 0 {
		bm_flip(bm)
	}
	return 1
}
