package utils

// GetLastItemString of slice
func GetLastItemString(slice []string) string {
	if len(slice) <= 0 {
		return ""
	}
	return slice[len(slice)-1]
}
