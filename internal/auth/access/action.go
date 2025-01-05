package access

import "fmt"

// Action represents possible actions on resources
type Action string

var (
	ActionGeneralCreate = NewResourceAction("general", "create")
	ActionGeneralRead   = NewResourceAction("general", "read")
	ActionGeneralUpdate = NewResourceAction("general", "update")
	ActionGeneralDelete = NewResourceAction("general", "delete")
)

func NewResourceAction(resourceType string, action string) Action {
	return Action(fmt.Sprintf("%s:%s", resourceType, action))
}
