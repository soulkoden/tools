# Implements the PSR Cache Interface for go using file storage

Example usage:

```go
package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/soulkoden/tools/cache"
)

func main() {
	// create new item
	itemPool, _ := cache.NewFilesystemItemPool(os.TempDir())
	// itemPool, _ := cache.NewInMemoryItemPool[MyType]()

	item := itemPool.GetItem("my_key")
	item.Set([]byte("My Data"))
	item.ExpiresAfter(24 * time.Hour)

	itemPool.Save(item)

	// working with exists item later
	item = itemPool.GetItem("my_key")
	if item.IsHit() {
		spew.Dump(item.Get())
	}

	// symfony (closure) style
	value, _ := cache.Cacheable[[]byte](itemPool, "my_key", func(item contract.Item[[]byte]) ([]byte, error) {
		// this closure used as factory for creating/updating cache only
		item.ExpiresAfter(24 * time.Hour)

		return []byte("Inline Data"), nil
	})

	spew.Dump(value)
}
```
