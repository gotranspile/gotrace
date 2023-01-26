package gotrace

const DIM_IN = 72
const DIM_CM = 72
const DIM_MM = 72
const DIM_PT = 1
const DEFAULT_DIM = 72
const DEFAULT_DIM_NAME = "inches"
const DEFAULT_PAPERWIDTH = 612
const DEFAULT_PAPERHEIGHT = 792
const DEFAULT_PAPERFORMAT = "letter"

type Dim struct {
	X float64
	D float64
}
type BackendInfo struct {
	Backend      *backend_s
	Param        *Param
	Debug        bool
	Width_d      Dim
	Height_d     Dim
	Rx           float64
	Ry           float64
	Sx           float64
	Sy           float64
	Stretch      float64
	Lmar_d       Dim
	Rmar_d       Dim
	Tmar_d       Dim
	Bmar_d       Dim
	Angle        float64
	Paperwidth   int
	Paperheight  int
	Tight        int
	Unit         float64
	Compress     int
	Pslevel      int
	Color        int
	Fillcolor    int
	Gamma        float64
	Longcoding   int
	Outfile      *byte
	Infiles      **byte
	Infilecount  int
	Some_infiles int
	Blacklevel   float64
	Invert       int
	Opaque       bool
	Grouping     int
	Progress     int
	Progress_bar *progress_bar_t
}
type ImgInfo struct {
	Pixwidth  int
	Pixheight int
	Width     float64
	Height    float64
	Lmar      float64
	Rmar      float64
	Tmar      float64
	Bmar      float64
	Trans     Trans
}
