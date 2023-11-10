package cache

import (
	"sync"

	"github.com/soulkoden/tools/contract"
)

type InMemoryItemPool[T any] struct {
	storage sync.Map
	queue   []contract.Item[T]
}

func NewInMemoryItemPool[T any]() (contract.ItemPool[T], error) {
	return &InMemoryItemPool[T]{
		queue: make([]contract.Item[T], 0),
	}, nil
}

func (t *InMemoryItemPool[T]) GetItem(key string) contract.Item[T] {
	exists, ok := t.storage.Load(key)
	if !ok {
		return NewInMemoryItem[T](key)
	}

	return exists.(contract.Item[T])
}

func (t *InMemoryItemPool[T]) GetItems(keys []string) []contract.Item[T] {
	var items = make([]contract.Item[T], len(keys))
	for i, key := range keys {
		items[i] = t.GetItem(key)
	}

	return items
}

func (t *InMemoryItemPool[T]) HasItem(key string) bool {
	_, ok := t.storage.Load(key)

	return ok
}

func (t *InMemoryItemPool[T]) Clear() bool {
	t.storage.Range(func(key, _ any) bool {
		return t.DeleteItem(key.(string))
	})

	return true
}

func (t *InMemoryItemPool[T]) DeleteItem(key string) bool {
	t.storage.Delete(key)

	return true
}

func (t *InMemoryItemPool[T]) DeleteItems(keys []string) bool {
	success := true

	for _, key := range keys {
		if !t.DeleteItem(key) {
			success = false
		}
	}

	return success
}

func (t *InMemoryItemPool[T]) Save(item contract.Item[T]) bool {
	t.storage.Store(item.GetKey(), item)

	return true
}

func (t *InMemoryItemPool[T]) SaveDeferred(item contract.Item[T]) bool {
	t.queue = append(t.queue, item)

	return true
}

func (t *InMemoryItemPool[T]) Commit() bool {
	success := true

	for _, item := range t.queue {
		if !t.Save(item) {
			success = false
		}
	}

	if success {
		t.queue = make([]contract.Item[T], 0)
	}

	return success
}
