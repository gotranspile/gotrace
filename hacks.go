package gotrace

import (
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/stdio"
)

const sizeofWord = unsafe.Sizeof(Word(0))

// Compression relies on C zlib, so we disable it.

func dummy_xship(f *stdio.File, filter int, s *byte, len_ int) int {
	f.WriteN(s, 1, len_)
	return len_
}
func pdf_xship(f *stdio.File, filter int, s *byte, len_ int) int {
	return dummy_xship(f, filter, s, len_)
}

// C flips bitmaps by using negative bitmap strides, which we cannot represent in Go with slices.

func bm_flip(bm *Bitmap) {
	dy := bm.Dy
	if dy < 0 {
		dy = -dy
	}
	// TODO: optimize
	bm2 := NewBitmap(bm.W, bm.H)
	for y := 0; y < bm.H; y++ {
		for x := 0; x < bm.W; x++ {
			bm2.Put(x, bm.H-1-y, bm.Get(x, y))
		}
	}
	*bm = *bm2
}

func gm_flip(gm *Greymap) {
	// TODO: implement
}

// unused in the code

type potrace_privstate_s struct{}
