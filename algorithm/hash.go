package algorithm

// DJB算法
func DJBHash(b []byte) int {
	var hash int = 5381
	for i, _ := range b {
		hash = ((hash << 5) + hash) + int(b[i])
	}
	return hash & 0x7FFFFFFF
}
