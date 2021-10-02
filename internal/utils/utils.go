package utils

// BelongsToIgnoreList is fastest way we can do a string compare on a list
func BelongsToIgnoreList(needle string) bool {
	switch needle {
	case
		"node_modules",
		"vendor",
		"scripts",
		"docs",
		"test",
		"tests",
		".git":
		return true
	}
	return false
}

func CharIsDigit(c string) bool {
	if len(c) == 0 {
		return false
	}

	if c[0] < '0' || c[0] > '9' {
		return false
	}
	return true
}
