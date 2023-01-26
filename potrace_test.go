package gotrace_test

import (
	"crypto/sha1"
	"encoding/hex"
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

func TestPotrace(t *testing.T) {
	const dir = "./testdata/"

	bm, err := gotrace.BitmapReadFile(filepath.Join(dir, "stanford.pbm"), 0.5)
	checkErr(t, err)

	plist, err := gotrace.Trace(bm, nil)
	checkErr(t, err)

	bi := gotrace.NewRenderConf()
	fname := filepath.Join(dir, "stanford.svg")
	err = gotrace.RenderFile("svg", bi, fname, plist, bm.W, bm.H)
	checkErr(t, err)
	if h := hashFile(t, fname); h != "9aee25cf092f7228c4f14ba4607592487f4fab07" {
		t.Errorf("unexpected hash for SVG: %s", h)
	}

	fname = filepath.Join(dir, "stanford.pdf")
	err = gotrace.RenderFile("pdf", bi, fname, plist, bm.W, bm.H)
	checkErr(t, err)
}

func hashFile(t testing.TB, path string) string {
	f, err := os.Open(path)
	checkErr(t, err)
	defer f.Close()
	h := sha1.New()
	_, err = io.Copy(h, f)
	checkErr(t, err)
	return hex.EncodeToString(h.Sum(nil))
}
