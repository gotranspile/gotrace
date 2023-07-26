# GoTrace

Pure Go implementation of [Potrace](https://potrace.sourceforge.net/potracelib.pdf) vectorization library.
Supports SVG, PDF and DXF output generation.

It is a direct machine translation (transpilation) of potrace using [cxgo](https://github.com/gotranspile/cxgo).

**Original image**

![Original](http://potrace.sourceforge.net/img/stanford-orig2.png)

**Vectorized image**

![Vectorized](http://potrace.sourceforge.net/img/stanford-smooth2.png)

# Installation

## Tool
```
go install github.com/gotranspile/gotrace@latest
```

### Usage

Convert PNG image to SVG:
```
gotrace -s -o ./testdata/stanford.svg ./testdata/stanford.png
```

## Library
```
go get -u github.com/gotranspile/gotrace
```

### Usage

Minimal example:

```go
func traceImage(outPath string, img image.Image) error {
    bm := gotrace.BitmapFromImage(img, nil)
    
    paths, err := gotrace.Trace(bm, nil)
    if err != nil {
        return err
    }
    
    sz := img.Bounds().Size()
    return gotrace.RenderFile("svg", nil, outPath, paths, sz.X, sz.Y)
}
```

For a full example, see [example_test.go](./example_test.go).

## Updating the source

This library uses `cxgo` to translate C source code directly to Go. See [cxgo.yml](./cxgo.yml) for the config.

To regenerate source, install `cxgo` and `goimports` and run:

```
go generate
```