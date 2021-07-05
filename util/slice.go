package util

// InArrayString 查询某一个值是否在切片中
func InArrayString(search string, slice []string) bool {
	return inArray(search, slice)
}

// InArrayInt 查询某一个值是否在切片中
func InArrayInt(search int, slice []int) bool {
	return inArray(search, slice)
}

// InArrayInt64 查询某一个值是否在切片中
func InArrayInt64(search int64, slice []int64) bool {
	return inArray(search, slice)
}

func inArray(needle interface{}, haystacks interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range haystacks.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range haystacks.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range haystacks.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}
