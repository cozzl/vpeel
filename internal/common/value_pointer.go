package common

// 辅助函数，返回字符串指针
func StringPointer(s string) *string {
	return &s
}

// 辅助函数，返回整数指针
func IntPointer(i int) *int {
	return &i
}