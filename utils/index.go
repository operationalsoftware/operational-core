package utils

// CheckRole checks if a role is in a list of roles.
//
// Parameters:
// @roles - []string
// @role - string
//
// Returns:
// bool - Whether the role is in the list of roles.
func CheckRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
