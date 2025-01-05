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

type Condition struct {
	AllOf []ConditionRule `json:"all_of,omitempty"` // AND conditions
	AnyOf []ConditionRule `json:"any_of,omitempty"` // OR conditions
	Not   *Condition      `json:"not,omitempty"`    // NOT conditions
}

// ConditionRule represents a single rule in a condition
type ConditionRule struct {
	Resource *AttributeRule `json:"resource,omitempty"`
	Subject  *AttributeRule `json:"subject,omitempty"`
	Context  *AttributeRule `json:"context,omitempty"`
}

// AttributeRule defines how to evaluate an attribute
type AttributeRule struct {
	Attribute string      `json:"attribute"`
	Operator  Operator    `json:"operator"`
	Value     interface{} `json:"value"`
}
