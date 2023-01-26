package gotrace

import "errors"

var (
	ErrIncomplete = errors.New("tracing incomplete")
)

// Trace a bitmap.
func Trace(bm *Bitmap, conf *Config) (*Path, error) {
	if conf == nil {
		conf = DefaultConfig()
	}
	st := traceBitmap(conf, bm)
	if st == nil {
		return nil, errors.New("tracing failed")
	}
	switch st.Status {
	case statusOK:
		return st.Plist, nil
	case statusIncomplete:
		return st.Plist, ErrIncomplete
	default:
		return st.Plist, errors.New("tracing failed")
	}
}
