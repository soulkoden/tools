package cache_test

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/soulkoden/tools/cache"
	"github.com/soulkoden/tools/contract"
	"github.com/stretchr/testify/assert"
)

func TestFilesystemCache(t *testing.T) {
	t.Parallel()

	itemPool, err := cache.NewFilesystemItemPool(os.TempDir())
	assert.NoError(t, err)

	const item1Key = "example~"

	_ = os.Remove(path.Join(os.TempDir(), item1Key))

	assert.False(t, itemPool.HasItem(item1Key))

	item1 := itemPool.GetItem(item1Key)
	item1.Set([]byte("some value"))

	assert.True(t, itemPool.Save(item1))

	item2 := itemPool.GetItem(item1Key)
	assert.True(t, item2.IsHit())

	assert.True(t, itemPool.DeleteItem(item2.GetKey()))

	assert.False(t, item2.IsHit())

	item3 := itemPool.GetItem(item1Key)
	item3.Set([]byte("other value"))
	item3.ExpiresAt(time.Now())
	itemPool.Save(item3)

	time.Sleep(1 * time.Second)

	assert.False(t, item3.IsHit())
}

func TestInMemoryCache(t *testing.T) {
	t.Parallel()

	itemPool, err := cache.NewInMemoryItemPool[map[string]string]()
	assert.NoError(t, err)

	const item3Key = "example3~"
	const item4Key = "example4~"

	assert.False(t, itemPool.HasItem(item3Key))
	assert.False(t, itemPool.HasItem(item4Key))

	item1 := itemPool.GetItem(item3Key)
	item1.Set(map[string]string{"key1": "value1"})
	item1.ExpiresAfter(1 * time.Second)

	assert.True(t, itemPool.SaveDeferred(item1))

	item2 := itemPool.GetItem(item4Key)
	item2.Set(map[string]string{"key2": "value2"})
	item2.ExpiresAfter(2 * time.Second)

	assert.True(t, itemPool.SaveDeferred(item2))

	assert.False(t, itemPool.HasItem(item3Key))
	assert.False(t, itemPool.HasItem(item4Key))

	assert.True(t, itemPool.Commit())

	assert.True(t, itemPool.HasItem(item3Key))
	assert.True(t, itemPool.HasItem(item4Key))

	item3 := itemPool.GetItem(item3Key)

	assert.True(t, item3.IsHit())
	assert.Equal(t, map[string]string{"key1": "value1"}, item3.Get())

	time.Sleep(1 * time.Second)

	items := itemPool.GetItems([]string{item3Key, item4Key})
	assert.Len(t, items, 2)

	item4, item5 := items[0], items[1]

	assert.False(t, item4.IsHit())
	assert.True(t, item5.IsHit())
}

func TestCacheable(t *testing.T) {
	t.Parallel()

	itemPool, err := cache.NewFilesystemItemPool(os.TempDir())
	assert.NoError(t, err)

	const item2Key = "example2~"

	_ = os.Remove(path.Join(os.TempDir(), item2Key))

	value, err := cache.Cacheable[[]byte](itemPool, item2Key, func(item contract.Item[[]byte]) ([]byte, error) {
		item.ExpiresAfter(1 * time.Second)

		return []byte("some value"), nil
	})

	assert.NoError(t, err)
	assert.Equal(t, []byte("some value"), value)

	item2 := itemPool.GetItem(item2Key)
	assert.True(t, item2.IsHit())
}
