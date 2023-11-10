package cache_test

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/soulkoden/tools/cache"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
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
