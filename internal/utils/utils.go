package utils

// fastest way we can do a string compare on a list
func BelongsToIgnoreList(needle string) bool {
	switch needle {
	case
		"node_modules",
		"vendor",
		"scripts",
		"docs",
		".git":
		return true
	}
	return false
}
