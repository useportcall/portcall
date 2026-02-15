package harness

import (
	"context"
	"fmt"
	"sync"
)

// MemoryStore is an in-memory replacement for S3/object storage.
type MemoryStore struct {
	mu         sync.RWMutex
	signatures map[string][]byte
	icons      map[string][]byte
}

// NewMemoryStore returns a fresh in-memory object store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		signatures: map[string][]byte{},
		icons:      map[string][]byte{},
	}
}

func (m *MemoryStore) GetFromSignatureBucket(id string, _ context.Context) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, ok := m.signatures[id]
	if !ok {
		return nil, fmt.Errorf("signature not found")
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	return cp, nil
}

func (m *MemoryStore) PutInSignatureBucket(id string, data []byte, _ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make([]byte, len(data))
	copy(cp, data)
	m.signatures[id] = cp
	return nil
}

func (m *MemoryStore) GetFromIconLogoBucket(id string, _ context.Context) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, ok := m.icons[id]
	if !ok {
		return nil, fmt.Errorf("icon not found")
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	return cp, nil
}

func (m *MemoryStore) PutInIconLogoBucket(id string, data []byte, _ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make([]byte, len(data))
	copy(cp, data)
	m.icons[id] = cp
	return nil
}

func (m *MemoryStore) DeleteFromIconLogoBucket(id string, _ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.icons, id)
	return nil
}
