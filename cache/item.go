package cache

import (
	"os"
	"path"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type Item struct {
	dir        string
	key        string
	value      []byte
	expiration *time.Time
}

func NewItem(dir string, key string) *Item {
	return &Item{
		dir:   dir,
		key:   key,
		value: nil,
	}
}

func (t *Item) GetKey() string {
	return t.key
}

func (t *Item) Get() []byte {
	if t.value == nil {
		var (
			err      error
			filename = path.Join(t.dir, t.GetKey())
		)

		if t.value, err = os.ReadFile(filename); err != nil {
			logrus.WithError(err).WithField("filename", filename).Error("failed to read cache file")
		}
	}

	return t.value
}

func (t *Item) IsHit() bool {
	filename := path.Join(t.dir, t.GetKey())
	stat, err := os.Stat(filename)

	if err != nil {
		if !os.IsNotExist(err) {
			logrus.WithError(err).WithField("filename", filename).Errorf("failed to get file stat")
		}

		return false
	}

	if t.expiration == nil {
		return true
	}

	updatedAtSpec := stat.Sys().(*syscall.Stat_t).Mtimespec
	updatedAt := time.Unix(updatedAtSpec.Sec, updatedAtSpec.Nsec)

	return t.expiration.Sub(updatedAt) >= 0
}

func (t *Item) Set(value []byte) {
	t.value = value
}

func (t *Item) ExpiresAt(expiration time.Time) {
	t.expiration = &expiration
}

func (t *Item) ExpiresAfter(duration time.Duration) {
	t.ExpiresAt(time.Now().Add(duration))
}
