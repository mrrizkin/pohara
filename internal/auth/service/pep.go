package service

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/mrrizkin/pohara/internal/auth/entity"
	"github.com/mrrizkin/pohara/internal/ports"
	"gorm.io/gorm"
)

// AuthContext holds contextual information for authorization
type AuthContext struct {
	Time     time.Time
	IP       string
	Location string
	// Add other context attributes as needed
}

func (c AuthContext) Hash() string {
	return fmt.Sprintf("%d:%s:%s", c.Time.Unix(), c.IP, c.Location)
}

type Enforcer struct {
	db    ports.Database
	cache ports.Cache
}

type EnforcerDependencies struct {
	Database ports.Database
	Cache    ports.Cache
}

func NewEnforcer(deps EnforcerDependencies) *Enforcer {
	return &Enforcer{
		db:    deps.Database,
		cache: deps.Cache,
	}
}

func (e *Enforcer) IsAllowed(subject entity.Subject, action entity.Action, resource entity.Resource, ctx AuthContext) bool {
	start := time.Now()
	traceID := "traceID"

	auditLog := &entity.AuditAuthLog{
		TraceID:      traceID,
		Timestamp:    start,
		SubjectID:    subject.GetSubjectID(),
		SubjectName:  subject.GetSubjectName(),
		Action:       action,
		ResourceType: resource.GetResourceType(),
		ResourceID:   resource.GetResourceID(),
	}

	defer func() {
		auditLog.Duration = time.Since(start).Nanoseconds()
		e.db.Create(auditLog)
	}()

	cacheKey := e.generateCacheKey(subject, action, resource, ctx)
	if cachedDecision, ok := e.checkCache(cacheKey); ok {
		auditLog.Decision = "Cache hit"
		auditLog.Effect = cachedDecision
		return cachedDecision
	}

	decision, matchedPolicyID, policyIDs, explanation := e.evaluateAccess(subject, action, resource, ctx)

	auditLog.Effect = decision
	auditLog.PolicyIDs = policyIDs
	auditLog.MatchedPolicyID = matchedPolicyID
	auditLog.Decision = explanation

	e.cacheDecision(cacheKey, decision)

	return false
}

// evaluateAccess makes and explains the access decision
func (e *Enforcer) evaluateAccess(
	subject entity.Subject,
	action entity.Action,
	resource entity.Resource,
	ctx AuthContext,
) (bool, *uint, []uint, string) {
	var policies []entity.Policy
	var policyIDs []uint
	var matchedPolicyID *uint
	var explanation strings.Builder

	db := e.db.GetDB().(*gorm.DB)

	result := db.Preload("Roles").
		Joins("JOIN role_policies ON role_policies.policy_id = policies.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_policies.role_id").
		Where("user_roles.subject_id = ? AND policies.resource_type = ? AND policies.action = ?",
			subject.GetSubjectID(), resource.GetResourceType(), action).
		Order("policies.priority DESC").
		Find(&policies)

	if result.Error != nil {
		explanation.WriteString("Error fetching policies")
		return false, nil, nil, explanation.String()
	}

	if len(policies) == 0 {
		explanation.WriteString("No applicable policies found")
		return false, nil, nil, explanation.String()
	}

	resourceAttrs := resource.GetAttributes()
	subjectAttrs := subject.GetAttributes()
	contextAttrs := getContextAttributes(ctx)

	highestPriority := policies[0].Priority
	var allowPolicies, denyPolicies []entity.Policy

	for _, policy := range policies {
		policyIDs = append(policyIDs, policy.ID)

		if policy.Priority != highestPriority {
			continue
		}

		if e.evaluateCondition(policy.Conditions, resourceAttrs, subjectAttrs, contextAttrs) {
			if policy.Effect == entity.EffectAllow {
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
func (e *Enforcer) evaluateCondition(
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
func (e *Enforcer) evaluateRule(
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

// Helper functions for value comparison
func compareValues(a, b interface{}) int {
	switch a := a.(type) {
	case int:
		b := b.(int)
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	case float64:
		b := b.(float64)
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	case string:
		b := b.(string)
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func containsValue(container, item interface{}) bool {
	switch container := container.(type) {
	case string:
		return strings.Contains(container, fmt.Sprint(item))
	case []string:
		for _, v := range container {
			if v == item {
				return true
			}
		}
	case []int:
		for _, v := range container {
			if v == item {
				return true
			}
		}
	case []interface{}:
		for _, v := range container {
			if reflect.DeepEqual(v, item) {
				return true
			}
		}
	}
	return false
}

func valueInSlice(value, slice interface{}) bool {
	switch slice := slice.(type) {
	case []string:
		for _, v := range slice {
			if v == value {
				return true
			}
		}
	case []int:
		for _, v := range slice {
			if v == value {
				return true
			}
		}
	case []interface{}:
		for _, v := range slice {
			if reflect.DeepEqual(v, value) {
				return true
			}
		}
	}
	return false
}

func getContextAttributes(ctx AuthContext) map[string]interface{} {
	return map[string]interface{}{
		"time":     ctx.Time,
		"ip":       ctx.IP,
		"location": ctx.Location,
	}
}

// Cache related methods
func (e *Enforcer) generateCacheKey(subject entity.Subject, action entity.Action, resource entity.Resource, ctx AuthContext) string {
	return fmt.Sprintf("auth:%d:%s:%s:%d:%s",
		subject.GetSubjectID(),
		action,
		resource.GetResourceType(),
		resource.GetResourceID(),
		ctx.Hash(),
	)
}

func (e *Enforcer) checkCache(key string) (bool, bool) {
	ctx := context.Background()
	val, ok := e.cache.Get(ctx, key)
	if !ok {
		return false, false
	}
	return val == "true", true
}

func (e *Enforcer) cacheDecision(key string, decision bool) {
	ctx := context.Background()
	e.cache.Set(ctx, key, fmt.Sprintf("%t", decision))
}
