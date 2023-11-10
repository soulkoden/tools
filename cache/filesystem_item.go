package cache

import (
	"os"
	"path"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/soulkoden/tools/contract"
)

type FilesystemItem struct {
	dir        string
	key        string
	value      []byte
	expiration *time.Time
}

func NewFilesystemItem(dir string, key string) contract.Item[[]byte] {
	return &FilesystemItem{
		dir:   dir,
		key:   key,
		value: nil,
	}
}

func (t *FilesystemItem) GetKey() string {
	return t.key
}

func (t *FilesystemItem) Get() []byte {
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

func (t *FilesystemItem) IsHit() bool {
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

	stat_t, ok := stat.Sys().(*syscall.Stat_t)

	if !ok {
		logrus.Errorf("os no supports Stat_t")

		return true
	}

	updatedAtSpec := stat_t.Mtimespec
	updatedAt := time.Unix(updatedAtSpec.Sec, updatedAtSpec.Nsec)

	return t.expiration.Sub(updatedAt) >= 0
}

func (t *FilesystemItem) Set(value []byte) {
	t.value = value
}

func (t *FilesystemItem) ExpiresAt(expiration time.Time) {
	t.expiration = &expiration
}

func (t *FilesystemItem) ExpiresAfter(duration time.Duration) {
	t.ExpiresAt(time.Now().Add(duration))
}
