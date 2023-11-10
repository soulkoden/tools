package contract

import "time"

type Item[T any] interface {
	GetKey() string
	Get() T
	IsHit() bool
	Set(value T)
	ExpiresAt(expiration time.Time)
	ExpiresAfter(duration time.Duration)
}
