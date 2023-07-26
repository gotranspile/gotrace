package gotrace

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/gotranspile/cxgo/runtime/stdio"
)

type backendStateFunc func(info *RenderConf, fout *stdio.File) int
type backendRenderFunc func(info *RenderConf, fout *stdio.File, plist *Path, imginfo *imgInfo) int

type backend struct {
	Name   string
	Info   BackendInfo
	Init   backendStateFunc
	Render backendRenderFunc
	Term   backendStateFunc
}

var backends = make(map[string]backend)

func registerBackend(name string, info BackendInfo, init backendStateFunc, render backendRenderFunc, term backendStateFunc) {
	if _, ok := backends[name]; ok {
		panic("already registered: " + name)
	}
	backends[name] = backend{
		Name:   name,
		Info:   info,
		Init:   init,
		Render: render,
		Term:   term,
	}
}

// Backends lists names of supported backends.
func Backends() []string {
	var names []string
	for k := range backends {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

func init() {
	registerBackend("svg", BackendInfo{Fixed: false, Pixel: false}, nil, page_svg, nil)
	registerBackend("gimppath", BackendInfo{Fixed: false, Pixel: true}, nil, page_gimp, nil)
	registerBackend("pdf", BackendInfo{Fixed: false, Pixel: false}, init_pdf, page_pdf, term_pdf)
	registerBackend("pdfpage", BackendInfo{Fixed: true, Pixel: false}, init_pdf, page_pdfpage, term_pdf)
	registerBackend("dxf", BackendInfo{Fixed: false, Pixel: true}, nil, page_dxf, nil)
}

const undef = 1e30

type BackendInfo struct {
	Fixed bool
	Pixel bool
}

// NewRenderConf creates a default render config.
func NewRenderConf() *RenderConf {
	return &RenderConf{
		Width_d:  Dim{X: undef},
		Height_d: Dim{X: undef},
		Rx:       undef, Ry: undef,
		Sx: undef, Sy: undef,
		Stretch:     1,
		Lmar_d:      Dim{X: undef},
		Rmar_d:      Dim{X: undef},
		Tmar_d:      Dim{X: undef},
		Bmar_d:      Dim{X: undef},
		Paperwidth:  DefaultPaperWidth,
		Paperheight: DefaultPaperHeight,
		Unit:        10,
		Compress:    false,
		Pslevel:     2,
		Color:       0x000000,
		Gamma:       2.2,
		Backend: &BackendInfo{
			Pixel: false,
			Fixed: false,
		},
		Blacklevel: 0.5,
		Grouping:   1,
		Fillcolor:  0xffffff,
	}
}

// RenderFile writes paths with a given backend to a file.
func RenderFile(backend string, conf *RenderConf, out string, paths *Path, width, height int) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = Render(backend, conf, f, paths, width, height); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}

// Render paths with a given backend.
func Render(backend string, conf *RenderConf, out io.Writer, paths *Path, width, height int) error {
	if f, ok := out.(*os.File); ok {
		return renderFile(backend, conf, f, paths, width, height)
	}
	pr, pw, err := os.Pipe()
	if err != nil {
		return err
	}
	errc := make(chan error, 1)
	go func() {
		defer close(errc)
		defer pr.Close()
		if _, err := io.Copy(out, pr); err != nil {
			errc <- err
		}
	}()
	defer pw.Close()
	if err = renderFile(backend, conf, pw, paths, width, height); err != nil {
		return err
	}
	if err = pw.Close(); err != nil {
		return err
	}
	return <-errc
}

func renderFile(backend string, conf *RenderConf, out *os.File, paths *Path, width, height int) error {
	b, ok := backends[backend]
	if !ok {
		return fmt.Errorf("unsupported backend: %q", backend)
	}
	if conf == nil {
		conf = NewRenderConf()
		info := b.Info
		conf.Backend = &info
	} else if conf.Backend == nil {
		info := b.Info
		conf.Backend = &info
	}

	cout := stdio.OpenFrom(out)
	var iinfo imgInfo
	iinfo.Pixwidth = width
	iinfo.Pixheight = height
	calc_dimensions(conf, &iinfo, paths)
	if b.Init != nil {
		b.Init(conf, cout)
	}
	b.Render(conf, cout, paths, &iinfo)
	if b.Term != nil {
		b.Term(conf, cout)
	}
	return nil
}
