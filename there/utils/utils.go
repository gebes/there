package there


//CheckArrayContains checks if a string array contains a specific element
func CheckArrayContains(slice []string, toSearch string) bool {
	for _, s := range slice {
		if toSearch == s {
			return true
		}
	}
	return false
}


