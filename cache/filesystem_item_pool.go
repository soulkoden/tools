package cache

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/soulkoden/tools/contract"
)

type FilesystemItemPool struct {
	dir string

	queue []contract.Item[[]byte]
}

func NewFilesystemItemPool(dir string) (contract.ItemPool[[]byte], error) {
	stat, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("cache directory is unreadable: %w", err)
	}

	if !stat.IsDir() {
		return nil, errors.New("cache directory is not a directory")
	}

	return &FilesystemItemPool{
		dir:   dir,
		queue: make([]contract.Item[[]byte], 0),
	}, nil
}

func (t *FilesystemItemPool) GetItem(key string) contract.Item[[]byte] {
	return NewFilesystemItem(t.dir, key)
}

func (t *FilesystemItemPool) GetItems(keys []string) []contract.Item[[]byte] {
	var items = make([]contract.Item[[]byte], len(keys))
	for i, key := range keys {
		items[i] = t.GetItem(key)
	}

	return items
}

func (t *FilesystemItemPool) HasItem(key string) bool {
	filename := path.Join(t.dir, key)

	if _, err := os.Stat(filename); err != nil {
		if !os.IsNotExist(err) {
			logrus.WithError(err).WithField("filename", filename).Errorf("failed to get file stat")
		}

		return false
	}

	return true
}

func (t *FilesystemItemPool) Clear() bool {
	files, err := os.ReadDir(t.dir)

	if err != nil {
		logrus.WithError(err).WithField("dir", t.dir).Errorf("failed to read dir")

		return false
	}

	success := true

	for _, file := range files {
		if !t.DeleteItem(file.Name()) {
			success = false
		}
	}

	return success
}

func (t *FilesystemItemPool) DeleteItem(key string) bool {
	filename := path.Join(t.dir, key)

	if err := os.Remove(filename); err != nil {
		logrus.WithError(err).WithField("filename", filename).Error("cannot delete cache item")

		return false
	}

	return true
}

func (t *FilesystemItemPool) DeleteItems(keys []string) bool {
	success := true

	for _, key := range keys {
		if !t.DeleteItem(key) {
			success = false
		}
	}

	return success
}

func (t *FilesystemItemPool) Save(item contract.Item[[]byte]) bool {
	filename := path.Join(t.dir, item.GetKey())

	if err := os.WriteFile(filename, item.Get(), 0644); err != nil {
		logrus.WithError(err).WithField("filename", filename).Error("failed to write file cache")

		return false
	}

	return true
}

func (t *FilesystemItemPool) SaveDeferred(item contract.Item[[]byte]) bool {
	t.queue = append(t.queue, item)

	return true
}

func (t *FilesystemItemPool) Commit() bool {
	success := true

	for _, item := range t.queue {
		if !t.Save(item) {
			success = false
		}
	}

	if success {
		t.queue = make([]contract.Item[[]byte], 0)
	}

	return success
}
