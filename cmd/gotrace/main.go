package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gotranspile/gotrace"
)

var Root = &cobra.Command{
	Use:  "gotrace [options] [filename...]",
	RunE: run,
}

var (
	fOut     *string
	fBackend *string
	fSVG     *bool
)

func main() {
	fOut = Root.Flags().StringP("output", "o", "", "write all output to this file")
	fBackend = Root.Flags().StringP("backend", "b", "", "select backend by name")
	fSVG = Root.Flags().BoolP("svg", "s", false, "SVG backend (scalable vector graphics)")
	if err := Root.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if *fBackend == "" {
		if *fSVG {
			*fBackend = "svg"
		} else {
			*fBackend = "svg" // default
		}
	}
	conf := gotrace.DefaultConfig()
	if len(args) == 0 {
		return process(*fOut, conf, "-", os.Stdin)
	} else if len(args) == 1 {
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close()
		return process(*fOut, conf, args[0], f)
	}
	for _, fname := range args {
		f, err := os.Open(fname)
		if err != nil {
			return err
		}
		err = process(*fOut, conf, fname, f)
		_ = f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func process(out string, conf *gotrace.Config, fname string, r io.Reader) error {
	bm, err := gotrace.BitmapRead(r, 0.5)
	if err != nil {
		return err
	}
	paths, err := gotrace.Trace(bm, conf)
	if err != nil {
		return err
	}
	c := gotrace.NewRenderConf()
	if out == "-" || (out == "" && fname == "-") {
		return gotrace.Render(*fBackend, c, os.Stdout, paths, bm.W, bm.H)
	}
	if out == "" {
		ext := filepath.Ext(fname)
		out = strings.TrimSuffix(fname, ext)
		switch *fBackend {
		case "svg", "gimppath":
			out += ".svg"
		case "pdf", "pdfpage":
			out += ".pdf"
		}
	}
	return gotrace.RenderFile(*fBackend, c, out, paths, bm.W, bm.H)
}
