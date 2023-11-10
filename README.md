# Implements the PSR Cache Interface for go using file storage


Example usage:

```go
itemPool, _ := cache.NewItemPool(os.TempDir())

item := itemPool.GetItem("my_key")
item.Set([]byte("My Data"))
item.ExpiresAfter(24 * time.Hour)

itemPool.Save()
```
