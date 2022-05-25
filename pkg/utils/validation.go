package utils

func Validate(d, c []string) bool {
	if len(d) != len(c) {
		return false
	}
	for i, v := range d {
		if v != c[i] {
			return false
		}
	}
	return true
}
