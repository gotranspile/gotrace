package gotrace

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gotranspile/cxgo/runtime/stdio"
)

func checkErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPotrace(t *testing.T) {
	resp, err := http.Get("https://potrace.sourceforge.net/img/stanford.pbm")
	checkErr(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatal(resp.Status)
	}

	tfile, err := os.CreateTemp("", "potrace_")
	checkErr(t, err)
	defer func() {
		_ = tfile.Close()
		_ = os.Remove(tfile.Name())
	}()

	_, err = io.Copy(tfile, resp.Body)
	checkErr(t, err)
	tfile.Seek(0, io.SeekStart)

	var bm *Bitmap
	e := BitmapRead(stdio.OpenFrom(tfile), 0.5, &bm)
	if e != 0 {
		t.Fatal(e)
	}

	p := ParamDefault()
	st := Trace(p, bm)
	if st.Status != 0 {
		t.Fatal(st.Status)
	}
	var tr Trans
	trans_from_rect(&tr, float64(bm.W), float64(bm.H))
	bi := &BackendInfo{
		Unit: 1,
	}
	iinfo := &ImgInfo{
		Pixwidth: bm.W, Pixheight: bm.H,
		Width: float64(bm.W), Height: float64(bm.H),
		Trans: tr,
	}

	svgOut, err := os.Create("stanford.svg")
	checkErr(t, err)
	defer svgOut.Close()
	page_svg(bi, stdio.OpenFrom(svgOut), st.Plist, iinfo)

	pdfOut, err := os.Create("stanford.pdf")
	checkErr(t, err)
	defer pdfOut.Close()
	init_pdf(bi, stdio.OpenFrom(pdfOut))
	page_pdf(bi, stdio.OpenFrom(pdfOut), st.Plist, iinfo)
	term_pdf(bi, stdio.OpenFrom(pdfOut))
}
