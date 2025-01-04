package entity

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
