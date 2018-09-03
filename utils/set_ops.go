package utils

func Combine(v1, v2 []int64) []int64 {

	m := make(map[int64]bool)
	for _, v := range v1 {
		m[v] = true
	}

	for _, v := range v2 {
		if _, ok := m[v]; ok {
			m[v] = true
		}
	}

	ret := make([]int64, 0)
	for k := range m {
		if m[k] {
			ret = append(ret, k)
		}
	}

	return ret
}
