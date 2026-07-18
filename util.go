package revcatgo

import "slices"

func contains(s []string, e string) bool {
	return slices.Contains(s, e)
}
