package authx

import (
	"sync"

	"github.com/clerk/clerk-sdk-go/v2"
)

// Sample interface for JSON Web Key storage.
// Implementation may vary.
type JWKStore interface {
	GetJWK() *clerk.JSONWebKey
	SetJWK(*clerk.JSONWebKey)
}

// MemoryJWKStore is a thread-safe in-memory implementation of JWKStore
type MemoryJWKStore struct {
	mu  sync.RWMutex
	jwk *clerk.JSONWebKey
}

// GetJWK returns the stored JWK
func (store *MemoryJWKStore) GetJWK() *clerk.JSONWebKey {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.jwk
}

// SetJWK sets the JWK in the store
func (store *MemoryJWKStore) SetJWK(jwk *clerk.JSONWebKey) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.jwk = jwk
}

func NewJWKStore() JWKStore {
	// Implementation may vary. This can be an
	// in-memory store, database, caching layer,...
	return &MemoryJWKStore{}
}
