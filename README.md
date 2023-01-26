# GoTrace

Pure Go implementation of [Potrace](https://potrace.sourceforge.net/potracelib.pdf) vectorization library.
Supports simple SVG and PDF output generation.

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

## Library
```
go get -u github.com/gotranspile/gotrace
```
