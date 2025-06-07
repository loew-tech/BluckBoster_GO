package gql

func SetToList[T comparable](set map[T]bool) []T {
	list := make([]T, 0, len(set))
	for item := range set {
		list = append(list, item)
	}
	return list
}
