package saga

import (
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx"
	"gorm.io/gorm"
)

// --- Stubs and simple IORM methods ---

func (s *StubDB) Delete(value, query any, args ...any) error {
	if s.DeleteFn != nil {
		return s.DeleteFn(value, query, args)
	}
	return nil
}

func (s *StubDB) Txn(fn func(orm dbx.IORM) error) error { return fn(s) }

func (s *StubDB) ListWithOrder(any, string, ...any) error { return nil }
func (s *StubDB) ListWithOrderAndLimit(dest any, order string, limit int, conds ...any) error {
	if s.ListWithOrderAndLimitFn != nil {
		return s.ListWithOrderAndLimitFn(dest, order, limit, conds)
	}
	return nil
}
func (s *StubDB) ListWithOrderLimitOffset(any, string, int, int, ...any) error {
	return nil
}
func (s *StubDB) ListForAppID(uint, any, *int) error          { return nil }
func (s *StubDB) ListForPlanID(uint, uint, any, string) error { return nil }
func (s *StubDB) GetForPublicID(uint, string, any) error      { return gorm.ErrRecordNotFound }
func (s *StubDB) FindFirstForAppID(uint, any) error           { return gorm.ErrRecordNotFound }
func (s *StubDB) FindFirstOrNil(any, ...any) error            { return nil }
func (s *StubDB) Update(any, ...any) error                    { return nil }
func (s *StubDB) UpdateForPublicID(uint, string, any) error   { return nil }
func (s *StubDB) Upsert(any, any, ...any) error               { return nil }
func (s *StubDB) UpsertForPublicID(uint, string, any) error   { return nil }
func (s *StubDB) RemoveForPublicID(uint, string, any) error   { return nil }
func (s *StubDB) DeleteForID(any) error                       { return nil }
func (s *StubDB) IncrementCount(any, string, int64) error     { return nil }
func (s *StubDB) Exec(string, ...any) error                   { return nil }
func (s *StubDB) AutoMigrate(...any) error                    { return nil }

// FindCreated returns the first created entity of the given type, or an error.
func FindCreated[T any](db *StubDB) (*T, error) {
	for _, c := range db.Created {
		if v, ok := c.(*T); ok {
			return v, nil
		}
	}
	return nil, fmt.Errorf("no created entity of type %T", (*T)(nil))
}
