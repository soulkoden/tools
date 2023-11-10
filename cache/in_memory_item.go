package cache

import (
	"time"

	"github.com/soulkoden/tools/contract"
)

type InMemoryItem[T any] struct {
	key        string
	value      T
	expiration *time.Time
}

func NewInMemoryItem[T any](key string) contract.Item[T] {
	return &InMemoryItem[T]{
		key: key,
	}
}

func (t *InMemoryItem[T]) GetKey() string {
	return t.key
}

func (t *InMemoryItem[T]) Get() T {
	return t.value
}

func (t *InMemoryItem[T]) IsHit() bool {
	if t.expiration == nil {
		return true
	}

	return t.expiration.Sub(time.Now()) > 0
}

func (t *InMemoryItem[T]) Set(value T) {
	t.value = value
}

func (t *InMemoryItem[T]) ExpiresAt(expiration time.Time) {
	t.expiration = &expiration
}

func (t *InMemoryItem[T]) ExpiresAfter(duration time.Duration) {
	t.ExpiresAt(time.Now().Add(duration))
}
