package checkutil

func IsLegalID(i int, floor int, ceil int) bool {
	if i <= floor || i > ceil {
		return false
	}
	return true
}
