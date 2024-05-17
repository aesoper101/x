package factory

import (
	"fmt"
	"strings"
	"sync"
)

type safeFactoryFn[T any] map[string]func() (T, error)

// SafeFactory is a factory that can create objects of type T. It is safe for concurrent use.
type SafeFactory[T any] struct {
	rw sync.RWMutex
	m  safeFactoryFn[T]
}

func NewSafeFactory[T any]() *SafeFactory[T] {
	return &SafeFactory[T]{
		m: make(safeFactoryFn[T]),
	}
}

// Register registers a new factory functionx for the given type.
// The name is used to identify the type in the factory. If the name is already registered, returns an error. it is not case sensitive.
func (f *SafeFactory[T]) Register(name string, fn func() (T, error)) error {
	f.rw.Lock()
	defer f.rw.Unlock()

	name = strings.ToLower(name)

	if _, ok := f.m[name]; ok {
		return fmt.Errorf("factory: %s already registered", name)
	}

	f.m[name] = fn

	return nil
}

func (f *SafeFactory[T]) Unregister(name string) error {
	f.rw.Lock()
	defer f.rw.Unlock()

	name = strings.ToLower(name)

	if _, ok := f.m[name]; !ok {
		return fmt.Errorf("factory: %s not registered", name)
	}

	delete(f.m, name)
	return nil
}

func (f *SafeFactory[T]) Get(name string) (func() (T, error), error) {
	f.rw.RLock()
	defer f.rw.RUnlock()

	name = strings.ToLower(name)

	fn, ok := f.m[name]
	if !ok {
		return nil, fmt.Errorf("factory: %s not registered", name)
	}

	return fn, nil
}
