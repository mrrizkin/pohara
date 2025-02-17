package migration

type Migration interface {
	Up(schema *Schema)
	Down(schema *Schema)
	ID() string
}
