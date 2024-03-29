// Code generated by cxgo. DO NOT EDIT.

package gotrace

const DefaultPaperWidth = 612
const DefaultPaperHeight = 792

type Dim struct {
	X float64
	D float64
}
type RenderConf struct {
	Backend *BackendInfo

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
	Compress     bool
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
}
type imgInfo struct {
	Pixwidth  int
	Pixheight int
	Width     float64
	Height    float64
	Lmar      float64
	Rmar      float64
	Tmar      float64
	Bmar      float64
	Trans     transT
}

func double_of_dim(d Dim, def float64) float64 {
	if d.D != 0 {
		return d.X * d.D
	} else {
		return d.X * def
	}
}

func calc_dimensions(info *RenderConf, imginfo *imgInfo, plist *Path) {
	var (
		dim_def         float64
		maxwidth        float64
		maxheight       float64
		sc              float64
		default_scaling int = 0
	)
	if imginfo.Pixwidth == 0 {
		imginfo.Pixwidth = 1
	}
	if imginfo.Pixheight == 0 {
		imginfo.Pixheight = 1
	}
	if info.Backend.Pixel {
		dim_def = 1
	} else {
		dim_def = 72
	}
	if info.Width_d.X == undef {
		imginfo.Width = undef
	} else {
		imginfo.Width = double_of_dim(info.Width_d, dim_def)
	}
	if info.Height_d.X == undef {
		imginfo.Height = undef
	} else {
		imginfo.Height = double_of_dim(info.Height_d, dim_def)
	}
	if info.Lmar_d.X == undef {
		imginfo.Lmar = undef
	} else {
		imginfo.Lmar = double_of_dim(info.Lmar_d, dim_def)
	}
	if info.Rmar_d.X == undef {
		imginfo.Rmar = undef
	} else {
		imginfo.Rmar = double_of_dim(info.Rmar_d, dim_def)
	}
	if info.Tmar_d.X == undef {
		imginfo.Tmar = undef
	} else {
		imginfo.Tmar = double_of_dim(info.Tmar_d, dim_def)
	}
	if info.Bmar_d.X == undef {
		imginfo.Bmar = undef
	} else {
		imginfo.Bmar = double_of_dim(info.Bmar_d, dim_def)
	}
	trans_from_rect(&imginfo.Trans, float64(imginfo.Pixwidth), float64(imginfo.Pixheight))
	if info.Tight != 0 {
		trans_tighten(&imginfo.Trans, plist)
	}
	if info.Backend.Pixel {
		if imginfo.Width == undef && info.Sx != undef {
			imginfo.Width = imginfo.Trans.Bb[0] * info.Sx
		}
		if imginfo.Height == undef && info.Sy != undef {
			imginfo.Height = imginfo.Trans.Bb[1] * info.Sy
		}
	} else {
		if imginfo.Width == undef && info.Rx != undef {
			imginfo.Width = imginfo.Trans.Bb[0] / info.Rx * 72
		}
		if imginfo.Height == undef && info.Ry != undef {
			imginfo.Height = imginfo.Trans.Bb[1] / info.Ry * 72
		}
	}
	if imginfo.Width == undef && imginfo.Height != undef {
		imginfo.Width = imginfo.Height / imginfo.Trans.Bb[1] * imginfo.Trans.Bb[0] / info.Stretch
	} else if imginfo.Width != undef && imginfo.Height == undef {
		imginfo.Height = imginfo.Width / imginfo.Trans.Bb[0] * imginfo.Trans.Bb[1] * info.Stretch
	}
	if imginfo.Width == undef && imginfo.Height == undef {
		imginfo.Width = imginfo.Trans.Bb[0]
		imginfo.Height = imginfo.Trans.Bb[1] * info.Stretch
		default_scaling = 1
	}
	trans_scale_to_size(&imginfo.Trans, imginfo.Width, imginfo.Height)
	if info.Angle != 0.0 {
		trans_rotate(&imginfo.Trans, info.Angle)
		if info.Tight != 0 {
			trans_tighten(&imginfo.Trans, plist)
		}
	}
	if default_scaling != 0 && info.Backend.Fixed {
		maxwidth = undef
		maxheight = undef
		if imginfo.Lmar != undef && imginfo.Rmar != undef {
			maxwidth = float64(info.Paperwidth) - imginfo.Lmar - imginfo.Rmar
		}
		if imginfo.Bmar != undef && imginfo.Tmar != undef {
			maxheight = float64(info.Paperheight) - imginfo.Bmar - imginfo.Tmar
		}
		if maxwidth == undef && maxheight == undef {
			if float64(info.Paperwidth-144) > (float64(info.Paperwidth) * 0.75) {
				maxwidth = float64(info.Paperwidth - 144)
			} else {
				maxwidth = float64(info.Paperwidth) * 0.75
			}
			if float64(info.Paperheight-144) > (float64(info.Paperheight) * 0.75) {
				maxheight = float64(info.Paperheight - 144)
			} else {
				maxheight = float64(info.Paperheight) * 0.75
			}
		}
		if maxwidth == undef {
			sc = maxheight / imginfo.Trans.Bb[1]
		} else if maxheight == undef {
			sc = maxwidth / imginfo.Trans.Bb[0]
		} else {
			if (maxwidth / imginfo.Trans.Bb[0]) < (maxheight / imginfo.Trans.Bb[1]) {
				sc = maxwidth / imginfo.Trans.Bb[0]
			} else {
				sc = maxheight / imginfo.Trans.Bb[1]
			}
		}
		imginfo.Width *= sc
		imginfo.Height *= sc
		trans_rescale(&imginfo.Trans, sc)
	}
	if info.Backend.Fixed {
		if imginfo.Lmar == undef && imginfo.Rmar == undef {
			imginfo.Lmar = (float64(info.Paperwidth) - imginfo.Trans.Bb[0]) / 2
		} else if imginfo.Lmar == undef {
			imginfo.Lmar = float64(info.Paperwidth) - imginfo.Trans.Bb[0] - imginfo.Rmar
		} else if imginfo.Lmar != undef && imginfo.Rmar != undef {
			imginfo.Lmar += (float64(info.Paperwidth) - imginfo.Trans.Bb[0] - imginfo.Lmar - imginfo.Rmar) / 2
		}
		if imginfo.Bmar == undef && imginfo.Tmar == undef {
			imginfo.Bmar = (float64(info.Paperheight) - imginfo.Trans.Bb[1]) / 2
		} else if imginfo.Bmar == undef {
			imginfo.Bmar = float64(info.Paperheight) - imginfo.Trans.Bb[1] - imginfo.Tmar
		} else if imginfo.Bmar != undef && imginfo.Tmar != undef {
			imginfo.Bmar += (float64(info.Paperheight) - imginfo.Trans.Bb[1] - imginfo.Bmar - imginfo.Tmar) / 2
		}
	} else {
		if imginfo.Lmar == undef {
			imginfo.Lmar = 0
		}
		if imginfo.Rmar == undef {
			imginfo.Rmar = 0
		}
		if imginfo.Bmar == undef {
			imginfo.Bmar = 0
		}
		if imginfo.Tmar == undef {
			imginfo.Tmar = 0
		}
	}
}
