package saga

import (
	"reflect"
	"sync"
)

// StubDB is a minimal in-memory dbx.IORM used by integration tests.
// Entities are stored by reflect.Type and uint ID. Query predicates are
// resolved via caller-supplied FindFirst/List callbacks.
type StubDB struct {
	mu       sync.Mutex
	entities map[reflect.Type]map[uint]any // type → id → entity
	nextID   uint

	// FindFirstFn is called for FindFirst queries. Return (entity, nil) to
	// populate dest, or (nil, gorm.ErrRecordNotFound) to signal not found.
	FindFirstFn func(dest any, conds []any) error

	// ListFn handles List calls.
	ListFn func(dest any, conds []any) error

	// ListIDsFn handles ListIDs calls.
	ListIDsFn func(table string, dest any, conds []any) error

	// CountFn handles Count calls.
	CountFn func(count *int64, dest any, query string, args []any) error

	// DeleteFn handles Delete calls.
	DeleteFn func(value any, query any, args []any) error

	// ListWithOrderAndLimitFn handles ListWithOrderAndLimit calls.
	ListWithOrderAndLimitFn func(dest any, order string, limit int, conds []any) error

	// Created collects every entity passed to Create, in order.
	Created []any

	// Saved collects every entity passed to Save, in order.
	Saved []any
}

// NewStubDB creates an empty StubDB.
func NewStubDB() *StubDB {
	return &StubDB{
		entities: make(map[reflect.Type]map[uint]any),
		nextID:   1,
	}
}

// Store adds an entity to the in-memory store keyed by its ID.
// ID is read from the struct's "ID" field via reflection.
func (s *StubDB) Store(entity any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		// Try embedded Model.ID
		modelField := v.FieldByName("Model")
		if modelField.IsValid() {
			idField = modelField.FieldByName("ID")
		}
	}
	id := uint(idField.Uint())
	if _, ok := s.entities[t]; !ok {
		s.entities[t] = make(map[uint]any)
	}
	s.entities[t][id] = entity
}


