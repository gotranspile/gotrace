package gotrace

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"unsafe"
)

func fgetc_ws(f *stdio.File) int {
	var c int
	for {
		c = f.GetC()
		if c == '#' {
			for {
				c = f.GetC()
				if c == '\n' || c == stdio.EOF {
					break
				}
			}
		}
		if c != ' ' && c != '\t' && c != '\r' && c != '\n' && c != 12 {
			return c
		}
	}
}
func readnum(f *stdio.File) int {
	var (
		c   int
		acc uint64
	)
	for {
		c = fgetc_ws(f)
		if c == stdio.EOF {
			return -1
		}
		if c >= '0' && c <= '9' {
			break
		}
	}
	acc = uint64(c - '0')
	for {
		c = f.GetC()
		if c == stdio.EOF {
			break
		}
		if c < '0' || c > '9' {
			f.UnGetC(c)
			break
		}
		acc *= 10
		acc += uint64(c - '0')
		if acc > math.MaxInt32 {
			return -1
		}
	}
	return int(acc)
}
func readbit(f *stdio.File) int {
	var c int
	for {
		c = fgetc_ws(f)
		if c == stdio.EOF {
			return -1
		}
		if c >= '0' && c <= '1' {
			break
		}
	}
	return c - '0'
}

var bm_read_error *byte = nil

func bitmapRead(f *stdio.File, threshold float64, bmp **Bitmap) int {
	var magic [2]int
	magic[0] = fgetc_ws(f)
	if magic[0] == stdio.EOF {
		return -3
	}
	magic[1] = f.GetC()
	if magic[0] == 'P' && magic[1] >= '1' && magic[1] <= '6' {
		return bm_readbody_pnm(f, threshold, bmp, magic[1])
	}
	if magic[0] == 'B' && magic[1] == 'M' {
		return bm_readbody_bmp(f, threshold, bmp)
	}
	return -4
}
func bm_readbody_pnm(f *stdio.File, threshold float64, bmp **Bitmap, magic int) int {
	var (
		bm         *Bitmap
		x          int
		y          int
		i          int
		b          int
		b1         int
		sum        int
		bpr        int
		w          int
		h          int
		max        int
		realheight int
	)
	bm = nil
	w = readnum(f)
	if w < 0 {
		goto format_error
	}
	h = readnum(f)
	if h < 0 {
		goto format_error
	}
	bm = NewBitmap(w, h)
	if bm == nil {
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
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> (x & ((8 * (int(sizeofWord))) - 1)))
				} else {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> (x & ((8 * (int(sizeofWord))) - 1))))
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
				if (func() int {
					if float64(b) > threshold*float64(max) {
						return 0
					}
					return 1
				}()) != 0 {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> (x & ((8 * (int(sizeofWord))) - 1)))
				} else {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> (x & ((8 * (int(sizeofWord))) - 1))))
				}
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
				if (func() int {
					if float64(sum) > threshold*3*float64(max) {
						return 0
					}
					return 1
				}()) != 0 {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> (x & ((8 * (int(sizeofWord))) - 1)))
				} else {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> (x & ((8 * (int(sizeofWord))) - 1))))
				}
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
				bm.Map[int64(y)*int64(bm.Dy)+int64((i*8)/(8*(int(sizeofWord))))] |= (Word(b)) << Word(((int(sizeofWord))-1-i%(int(sizeofWord)))*8)
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
				if (func() int {
					if float64(b) > threshold*float64(max) {
						return 0
					}
					return 1
				}()) != 0 {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> (x & ((8 * (int(sizeofWord))) - 1)))
				} else {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> (x & ((8 * (int(sizeofWord))) - 1))))
				}
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
				if (func() int {
					if float64(sum) > threshold*3*float64(max) {
						return 0
					}
					return 1
				}()) != 0 {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> (x & ((8 * (int(sizeofWord))) - 1)))
				} else {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> (x & ((8 * (int(sizeofWord))) - 1))))
				}
			}
		}
	}
	bm_flip(bm)
	*bmp = bm
	return 0
eof:
	if bm_resize(bm, realheight) != 0 {
		goto std_error
	}
	bm_flip(bm)
	*bmp = bm
	return 1
format_error:
	bm_free(bm)
	if magic == '1' || magic == '4' {
		bm_read_error = libc.CString("invalid pbm file")
	} else if magic == '2' || magic == '5' {
		bm_read_error = libc.CString("invalid pgm file")
	} else {
		bm_read_error = libc.CString("invalid ppm file")
	}
	return -2
std_error:
	bm_free(bm)
	return -1
}

type bmp_info_s struct {
	FileSize        uint
	Reserved        uint
	DataOffset      uint
	InfoSize        uint
	W               uint
	H               uint
	Planes          uint
	Bits            uint
	Comp            uint
	ImageSize       uint
	XpixelsPerM     uint
	YpixelsPerM     uint
	Ncolors         uint
	ColorsImportant uint
	RedMask         uint
	GreenMask       uint
	BlueMask        uint
	AlphaMask       uint
	Ctbits          uint
	Topdown         int
}
type bmp_info_t bmp_info_s

var bmp_count int = 0
var bmp_pos int = 0

func bmp_readint(f *stdio.File, n int, p *uint) int {
	var (
		i   int
		sum uint = 0
		b   int
	)
	for i = 0; i < n; i++ {
		b = f.GetC()
		if b == stdio.EOF {
			return 1
		}
		sum += uint(b) << uint(i*8)
	}
	bmp_count += n
	bmp_pos += n
	*p = sum
	return 0
}
func bmp_pad_reset() {
	bmp_count = 0
}
func bmp_pad(f *stdio.File) int {
	var (
		c int
		i int
		b int
	)
	c = (-bmp_count) & 3
	for i = 0; i < c; i++ {
		b = f.GetC()
		if b == stdio.EOF {
			return 1
		}
	}
	bmp_pos += c
	bmp_count = 0
	return 0
}
func bmp_forward(f *stdio.File, pos int) int {
	var b int
	for bmp_pos < pos {
		b = f.GetC()
		if b == stdio.EOF {
			return 1
		}
		bmp_pos++
		bmp_count++
	}
	return 0
}
func bm_readbody_bmp(f *stdio.File, threshold float64, bmp **Bitmap) int {
	var (
		bmpinfo    bmp_info_t
		coltable   *int
		b          uint
		c          uint
		i          uint
		bm         *Bitmap
		mask       int
		x          uint
		y          uint
		col        [2]int
		bitbuf     uint
		n          uint
		redshift   uint
		greenshift uint
		blueshift  uint
		col1       [2]int
		realheight int
	)
	bm_read_error = nil
	bm = nil
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
		bm_read_error = libc.CString("cannot handle bmp planes")
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
			if float64(c) > threshold*3*math.MaxUint8 {
				*(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr(i))) = 0
			} else {
				*(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr(i))) = 1
			}
			if i < 2 {
				col1[i] = int(c)
			}
		}
	}
	if bmpinfo.InfoSize != 12 {
		if bmp_forward(f, int(bmpinfo.DataOffset)) != 0 {
			goto try_error
		}
	}
	bm = NewBitmap(int(bmpinfo.W), int(bmpinfo.H))
	if bm == nil {
		goto std_error
	}
	realheight = 0
	switch bmpinfo.Bits + bmpinfo.Comp*0x100 {
	default:
		goto format_error
	case 0x1:
		if col1[0] < col1[1] {
			mask = math.MaxUint8
		} else {
			mask = 0
		}
		for y = 0; y < bmpinfo.H; y++ {
			realheight = int(y + 1)
			bmp_pad_reset()
			for i = 0; i*8 < bmpinfo.W; i++ {
				if bmp_readint(f, 1, &b) != 0 {
					goto eof
				}
				b ^= uint(mask)
				bm.Map[int64(y)*int64(bm.Dy)+int64((i*8)/uint(8*(int(sizeofWord))))] |= (Word(b)) << Word(((int(sizeofWord))-1-int(i%uint(int(sizeofWord))))*8)
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
				if (func() int {
					if b < bmpinfo.Ncolors {
						return *(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr(b)))
					}
					return 0
				}()) != 0 {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1)))
				} else {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1))))
				}
			}
			if bmp_pad(f) != 0 {
				goto try_error
			}
		}
	case 0x10:
		bm_read_error = libc.CString("cannot handle bmp 16-bit coding")
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
				if (func() int {
					if float64(c) > threshold*3*math.MaxUint8 {
						return 0
					}
					return 1
				}()) != 0 {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1)))
				} else {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1))))
				}
			}
			if bmp_pad(f) != 0 {
				goto try_error
			}
		}
	case 0x320:
		if bmpinfo.RedMask == 0 || bmpinfo.GreenMask == 0 || bmpinfo.BlueMask == 0 {
			goto format_error
		}
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
				if (func() int {
					if float64(c) > threshold*3*math.MaxUint8 {
						return 0
					}
					return 1
				}()) != 0 {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1)))
				} else {
					bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1))))
				}
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
					if int(x) >= 0 && int(x) < bm.W && (int(y) >= 0 && int(y) < bm.H) {
						if (col[i&1]) != 0 {
							bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1)))
						} else {
							bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1))))
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
					if int(x) >= 0 && int(x) < bm.W && (int(y) >= 0 && int(y) < bm.H) {
						if (func() int {
							if ((b >> (4 - (i&1)*4)) & 0xF) < bmpinfo.Ncolors {
								return *(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr((b>>(4-(i&1)*4))&0xF)))
							}
							return 0
						}()) != 0 {
							bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1)))
						} else {
							bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1))))
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
					if int(x) >= 0 && int(x) < bm.W && (int(y) >= 0 && int(y) < bm.H) {
						if (func() int {
							if c < bmpinfo.Ncolors {
								return *(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr(c)))
							}
							return 0
						}()) != 0 {
							bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1)))
						} else {
							bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1))))
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
					if int(x) >= 0 && int(x) < bm.W && (int(y) >= 0 && int(y) < bm.H) {
						if (func() int {
							if b < bmpinfo.Ncolors {
								return *(*int)(unsafe.Add(unsafe.Pointer(coltable), unsafe.Sizeof(int(0))*uintptr(b)))
							}
							return 0
						}()) != 0 {
							bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1)))
						} else {
							bm.Map[int64(y)*int64(bm.Dy)+int64(x/uint(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> int(x&uint((8*(int(sizeofWord)))-1))))
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
		bm_flip(bm)
	}
	*bmp = bm
	return 0
eof:
	if bm_resize(bm, realheight) != 0 {
		goto std_error
	}

	if bmpinfo.Topdown != 0 {
		bm_flip(bm)
	}
	*bmp = bm
	return 1
format_error:
try_error:
	;

	bm_free(bm)
	if bm_read_error == nil {
		bm_read_error = libc.CString("invalid bmp file")
	}
	return -2
std_error:

	bm_free(bm)
	return -1
}
func bm_writepbm(f *stdio.File, bm *Bitmap) {
	var (
		w   int
		h   int
		bpr int
		y   int
		i   int
		c   int
	)
	w = bm.W
	h = bm.H
	bpr = (w + 7) / 8
	stdio.Fprintf(f, "P4\n%d %d\n", w, h)
	for y = h - 1; y >= 0; y-- {
		for i = 0; i < bpr; i++ {
			c = int((bm.Map[int64(y)*int64(bm.Dy)+int64((i*8)/(8*(int(sizeofWord))))] >> Word(((int(sizeofWord))-1-i%(int(sizeofWord)))*8)) & math.MaxUint8)
			f.PutC(c)
		}
	}
	return
}
func bm_print(f *stdio.File, bm *Bitmap) int {
	var (
		x  int
		y  int
		xx int
		yy int
		d  int
		sw int
		sh int
	)
	if bm.W < 79 {
		sw = bm.W
	} else {
		sw = 79
	}
	if bm.W < 79 {
		sh = bm.H
	} else {
		sh = bm.H * sw * 44 / (bm.W * 79)
	}
	for yy = sh - 1; yy >= 0; yy-- {
		for xx = 0; xx < sw; xx++ {
			d = 0
			for x = xx * bm.W / sw; x < (xx+1)*bm.W/sw; x++ {
				for y = yy * bm.H / sh; y < (yy+1)*bm.H/sh; y++ {
					if func() bool {
						if x >= 0 && x < bm.W && (y >= 0 && y < bm.H) {
							return (bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] & Word((1<<((8*(int(sizeofWord)))-1))>>(x&((8*(int(sizeofWord)))-1)))) != 0
						}
						return false
					}() {
						d++
					}
				}
			}
			f.PutC(func() int {
				if d != 0 {
					return '*'
				}
				return ' '
			}())
		}
		f.PutC('\n')
	}
	return 0
}
