package gotrace

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"unsafe"
)

const GM_MODE_NONZERO = 1
const GM_MODE_ODD = 2
const GM_MODE_POSITIVE = 3
const GM_MODE_NEGATIVE = 4

type gm_sample_t int16
type Greymap struct {
	W    int
	H    int
	Dy   int
	Base []gm_sample_t
	Map  []gm_sample_t
}

func gm_getsize(dy int, h int) int64 {
	var size int64
	if dy < 0 {
		dy = -dy
	}
	size = int64(dy) * int64(h) * int64(unsafe.Sizeof(gm_sample_t(0)))
	if size < 0 || h != 0 && dy != 0 && size/int64(h)/int64(dy) != int64(unsafe.Sizeof(gm_sample_t(0))) {
		return -1
	}
	return size
}
func gm_size(gm *Greymap) int64 {
	return gm_getsize(gm.Dy, gm.H)
}
func NewGreymap(w int, h int) *Greymap {
	var (
		gm   *Greymap
		dy   int = w
		size int64
	)
	size = gm_getsize(dy, h)
	if size < 0 {
		panic("out of memory")
		return nil
	}
	if size == 0 {
		size = 1
	}
	gm = new(Greymap)
	if gm == nil {
		return nil
	}
	gm.W = w
	gm.H = h
	gm.Dy = dy
	gm.Base = make([]gm_sample_t, uintptr(size)/unsafe.Sizeof(gm_sample_t(0)))
	if gm.Base == nil {

		return nil
	}
	gm.Map = gm.Base
	return gm
}
func gm_free(gm *Greymap) {
	if gm != nil {

	}

}
func gm_dup(gm *Greymap) *Greymap {
	var (
		gm1 *Greymap = NewGreymap(gm.W, gm.H)
		y   int
	)
	if gm1 == nil {
		return nil
	}
	for y = 0; y < gm.H; y++ {
		libc.MemCpy(unsafe.Pointer(&gm1.Map[int64(y)*int64(gm1.Dy)]), unsafe.Pointer(&gm.Map[int64(y)*int64(gm.Dy)]), int(uint64(gm1.Dy)*uint64(unsafe.Sizeof(gm_sample_t(0)))))
	}
	return gm1
}
func gm_clear(gm *Greymap, b int) {
	var (
		size int64 = gm_size(gm)
		x    int
		y    int
	)
	if b == 0 {
		libc.MemSet(unsafe.Pointer(&gm.Base[0]), 0, int(size))
	} else {
		for y = 0; y < gm.H; y++ {
			for x = 0; x < gm.W; x++ {
				gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(b)
			}
		}
	}
}
func gm_resize(gm *Greymap, h int) int {
	var (
		dy      int = gm.Dy
		newsize int64
		newbase []gm_sample_t
	)
	if dy < 0 {
		gm_flip(gm)
	}
	newsize = gm_getsize(dy, h)
	if newsize < 0 {
		panic("out of memory")
		goto error
	}
	if newsize == 0 {
		newsize = 1
	}
	newbase = make([]gm_sample_t, uintptr(newsize)/unsafe.Sizeof(gm_sample_t(0)))
	copy(newbase, gm.Base)
	if newbase == nil {
		goto error
	}
	gm.Base = []gm_sample_t(newbase)
	gm.Map = []gm_sample_t(newbase)
	gm.H = h
	if dy < 0 {
		gm_flip(gm)
	}
	return 0
error:
	if dy < 0 {
		gm_flip(gm)
	}
	return 1
}

var gm_read_error *byte = nil

func GreymapRead(f *stdio.File, gmp **Greymap) int {
	var magic [2]int
	magic[0] = fgetc_ws(f)
	if magic[0] == stdio.EOF {
		return -3
	}
	magic[1] = f.GetC()
	if magic[0] == 'P' && magic[1] >= '1' && magic[1] <= '6' {
		return gm_readbody_pnm(f, gmp, magic[1])
	}
	if magic[0] == 'B' && magic[1] == 'M' {
		return gm_readbody_bmp(f, gmp)
	}
	return -4
}
func gm_readbody_pnm(f *stdio.File, gmp **Greymap, magic int) int {
	var (
		gm         *Greymap
		x          int
		y          int
		i          int
		j          int
		b          int
		b1         int
		sum        int
		bpr        int
		w          int
		h          int
		max        int
		realheight int
	)
	gm = nil
	w = readnum(f)
	if w < 0 {
		goto format_error
	}
	h = readnum(f)
	if h < 0 {
		goto format_error
	}
	gm = NewGreymap(w, h)
	if gm == nil {
		goto std_error
	}
	realheight = 0
	switch magic {
	default:
		goto format_error
	case '1':
		for y = 0; y < h; y++ {
			realheight = y + 1
			for x = 0; x < w; x++ {
				b = readbit(f)
				if b < 0 {
					goto eof
				}
				if b != 0 {
					gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = 0
				} else {
					gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = math.MaxUint8
				}
			}
		}
	case '2':
		max = readnum(f)
		if max < 1 {
			goto format_error
		}
		for y = 0; y < h; y++ {
			realheight = y + 1
			for x = 0; x < w; x++ {
				b = readnum(f)
				if b < 0 {
					goto eof
				}
				gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(b * math.MaxUint8 / max)
			}
		}
	case '3':
		max = readnum(f)
		if max < 1 {
			goto format_error
		}
		for y = 0; y < h; y++ {
			realheight = y + 1
			for x = 0; x < w; x++ {
				sum = 0
				for i = 0; i < 3; i++ {
					b = readnum(f)
					if b < 0 {
						goto eof
					}
					sum += b
				}
				gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(sum * (math.MaxUint8 / 3) / max)
			}
		}
	case '4':
		b = f.GetC()
		if b == stdio.EOF {
			goto format_error
		}
		bpr = (w + 7) / 8
		for y = 0; y < h; y++ {
			realheight = y + 1
			for i = 0; i < bpr; i++ {
				b = f.GetC()
				if b == stdio.EOF {
					goto eof
				}
				for j = 0; j < 8; j++ {
					if (i*8+j) >= 0 && (i*8+j) < gm.W && y >= 0 && y < gm.H {
						if b&(0x80>>j) != 0 {
							gm.Map[int64(y)*int64(gm.Dy)+int64(i*8+j)] = 0
						} else {
							gm.Map[int64(y)*int64(gm.Dy)+int64(i*8+j)] = math.MaxUint8
						}
					} else {
					}
				}
			}
		}
	case '5':
		max = readnum(f)
		if max < 1 {
			goto format_error
		}
		b = f.GetC()
		if b == stdio.EOF {
			goto format_error
		}
		for y = 0; y < h; y++ {
			realheight = y + 1
			for x = 0; x < w; x++ {
				b = f.GetC()
				if b == stdio.EOF {
					goto eof
				}
				if max >= 256 {
					b <<= 8
					b1 = f.GetC()
					if b1 == stdio.EOF {
						goto eof
					}
					b |= b1
				}
				gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(b * math.MaxUint8 / max)
			}
		}
	case '6':
		max = readnum(f)
		if max < 1 {
			goto format_error
		}
		b = f.GetC()
		if b == stdio.EOF {
			goto format_error
		}
		for y = 0; y < h; y++ {
			realheight = y + 1
			for x = 0; x < w; x++ {
				sum = 0
				for i = 0; i < 3; i++ {
					b = f.GetC()
					if b == stdio.EOF {
						goto eof
					}
					if max >= 256 {
						b <<= 8
						b1 = f.GetC()
						if b1 == stdio.EOF {
							goto eof
						}
						b |= b1
					}
					sum += b
				}
				gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(sum * (math.MaxUint8 / 3) / max)
			}
		}
	}
	gm_flip(gm)
	*gmp = gm
	return 0
eof:
	if gm_resize(gm, realheight) != 0 {
		goto std_error
	}
	gm_flip(gm)
	*gmp = gm
	return 1
format_error:
	gm_free(gm)
	if magic == '1' || magic == '4' {
		gm_read_error = libc.CString("invalid pbm file")
	} else if magic == '2' || magic == '5' {
		gm_read_error = libc.CString("invalid pgm file")
	} else {
		gm_read_error = libc.CString("invalid ppm file")
	}
	return -2
std_error:
	gm_free(gm)
	return -1
}
func gm_readbody_bmp(f *stdio.File, gmp **Greymap) int {
	var (
		bmpinfo    bmp_info_t
		coltable   *int
		b          uint
		c          uint
		i          uint
		j          uint
		gm         *Greymap
		x          uint
		y          uint
		col        [2]int
		bitbuf     uint
		n          uint
		redshift   uint
		greenshift uint
		blueshift  uint
		realheight int
	)
	gm_read_error = nil
	gm = nil
	coltable = nil
	bmp_pos = 2
	if bmp_readint(f, 4, &bmpinfo.FileSize) != 0 {
		goto try_error
	}
	if bmp_readint(f, 4, &bmpinfo.Reserved) != 0 {
		goto try_error
	}
	if bmp_readint(f, 4, &bmpinfo.DataOffset) != 0 {
		goto try_error
	}
	if bmp_readint(f, 4, &bmpinfo.InfoSize) != 0 {
		goto try_error
	}
	if bmpinfo.InfoSize == 40 || bmpinfo.InfoSize == 64 || bmpinfo.InfoSize == 108 || bmpinfo.InfoSize == 124 {
		bmpinfo.Ctbits = 32
		if bmp_readint(f, 4, &bmpinfo.W) != 0 {
			goto try_error
		}
		if bmp_readint(f, 4, &bmpinfo.H) != 0 {
			goto try_error
		}
		if bmp_readint(f, 2, &bmpinfo.Planes) != 0 {
			goto try_error
		}
		if bmp_readint(f, 2, &bmpinfo.Bits) != 0 {
			goto try_error
		}
		if bmp_readint(f, 4, &bmpinfo.Comp) != 0 {
			goto try_error
		}
		if bmp_readint(f, 4, &bmpinfo.ImageSize) != 0 {
			goto try_error
		}
		if bmp_readint(f, 4, &bmpinfo.XpixelsPerM) != 0 {
			goto try_error
		}
		if bmp_readint(f, 4, &bmpinfo.YpixelsPerM) != 0 {
			goto try_error
		}
		if bmp_readint(f, 4, &bmpinfo.Ncolors) != 0 {
			goto try_error
		}
		if bmp_readint(f, 4, &bmpinfo.ColorsImportant) != 0 {
			goto try_error
		}
		if bmpinfo.InfoSize >= 108 {
			if bmp_readint(f, 4, &bmpinfo.RedMask) != 0 {
				goto try_error
			}
			if bmp_readint(f, 4, &bmpinfo.GreenMask) != 0 {
				goto try_error
			}
			if bmp_readint(f, 4, &bmpinfo.BlueMask) != 0 {
				goto try_error
			}
			if bmp_readint(f, 4, &bmpinfo.AlphaMask) != 0 {
				goto try_error
			}
		}
		if bmpinfo.W > math.MaxInt32 {
			goto format_error
		}
		if bmpinfo.H > math.MaxInt32 {
			bmpinfo.H = (-bmpinfo.H) & math.MaxUint32
			bmpinfo.Topdown = 1
		} else {
			bmpinfo.Topdown = 0
		}
		if bmpinfo.H > math.MaxInt32 {
			goto format_error
		}
	} else if bmpinfo.InfoSize == 12 {
		bmpinfo.Ctbits = 24
		if bmp_readint(f, 2, &bmpinfo.W) != 0 {
			goto try_error
		}
		if bmp_readint(f, 2, &bmpinfo.H) != 0 {
			goto try_error
		}
		if bmp_readint(f, 2, &bmpinfo.Planes) != 0 {
			goto try_error
		}
		if bmp_readint(f, 2, &bmpinfo.Bits) != 0 {
			goto try_error
		}
		bmpinfo.Comp = 0
		bmpinfo.Ncolors = 0
		bmpinfo.Topdown = 0
	} else {
		goto format_error
	}
	if bmpinfo.Comp == 3 && bmpinfo.InfoSize < 108 {
		goto format_error
	}
	if bmpinfo.Comp > 3 || bmpinfo.Bits > 32 {
		goto format_error
	}
	if bmp_forward(f, int(bmpinfo.InfoSize+14)) != 0 {
		goto try_error
	}
	if bmpinfo.Planes != 1 {
		gm_read_error = libc.CString("cannot handle bmp planes")
		goto format_error
	}
	if bmpinfo.Ncolors == 0 && bmpinfo.Bits <= 8 {
		bmpinfo.Ncolors = 1 << bmpinfo.Bits
	}
	if bmpinfo.Bits <= 8 {
		coltable = &make([]int, int(bmpinfo.Ncolors))[0]
		if coltable == nil {
			goto std_error
		}
		for i = 0; i < bmpinfo.Ncolors; i++ {
			if bmp_readint(f, int(bmpinfo.Ctbits/8), &c) != 0 {
				goto try_error
			}
			c = ((c >> 16) & math.MaxUint8) + ((c >> 8) & math.MaxUint8) + (c & math.MaxUint8)
			*(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr(i))) = int(c / 3)
		}
	}
	if bmpinfo.InfoSize != 12 {
		if bmp_forward(f, int(bmpinfo.DataOffset)) != 0 {
			goto try_error
		}
	}
	gm = NewGreymap(int(bmpinfo.W), int(bmpinfo.H))
	if gm == nil {
		goto std_error
	}
	realheight = 0
	switch bmpinfo.Bits + bmpinfo.Comp*0x100 {
	default:
		goto format_error
	case 0x1:
		for y = 0; y < bmpinfo.H; y++ {
			realheight = int(y + 1)
			bmp_pad_reset()
			for i = 0; i*8 < bmpinfo.W; i++ {
				if bmp_readint(f, 1, &b) != 0 {
					goto eof
				}
				for j = 0; j < 8; j++ {
					if int(i*8+j) >= 0 && int(i*8+j) < gm.W && int(y) >= 0 && int(y) < gm.H {
						if b&(0x80>>j) != 0 {
							if 1 < bmpinfo.Ncolors {
								gm.Map[int64(y)*int64(gm.Dy)+int64(i*8+j)] = gm_sample_t(*(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*1)))
							} else {
								gm.Map[int64(y)*int64(gm.Dy)+int64(i*8+j)] = 0
							}
						} else if 0 < bmpinfo.Ncolors {
							gm.Map[int64(y)*int64(gm.Dy)+int64(i*8+j)] = gm_sample_t(*(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*0)))
						} else {
							gm.Map[int64(y)*int64(gm.Dy)+int64(i*8+j)] = 0
						}
					} else {
					}
				}
			}
			if bmp_pad(f) != 0 {
				goto try_error
			}
		}
	case 0x2:
		fallthrough
	case 0x3:
		fallthrough
	case 0x4:
		fallthrough
	case 0x5:
		fallthrough
	case 0x6:
		fallthrough
	case 0x7:
		fallthrough
	case 0x8:
		for y = 0; y < bmpinfo.H; y++ {
			realheight = int(y + 1)
			bmp_pad_reset()
			bitbuf = 0
			n = 0
			for x = 0; x < bmpinfo.W; x++ {
				if n < bmpinfo.Bits {
					if bmp_readint(f, 1, &b) != 0 {
						goto eof
					}
					bitbuf |= b << uint((8*unsafe.Sizeof(int(0)))-8-uintptr(n))
					n += 8
				}
				b = bitbuf >> uint((8*unsafe.Sizeof(int(0)))-uintptr(bmpinfo.Bits))
				bitbuf <<= bmpinfo.Bits
				n -= bmpinfo.Bits
				if b < bmpinfo.Ncolors {
					gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(*(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr(b))))
				} else {
					gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = 0
				}
			}
			if bmp_pad(f) != 0 {
				goto try_error
			}
		}
	case 0x10:
		gm_read_error = libc.CString("cannot handle bmp 16-bit coding")
		goto format_error
	case 0x18:
		fallthrough
	case 0x20:
		for y = 0; y < bmpinfo.H; y++ {
			realheight = int(y + 1)
			bmp_pad_reset()
			for x = 0; x < bmpinfo.W; x++ {
				if bmp_readint(f, int(bmpinfo.Bits/8), &c) != 0 {
					goto eof
				}
				c = ((c >> 16) & math.MaxUint8) + ((c >> 8) & math.MaxUint8) + (c & math.MaxUint8)
				gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(uint16(c / 3))
			}
			if bmp_pad(f) != 0 {
				goto try_error
			}
		}
	case 0x320:
		redshift = lobit(bmpinfo.RedMask)
		greenshift = lobit(bmpinfo.GreenMask)
		blueshift = lobit(bmpinfo.BlueMask)
		for y = 0; y < bmpinfo.H; y++ {
			realheight = int(y + 1)
			bmp_pad_reset()
			for x = 0; x < bmpinfo.W; x++ {
				if bmp_readint(f, int(bmpinfo.Bits/8), &c) != 0 {
					goto eof
				}
				c = ((c & bmpinfo.RedMask) >> redshift) + ((c & bmpinfo.GreenMask) >> greenshift) + ((c & bmpinfo.BlueMask) >> blueshift)
				gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(uint16(c / 3))
			}
			if bmp_pad(f) != 0 {
				goto try_error
			}
		}
	case 0x204:
		x = 0
		y = 0
		for {
			if bmp_readint(f, 1, &b) != 0 {
				goto eof
			}
			if bmp_readint(f, 1, &c) != 0 {
				goto eof
			}
			if b > 0 {
				if ((c >> 4) & 0xF) < bmpinfo.Ncolors {
					col[0] = *(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr((c>>4)&0xF)))
				} else {
					col[0] = 0
				}
				if (c & 0xF) < bmpinfo.Ncolors {
					col[1] = *(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr(c&0xF)))
				} else {
					col[1] = 0
				}
				for i = 0; i < b && x < bmpinfo.W; i++ {
					if x >= bmpinfo.W {
						x = 0
						y++
					}
					if x >= bmpinfo.W || y >= bmpinfo.H {
						break
					}
					realheight = int(y + 1)
					if int(x) >= 0 && int(x) < gm.W && int(y) >= 0 && int(y) < gm.H {
						gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(col[i&1])
					} else {
					}
					x++
				}
			} else if c == 0 {
				y++
				x = 0
			} else if c == 1 {
				break
			} else if c == 2 {
				if bmp_readint(f, 1, &b) != 0 {
					goto eof
				}
				if bmp_readint(f, 1, &c) != 0 {
					goto eof
				}
				x += b
				y += c
			} else {
				for i = 0; i < c; i++ {
					if (i & 1) == 0 {
						if bmp_readint(f, 1, &b) != 0 {
							goto eof
						}
					}
					if x >= bmpinfo.W {
						x = 0
						y++
					}
					if x >= bmpinfo.W || y >= bmpinfo.H {
						break
					}
					realheight = int(y + 1)
					if int(x) >= 0 && int(x) < gm.W && int(y) >= 0 && int(y) < gm.H {
						if ((b >> (4 - (i&1)*4)) & 0xF) < bmpinfo.Ncolors {
							gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(*(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr((b>>(4-(i&1)*4))&0xF))))
						} else {
							gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = 0
						}
					} else {
					}
					x++
				}
				if (c+1)&2 != 0 {
					if bmp_readint(f, 1, &b) != 0 {
						goto eof
					}
				}
			}
		}
	case 0x108:
		x = 0
		y = 0
		for {
			if bmp_readint(f, 1, &b) != 0 {
				goto eof
			}
			if bmp_readint(f, 1, &c) != 0 {
				goto eof
			}
			if b > 0 {
				for i = 0; i < b; i++ {
					if x >= bmpinfo.W {
						x = 0
						y++
					}
					if x >= bmpinfo.W || y >= bmpinfo.H {
						break
					}
					realheight = int(y + 1)
					if int(x) >= 0 && int(x) < gm.W && int(y) >= 0 && int(y) < gm.H {
						if c < bmpinfo.Ncolors {
							gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(*(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr(c))))
						} else {
							gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = 0
						}
					} else {
					}
					x++
				}
			} else if c == 0 {
				y++
				x = 0
			} else if c == 1 {
				break
			} else if c == 2 {
				if bmp_readint(f, 1, &b) != 0 {
					goto eof
				}
				if bmp_readint(f, 1, &c) != 0 {
					goto eof
				}
				x += b
				y += c
			} else {
				for i = 0; i < c; i++ {
					if bmp_readint(f, 1, &b) != 0 {
						goto eof
					}
					if x >= bmpinfo.W {
						x = 0
						y++
					}
					if x >= bmpinfo.W || y >= bmpinfo.H {
						break
					}
					realheight = int(y + 1)
					if int(x) >= 0 && int(x) < gm.W && int(y) >= 0 && int(y) < gm.H {
						if b < bmpinfo.Ncolors {
							gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = gm_sample_t(*(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr(b))))
						} else {
							gm.Map[int64(y)*int64(gm.Dy)+int64(x)] = 0
						}
					} else {
					}
					x++
				}
				if c&1 != 0 {
					if bmp_readint(f, 1, &b) != 0 {
						goto eof
					}
				}
			}
		}
	}
	bmp_forward(f, int(bmpinfo.FileSize))

	if bmpinfo.Topdown != 0 {
		gm_flip(gm)
	}
	*gmp = gm
	return 0
eof:
	if gm_resize(gm, realheight) != 0 {
		goto std_error
	}

	if bmpinfo.Topdown != 0 {
		gm_flip(gm)
	}
	*gmp = gm
	return 1
format_error:
try_error:
	;

	gm_free(gm)
	if gm_read_error == nil {
		gm_read_error = libc.CString("invalid bmp file")
	}
	return -2
std_error:

	gm_free(gm)
	return -1
}
func gm_writepgm(f *stdio.File, gm *Greymap, comment *byte, raw int, mode int, gamma float64) int {
	var (
		x          int
		y          int
		v          int
		gammatable [256]int
	)
	if gamma != 1.0 {
		gammatable[0] = 0
		for v = 1; v < 256; v++ {
			gammatable[v] = int(math.Exp(math.Log(float64(v)/255.0)/gamma)*math.MaxUint8 + 0.5)
		}
	} else {
		for v = 0; v < 256; v++ {
			gammatable[v] = v
		}
	}
	stdio.Fprintf(f, func() string {
		if raw != 0 {
			return "P5\n"
		}
		return "P2\n"
	}())
	if comment != nil && *comment != 0 {
		stdio.Fprintf(f, "# %s\n", comment)
	}
	stdio.Fprintf(f, "%d %d 255\n", gm.W, gm.H)
	for y = gm.H - 1; y >= 0; y-- {
		for x = 0; x < gm.W; x++ {
			v = int(gm.Map[int64(y)*int64(gm.Dy)+int64(x)])
			if mode == GM_MODE_NONZERO {
				if v > math.MaxUint8 {
					v = 510 - v
				}
				if v < 0 {
					v = 0
				}
			} else if mode == GM_MODE_ODD {
				if v >= 510 {
					v = v % 510
				} else if v >= 0 {
					v = v
				} else {
					v = 510 - 1 - (int(-1-v))%510
				}
				if v > math.MaxUint8 {
					v = 510 - v
				}
			} else if mode == GM_MODE_POSITIVE {
				if v < 0 {
					v = 0
				} else if v > math.MaxUint8 {
					v = math.MaxUint8
				}
			} else if mode == GM_MODE_NEGATIVE {
				v = 510 - v
				if v < 0 {
					v = 0
				} else if v > math.MaxUint8 {
					v = math.MaxUint8
				}
			}
			v = gammatable[v]
			if raw != 0 {
				f.PutC(v)
			} else {
				stdio.Fprintf(f, func() string {
					if x == gm.W-1 {
						return "%d\n"
					}
					return "%d "
				}(), v)
			}
		}
	}
	return 0
}
func gm_print(f *stdio.File, gm *Greymap) int {
	var (
		x  int
		y  int
		xx int
		yy int
		d  int
		t  int
		sw int
		sh int
	)
	if gm.W < 79 {
		sw = gm.W
	} else {
		sw = 79
	}
	if gm.W < 79 {
		sh = gm.H
	} else {
		sh = gm.H * sw * 44 / (gm.W * 79)
	}
	for yy = sh - 1; yy >= 0; yy-- {
		for xx = 0; xx < sw; xx++ {
			d = 0
			t = 0
			for x = xx * gm.W / sw; x < (xx+1)*gm.W/sw; x++ {
				for y = yy * gm.H / sh; y < (yy+1)*gm.H/sh; y++ {
					if x >= 0 && x < gm.W && y >= 0 && y < gm.H {
						d += int(gm.Map[int64(y)*int64(gm.Dy)+int64(x)])
					} else {
						d += 0
					}
					t += 256
				}
			}
			f.PutC(int("*#=- "[d*5/t]))
		}
		f.PutC('\n')
	}
	return 0
}
