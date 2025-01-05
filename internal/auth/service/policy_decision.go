package service

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/mrrizkin/pohara/internal/auth/access"
	"github.com/mrrizkin/pohara/internal/auth/entity"
	"github.com/mrrizkin/pohara/internal/auth/repository"
)

type PolicyDecision struct {
	authRepo *repository.Repository
}

type PolicyDecisionDependencies struct {
	AuthRepository *repository.Repository
}

func NewPolicyDecision(deps PolicyDecisionDependencies) *PolicyDecision {
	return &PolicyDecision{
		authRepo: deps.AuthRepository,
	}
}

// EvaluateAccess makes and explains the access decision
func (e *PolicyDecision) EvaluateAccess(
	subject entity.Subject,
	action access.Action,
	resource entity.Resource,
	ctx entity.AuthContext,
) (bool, *uint, []uint, string) {
	var policyIDs []uint
	var matchedPolicyID *uint
	var explanation strings.Builder

	policies, err := e.authRepo.GetRolePolicy(subject, resource, action)
	if err != nil {
		explanation.WriteString("Error fetching policies")
		return false, nil, nil, explanation.String()
	}

	if len(policies) == 0 {
		explanation.WriteString("No applicable policies found")
		return false, nil, nil, explanation.String()
	}

	resourceAttrs := resource.GetAttributes()
	subjectAttrs := subject.GetAttributes()
	contextAttrs := ctx.GetAttributes()

	highestPriority := policies[0].Priority
	var allowPolicies, denyPolicies []entity.Policy

	for _, policy := range policies {
		policyIDs = append(policyIDs, policy.ID)

		if policy.Priority != highestPriority {
			continue
		}

		if e.evaluateCondition(policy.Conditions, resourceAttrs, subjectAttrs, contextAttrs) {
			if policy.Effect == access.EffectAllow {
				allowPolicies = append(allowPolicies, policy)
			} else {
				denyPolicies = append(denyPolicies, policy)
			}
		}
	}

	if len(denyPolicies) > 0 {
		matchedPolicyID = &denyPolicies[0].ID
		return false, matchedPolicyID, policyIDs, "Denied by policy"
	}

	if len(allowPolicies) > 0 {
		matchedPolicyID = &allowPolicies[0].ID
		return true, matchedPolicyID, policyIDs, "Allowed by policy"
	}

	return false, nil, policyIDs, "No matching policies"
}

// evaluateCondition evaluates a complex condition
func (e *PolicyDecision) evaluateCondition(
	condition entity.Condition,
	resourceAttrs, subjectAttrs, contextAttrs map[string]interface{},
) bool {
	if condition.Not != nil {
		return !e.evaluateCondition(*condition.Not, resourceAttrs, subjectAttrs, contextAttrs)
	}

	for _, rule := range condition.AllOf {
		if !e.evaluateRule(rule, resourceAttrs, subjectAttrs, contextAttrs) {
			return false
		}
	}

	if len(condition.AnyOf) > 0 {
		anyMatched := false
		for _, rule := range condition.AnyOf {
			if e.evaluateRule(rule, resourceAttrs, subjectAttrs, contextAttrs) {
				anyMatched = true
				break
			}
		}
		if !anyMatched {
			return false
		}
	}

	return true
}

// evaluateRule evaluates a single condition rule
func (e *PolicyDecision) evaluateRule(
	rule entity.ConditionRule,
	resourceAttrs, subjectAttrs, contextAttrs map[string]interface{},
) bool {
	if rule.Resource != nil {
		if !evaluateAttributeRule(*rule.Resource, resourceAttrs) {
			return false
		}
	}
	if rule.Subject != nil {
		if !evaluateAttributeRule(*rule.Subject, subjectAttrs) {
			return false
		}
	}
	if rule.Context != nil {
		if !evaluateAttributeRule(*rule.Context, contextAttrs) {
			return false
		}
	}
	return true
}

// evaluateAttributeRule evaluates a single attribute rule
func evaluateAttributeRule(rule entity.AttributeRule, attrs map[string]interface{}) bool {
	value, exists := attrs[rule.Attribute]
	if !exists {
		return rule.Operator == entity.OpNotExists
	}
	if rule.Operator == entity.OpExists {
		return true
	}

	switch rule.Operator {
	case entity.OpEquals:
		return reflect.DeepEqual(value, rule.Value)
	case entity.OpNotEquals:
		return !reflect.DeepEqual(value, rule.Value)
	case entity.OpGreaterThan:
		return compareValues(value, rule.Value) > 0
	case entity.OpGreaterThanEqual:
		return compareValues(value, rule.Value) >= 0
	case entity.OpLessThan:
		return compareValues(value, rule.Value) < 0
	case entity.OpLessThanEqual:
		return compareValues(value, rule.Value) <= 0
	case entity.OpContains:
		return containsValue(value, rule.Value)
	case entity.OpNotContains:
		return !containsValue(value, rule.Value)
	case entity.OpStartsWith:
		return strings.HasPrefix(fmt.Sprint(value), fmt.Sprint(rule.Value))
	case entity.OpEndsWith:
		return strings.HasSuffix(fmt.Sprint(value), fmt.Sprint(rule.Value))
	case entity.OpIn:
		return valueInSlice(value, rule.Value)
	case entity.OpNotIn:
		return !valueInSlice(value, rule.Value)
	case entity.OpMatches:
		pattern, ok := rule.Value.(string)
		if !ok {
			return false
		}
		matched, err := regexp.MatchString(pattern, fmt.Sprint(value))
		return err == nil && matched
	}
	return false
}

// Helper functions for value comparison and manipulation
func compareValues(a, b interface{}) int {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	switch aVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		aInt := aVal.Int()
		bInt := bVal.Int()
		if aInt < bInt {
			return -1
		} else if aInt > bInt {
			return 1
		}
		return 0
	case reflect.Float32, reflect.Float64:
		aFloat := aVal.Float()
		bFloat := bVal.Float()
		if aFloat < bFloat {
			return -1
		} else if aFloat > bFloat {
			return 1
		}
		return 0
	case reflect.String:
		aString := aVal.String()
		bString := bVal.String()
		return strings.Compare(aString, bString)
	default:
		return 0
	}
}

func containsValue(container, item interface{}) bool {
	containerVal := reflect.ValueOf(container)

	switch containerVal.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < containerVal.Len(); i++ {
			if reflect.DeepEqual(containerVal.Index(i).Interface(), item) {
				return true
			}
		}
	case reflect.String:
		if itemStr, ok := item.(string); ok {
			return strings.Contains(containerVal.String(), itemStr)
		}
	}
	return false
}

func valueInSlice(value, slice interface{}) bool {
	sliceVal := reflect.ValueOf(slice)
	if sliceVal.Kind() != reflect.Slice && sliceVal.Kind() != reflect.Array {
		return false
	}

	for i := 0; i < sliceVal.Len(); i++ {
		if reflect.DeepEqual(sliceVal.Index(i).Interface(), value) {
			return true
		}
	}
	return false
}
