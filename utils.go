package paypalutil

func isUpperNumber[T ~string | ~[]uint8](v T) bool {
	length := len(v)
	if length == 0 {
		return false
	}
	for i := 0; i < length; i++ {
		if 0x41 > v[i] || v[i] > 0x5A {
			if v[i] < '0' || v[i] > '9' {
				return false
			}
		}
	}
	return true
}
