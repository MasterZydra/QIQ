package common

func ImplodeSlice(slice []string, separator string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

func ImplodeStrSlice(slice []string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += ", "
		}
		result += "\"" + s + "\""
	}
	return result
}

func RemoveIndex[T any](s []T, index int) []T {
	return append(append(make([]T, 0), s[:index]...), s[index+1:]...)
}
