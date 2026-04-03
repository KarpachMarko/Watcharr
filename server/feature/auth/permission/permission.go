package permission

import "github.com/sbondCo/Watcharr/database/entity"

// If `perms` has `req(uired)Perm`.
func Has(perms int, reqPerm int) bool {
	// Admins have permission for everything.
	if perms&entity.PERM_ADMIN == entity.PERM_ADMIN {
		return true
	}
	return (perms & reqPerm) == reqPerm
}
