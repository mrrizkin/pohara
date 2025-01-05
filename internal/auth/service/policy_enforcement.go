package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mrrizkin/pohara/internal/auth/access"
	"github.com/mrrizkin/pohara/internal/auth/entity"
	"github.com/mrrizkin/pohara/internal/auth/repository"
	"github.com/mrrizkin/pohara/internal/common/nanoid"
	"github.com/mrrizkin/pohara/internal/ports"
)

type PolicyEnforcement struct {
	cache ports.Cache

	authRepo *repository.Repository
	pd       *PolicyDecision
}

type PolicyEnforcementDependencies struct {
	Cache ports.Cache

	AuthRepository *repository.Repository
	PolicyDecision *PolicyDecision
}

func NewPolicyEnforcement(deps PolicyEnforcementDependencies) *PolicyEnforcement {
	return &PolicyEnforcement{
		cache: deps.Cache,

		authRepo: deps.AuthRepository,
		pd:       deps.PolicyDecision,
	}
}

func (pe *PolicyEnforcement) IsAllowed(
	subject entity.Subject,
	action access.Action,
	resource entity.Resource,
	ctx entity.AuthContext,
) bool {
	start := time.Now()
	traceID := nanoid.New()

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
		pe.authRepo.CreateAuditAuthLog(auditLog)
	}()

	cacheKey := pe.generateCacheKey(subject, action, resource, ctx)
	if cachedDecision, ok := pe.checkCache(cacheKey); ok {
		auditLog.Decision = "Cache hit"
		auditLog.Effect = cachedDecision
		return cachedDecision
	}

	decision, matchedPolicyID, policyIDs, explanation := pe.pd.EvaluateAccess(subject, action, resource, ctx)

	auditLog.Effect = decision
	auditLog.PolicyIDs = policyIDs
	auditLog.MatchedPolicyID = matchedPolicyID
	auditLog.Decision = explanation

	pe.cacheDecision(cacheKey, decision)

	return false
}

// Cache related methods
func (e *PolicyEnforcement) generateCacheKey(
	subject entity.Subject,
	action access.Action,
	resource entity.Resource,
	ctx entity.AuthContext,
) string {
	return fmt.Sprintf("auth:%d:%s:%s:%d:%s",
		subject.GetSubjectID(),
		action,
		resource.GetResourceType(),
		resource.GetResourceID(),
		ctx.Hash(),
	)
}

func (e *PolicyEnforcement) checkCache(key string) (bool, bool) {
	ctx := context.Background()
	val, ok := e.cache.Get(ctx, key)
	if !ok {
		return false, false
	}
	return val == "true", true
}

func (e *PolicyEnforcement) cacheDecision(key string, decision bool) {
	ctx := context.Background()
	e.cache.Set(ctx, key, fmt.Sprintf("%t", decision))
}
