package cache

import (
	"errors"
	"fmt"

	"github.com/soulkoden/tools/contract"
)

func Cacheable[T any](itemPool contract.ItemPool[T], key string, factory func(contract.Item[T]) (T, error)) (T, error) {
	var def T

	item := itemPool.GetItem(key)
	if item.IsHit() {
		return item.Get(), nil
	}

	value, err := factory(item)
	if err != nil {
		return def, fmt.Errorf("failed to call factory: %w", err)
	}

	item.Set(value)

	if !itemPool.Save(item) {
		return def, errors.New("cannot save item into item pool")
	}

	return value, nil
}
