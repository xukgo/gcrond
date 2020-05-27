package psutil

func checkStringIsNumberic(str string) bool {
	for idx := range str {
		if str[idx] >= '0' && str[idx] <= '9' {
			continue
		}
		return false
	}
	return true
}
