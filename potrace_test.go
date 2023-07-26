package gotrace_test

import (
	"crypto/sha1"
	"encoding/hex"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gotranspile/gotrace"
)

func checkErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

const testdata = "./testdata/"

func TestPotrace(t *testing.T) {
	bm, err := gotrace.BitmapReadFile(filepath.Join(testdata, "stanford.pbm"), 0.5)
	checkErr(t, err)

	plist, err := gotrace.Trace(bm, nil)
	checkErr(t, err)

	// Test data generated with:
	//		potrace -b <backend> -o ./testdata/stanford.<ext> ./testdata/stanford.pbm
	//		sha1sum ./testdata/*

	bi := gotrace.NewRenderConf()
	fname := filepath.Join(testdata, "stanford.svg")
	err = gotrace.RenderFile("svg", bi, fname, plist, bm.W, bm.H)
	checkErr(t, err)
	if h := hashFile(fname); h != "9aee25cf092f7228c4f14ba4607592487f4fab07" {
		t.Errorf("unexpected hash for SVG: %s", h)
	}

	fname = filepath.Join(testdata, "stanford.pdf")
	err = gotrace.RenderFile("pdf", bi, fname, plist, bm.W, bm.H)
	checkErr(t, err)
	if h := hashFile(fname); h != "347917e85988fc7f881244cb2f8adfa7ce04d495" {
		t.Errorf("unexpected hash for PDF: %s", h)
	}

	fname = filepath.Join(testdata, "stanford.dxf")
	err = gotrace.RenderFile("dxf", bi, fname, plist, bm.W, bm.H)
	checkErr(t, err)
	if h := hashFile(fname); h != "8c03dff3ce9f41a9b726963b7c6dda4f3dfffbab" {
		t.Errorf("unexpected hash for DXF: %s", h)
	}
}

func TestPotracePNG(t *testing.T) {
	f, err := os.Open(filepath.Join(testdata, "stanford.png"))
	checkErr(t, err)
	defer f.Close()

	img, err := png.Decode(f)
	checkErr(t, err)

	bm := gotrace.BitmapFromImage(img, nil)

	plist, err := gotrace.Trace(bm, nil)
	checkErr(t, err)

	bi := gotrace.NewRenderConf()
	fname := filepath.Join(testdata, "stanford.png.svg")
	err = gotrace.RenderFile("svg", bi, fname, plist, bm.W, bm.H)
	checkErr(t, err)
	if h := hashFile(fname); h != "9aee25cf092f7228c4f14ba4607592487f4fab07" {
		t.Errorf("unexpected hash for SVG: %s", h)
	}
}

func hashFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}
