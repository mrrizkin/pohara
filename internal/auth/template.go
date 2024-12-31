package auth

import (
	"github.com/mrrizkin/pohara/internal/auth/entity"
	"github.com/mrrizkin/pohara/internal/web/template"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func Can() *exec.ControlStructureSet {
	return template.CustomIf("can", func(role entity.Role, permission string) bool {
		return role.Can(permission)
	})
}
