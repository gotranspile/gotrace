package gotrace_test

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/gotranspile/gotrace"
)

func ExamplePotracePNG() {
	f, err := os.Open(filepath.Join(testdata, "stanford.png"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	bm := gotrace.BitmapFromImage(img, nil)

	paths, err := gotrace.Trace(bm, nil)
	if err != nil {
		panic(err)
	}

	sz := img.Bounds().Size()
	out := filepath.Join(testdata, "stanford.svg.png")
	err = gotrace.RenderFile("svg", nil, out, paths, sz.X, sz.Y)
	if err != nil {
		panic(err)
	}
	fmt.Println(hashFile(out))
	// Output: 9aee25cf092f7228c4f14ba4607592487f4fab07
}
