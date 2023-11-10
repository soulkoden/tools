package contract

type ItemPool[T any] interface {
	GetItem(key string) Item[T]
	GetItems(keys []string) []Item[T]
	HasItem(key string) bool
	Clear() bool
	DeleteItem(key string) bool
	DeleteItems(keys []string) bool
	Save(item Item[T]) bool
	SaveDeferred(item Item[T]) bool
	Commit() bool
}
