package common

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
