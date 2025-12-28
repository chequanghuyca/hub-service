package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

var (
	ErrInvalidState = errors.New("invalid or expired state")
)

type stateEntry struct {
	createdAt   time.Time
	frontendURL string
}

// StateManager manages OAuth state tokens for CSRF protection
type StateManager struct {
	states map[string]stateEntry
	mu     sync.RWMutex
	ttl    time.Duration
}

// NewStateManager creates a new state manager with the given TTL
func NewStateManager(ttl time.Duration) *StateManager {
	sm := &StateManager{
		states: make(map[string]stateEntry),
		ttl:    ttl,
	}

	// Start cleanup goroutine
	go sm.cleanup()

	return sm
}

// Generate creates a new random state token with frontend URL
func (sm *StateManager) Generate(frontendURL string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)

	sm.mu.Lock()
	sm.states[state] = stateEntry{
		createdAt:   time.Now(),
		frontendURL: frontendURL,
	}
	sm.mu.Unlock()

	return state, nil
}

// Validate checks if a state token is valid and returns the frontend URL
func (sm *StateManager) Validate(state string) (bool, string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	entry, exists := sm.states[state]
	if !exists {
		return false, ""
	}

	// Check if expired
	if time.Since(entry.createdAt) > sm.ttl {
		delete(sm.states, state)
		return false, ""
	}

	// Get frontend URL before removing state
	frontendURL := entry.frontendURL

	// Remove state after validation (one-time use)
	delete(sm.states, state)
	return true, frontendURL
}

// cleanup periodically removes expired states
func (sm *StateManager) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sm.mu.Lock()
		now := time.Now()
		for state, entry := range sm.states {
			if now.Sub(entry.createdAt) > sm.ttl {
				delete(sm.states, state)
			}
		}
		sm.mu.Unlock()
	}
}
