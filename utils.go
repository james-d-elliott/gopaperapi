package gopaperapi

// IsStringInSlice checks if a single string is in an array of strings.
func IsStringInSlice(a string, list []string) (inSlice bool) {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
