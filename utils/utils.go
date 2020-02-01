package utils

// GetLastItemString of slice
func GetLastItemString(slice []string) string {
	if len(slice) <= 0 {
		return ""
	}
	return slice[len(slice)-1]
}

// CalcSizeInByte func
func CalcSizeInByte(number float64, unit string) int64 {
	switch unit {
	case "KB":
		return int64(number) * 1000
	case "MB":
		return int64(number) * 1000000
	case "GB":
		return int64(number) * 10000000000
	default:
		return int64(number)
	}
}
