package entity

// Operator represents comparison operations
type Operator string

const (
	OpEquals           Operator = "eq"
	OpNotEquals        Operator = "neq"
	OpGreaterThan      Operator = "gt"
	OpGreaterThanEqual Operator = "gte"
	OpLessThan         Operator = "lt"
	OpLessThanEqual    Operator = "lte"
	OpContains         Operator = "contains"
	OpNotContains      Operator = "not_contains"
	OpStartsWith       Operator = "starts_with"
	OpEndsWith         Operator = "ends_with"
	OpIn               Operator = "in"
	OpNotIn            Operator = "not_in"
	OpExists           Operator = "exists"
	OpNotExists        Operator = "not_exists"
	OpMatches          Operator = "matches" // For regex
)
