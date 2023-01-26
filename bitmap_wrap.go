package gotrace

import (
	"fmt"
	"io"
	"os"

	"github.com/gotranspile/cxgo/runtime/stdio"
)

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
