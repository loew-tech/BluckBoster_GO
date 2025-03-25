package data

func SliceContains(slice []string, val string) (bool, int) {
	for i, item := range slice {
		if item == val {
			return true, i
		}
	}
	return false, -1
}
