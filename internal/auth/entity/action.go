package entity

type Action int

const (
	NOOP Action = iota
	READ
	CREATE
	UPDATE
	DELETE
)

func (a Action) String() string {
	switch a {
	default:
		return "unrecognize"
	case NOOP:
		return "noop"
	case READ:
		return "read"
	case CREATE:
		return "create"
	case UPDATE:
		return "update"
	case DELETE:
		return "delete"
	}
}
