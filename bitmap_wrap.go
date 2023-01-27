package gotrace

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"os"

	"github.com/gotranspile/cxgo/runtime/stdio"
)

func BitmapFromGray(img *image.Gray, fnc func(c color.Gray) bool) *Bitmap {
	if fnc == nil {
		fnc = func(c color.Gray) bool {
			return c.Y < math.MaxUint8/2
		}
	}
	sz := img.Rect.Size()
	bm := NewBitmap(sz.X, sz.Y)
	for y := 0; y < sz.Y; y++ {
		for x := 0; x < sz.X; x++ {
			bm.Put(x, sz.Y-1-y, fnc(img.GrayAt(img.Rect.Min.X+x, img.Rect.Min.Y+y)))
		}
	}
	return bm
}

func BitmapFromGray16(img *image.Gray16, fnc func(c color.Gray16) bool) *Bitmap {
	if fnc == nil {
		fnc = func(c color.Gray16) bool {
			return c.Y < math.MaxUint16/2
		}
	}
	sz := img.Rect.Size()
	bm := NewBitmap(sz.X, sz.Y)
	for y := 0; y < sz.Y; y++ {
		for x := 0; x < sz.X; x++ {
			bm.Put(x, sz.Y-1-y, fnc(img.Gray16At(img.Rect.Min.X+x, img.Rect.Min.Y+y)))
		}
	}
	return bm
}

func BitmapFromNRGBA(img *image.NRGBA, fnc func(c color.NRGBA) bool) *Bitmap {
	if fnc == nil {
		fnc = func(c color.NRGBA) bool {
			if c.A < 128 {
				return true
			}
			return c.R+c.G+c.B < 128
		}
	}
	sz := img.Rect.Size()
	bm := NewBitmap(sz.X, sz.Y)
	for y := 0; y < sz.Y; y++ {
		for x := 0; x < sz.X; x++ {
			bm.Put(x, sz.Y-1-y, fnc(img.NRGBAAt(img.Rect.Min.X+x, img.Rect.Min.Y+y)))
		}
	}
	return bm
}

func BitmapFromRGBA(img *image.RGBA, fnc func(c color.RGBA) bool) *Bitmap {
	if fnc == nil {
		fnc = func(c color.RGBA) bool {
			if c.A < 128 {
				return true
			}
			return c.R+c.G+c.B < 128
		}
	}
	sz := img.Rect.Size()
	bm := NewBitmap(sz.X, sz.Y)
	for y := 0; y < sz.Y; y++ {
		for x := 0; x < sz.X; x++ {
			bm.Put(x, sz.Y-1-y, fnc(img.RGBAAt(img.Rect.Min.X+x, img.Rect.Min.Y+y)))
		}
	}
	return bm
}

func BitmapFromImage(img image.Image, fnc func(c color.Color) bool) *Bitmap {
	switch img := img.(type) {
	case *image.Gray:
		if fnc == nil {
			return BitmapFromGray(img, nil)
		}
		return BitmapFromGray(img, func(c color.Gray) bool {
			return fnc(c)
		})
	case *image.Gray16:
		if fnc == nil {
			return BitmapFromGray16(img, nil)
		}
		return BitmapFromGray16(img, func(c color.Gray16) bool {
			return fnc(c)
		})
	case *image.NRGBA:
		if fnc == nil {
			return BitmapFromNRGBA(img, nil)
		}
		return BitmapFromNRGBA(img, func(c color.NRGBA) bool {
			return fnc(c)
		})
	case *image.RGBA:
		if fnc == nil {
			return BitmapFromRGBA(img, nil)
		}
		return BitmapFromRGBA(img, func(c color.RGBA) bool {
			return fnc(c)
		})
	}
	if fnc == nil {
		fnc = func(c color.Color) bool {
			r, g, b, a := c.RGBA()
			if a < 128 {
				return true
			}
			return r+g+b < 128
		}
	}
	rect := img.Bounds()
	sz := rect.Size()
	bm := NewBitmap(sz.X, sz.Y)
	for y := 0; y < sz.Y; y++ {
		for x := 0; x < sz.X; x++ {
			bm.Put(x, sz.Y-1-y, fnc(img.At(rect.Min.X+x, rect.Min.Y+y)))
		}
	}
	return bm
}

func BitmapReadFile(path string, threshold float64) (*Bitmap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return bitmapReadFile(f, threshold)
}

func BitmapRead(r io.Reader, threshold float64) (*Bitmap, error) {
	if f, ok := r.(*os.File); ok {
		return bitmapReadFile(f, threshold)
	}
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	errc := make(chan error, 1)
	go func() {
		defer close(errc)
		defer pw.Close()
		if _, err := io.Copy(pw, r); err != nil {
			errc <- err
		}
	}()
	defer pr.Close()
	bm, err := bitmapReadFile(pr, threshold)
	if err != nil {
		return nil, err
	}
	err = <-errc
	return bm, err
}

func bitmapReadFile(f *os.File, threshold float64) (*Bitmap, error) {
	var bm *Bitmap
	e := bitmapRead(stdio.OpenFrom(f), threshold, &bm)
	if e != 0 {
		return nil, fmt.Errorf("error code: %d", e)
	}
	return bm, nil
}
