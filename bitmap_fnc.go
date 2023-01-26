package gotrace

func (bm *Bitmap) Get(x int, y int) bool {
	if x >= 0 && x < bm.W && (y >= 0 && y < bm.H) {
		return (bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] & Word((1<<((8*(int(sizeofWord)))-1))>>(x&((8*(int(sizeofWord)))-1)))) != 0
	}
	return false
}
func (bm *Bitmap) Put(x int, y int, b bool) {
	if x >= 0 && x < bm.W && (y >= 0 && y < bm.H) {
		if b {
			bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] |= Word((1 << ((8 * (int(sizeofWord))) - 1)) >> (x & ((8 * (int(sizeofWord))) - 1)))
		} else {
			bm.Map[int64(y)*int64(bm.Dy)+int64(x/(8*(int(sizeofWord))))] &= Word(^((1 << ((8 * (int(sizeofWord))) - 1)) >> (x & ((8 * (int(sizeofWord))) - 1))))
		}
	} else {
	}
}
