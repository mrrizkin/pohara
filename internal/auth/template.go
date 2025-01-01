package auth

import "github.com/mrrizkin/pohara/internal/auth/entity"

func Can(role entity.Role, permission string) bool {
	return role.Can(permission)
}
