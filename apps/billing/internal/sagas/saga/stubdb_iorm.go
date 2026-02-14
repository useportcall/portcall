package saga

import (
	"reflect"

	"gorm.io/gorm"
)

func (s *StubDB) FindForID(id uint, dest any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := reflect.ValueOf(dest)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	byID, ok := s.entities[t]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	entity, ok := byID[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	ev := reflect.ValueOf(entity)
	if ev.Kind() == reflect.Ptr {
		ev = ev.Elem()
	}
	v.Set(ev)
	return nil
}

func (s *StubDB) FindFirst(dest any, conds ...any) error {
	if s.FindFirstFn != nil {
		return s.FindFirstFn(dest, conds)
	}
	return gorm.ErrRecordNotFound
}

func (s *StubDB) Create(value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Auto-assign ID if zero
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	idField := v.FieldByName("ID")
	if idField.IsValid() && idField.Uint() == 0 {
		idField.SetUint(uint64(s.nextID))
		s.nextID++
	}
	// Store the entity
	t := v.Type()
	id := uint(idField.Uint())
	if _, ok := s.entities[t]; !ok {
		s.entities[t] = make(map[uint]any)
	}
	s.entities[t][id] = value
	s.Created = append(s.Created, value)
	return nil
}

func (s *StubDB) Save(dest any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Saved = append(s.Saved, dest)
	return nil
}

func (s *StubDB) List(dest any, conds ...any) error {
	if s.ListFn != nil {
		return s.ListFn(dest, conds)
	}
	return nil
}

func (s *StubDB) ListIDs(table string, dest any, conds ...any) error {
	if s.ListIDsFn != nil {
		return s.ListIDsFn(table, dest, conds)
	}
	return nil
}

func (s *StubDB) Count(count *int64, dest any, query string, args ...any) error {
	if s.CountFn != nil {
		return s.CountFn(count, dest, query, args)
	}
	*count = 0
	return nil
}
