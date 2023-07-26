package gotrace

func dpoint(p Point) DPoint {
	var res DPoint
	res.X = float64(p.X)
	res.Y = float64(p.Y)
	return res
}
func aux_interval(lambda float64, a DPoint, b DPoint) DPoint {
	var res DPoint
	res.X = a.X + lambda*(b.X-a.X)
	res.Y = a.Y + lambda*(b.Y-a.Y)
	return res
}
func mod(a int, n int) int {
	if a >= n {
		return a % n
	}
	if a >= 0 {
		return a
	}
	return n - 1 - (int(-1-a))%n
}
func floordiv(a int, n int) int {
	if a >= 0 {
		return a / n
	}
	return int(-1 - (int(-1-a))/n)
}
