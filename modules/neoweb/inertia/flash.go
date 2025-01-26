package inertia

import (
	"context"
	"fmt"

	"github.com/romsar/gonertia"
)

type SimpleFlashProvider struct {
	errors       map[string]gonertia.ValidationErrors
	clearHistory map[string]bool
}

func NewSimpleFlashProvider() *SimpleFlashProvider {
	return &SimpleFlashProvider{
		errors:       make(map[string]gonertia.ValidationErrors),
		clearHistory: make(map[string]bool),
	}
}

func (f *SimpleFlashProvider) FlashErrors(ctx context.Context, errors gonertia.ValidationErrors) error {
	sessionID := f.getSessionIDFromContext(ctx)
	f.errors[sessionID] = errors
	return nil
}

func (f *SimpleFlashProvider) GetErrors(ctx context.Context) (gonertia.ValidationErrors, error) {
	sessionID := f.getSessionIDFromContext(ctx)
	errors, ok := f.errors[sessionID]
	if !ok {
		return nil, fmt.Errorf("History doesn't exist with that session id")
	}
	delete(f.errors, sessionID)
	return errors, nil
}

func (f *SimpleFlashProvider) FlashClearHistory(ctx context.Context) error {
	sessionID := f.getSessionIDFromContext(ctx)
	f.clearHistory[sessionID] = true
	return nil
}

func (f *SimpleFlashProvider) ShouldClearHistory(ctx context.Context) (bool, error) {
	sessionID := f.getSessionIDFromContext(ctx)
	clearHistory, ok := f.clearHistory[sessionID]
	if !ok {
		return false, fmt.Errorf("History doesn't exist with that session id")
	}
	delete(f.clearHistory, sessionID)
	return clearHistory, nil
}

func (f *SimpleFlashProvider) getSessionIDFromContext(ctx context.Context) string {
	sessionId, ok := ctx.Value("session_id").(string)
	if !ok {
		return ""
	}

	return sessionId
}
