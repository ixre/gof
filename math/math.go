package math

import (
	"math"
)

// 四舍五入计算,并保留n位精度位
// 如果n小于０，则四舍五入到整数位
func Round(val float64, n int) float64 {
	if n <= 0 {
		if val < 0 {
			return math.Ceil(val - 0.5)
		}
		return math.Floor(val + 0.5)
	}
	digit := math.Pow10(n)
	if val < 0 {
		return math.Ceil(val*digit-0.5) / digit
	}
	return math.Floor(val*digit+0.5) / digit
}

// 四舍五入计算,并保留n位精度位
func Round32(val float32, n int) float32 {
	return float32(Round(float64(val), n))
}

// 普通近似值计算, 不四舍五入,n为小数点精度
func FixedFloat(v float64, n int) float64 {
	return math.Floor(v*math.Pow10(n)) / math.Pow10(n)
}
