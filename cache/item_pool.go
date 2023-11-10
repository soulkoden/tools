package cache

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

type ItemPool struct {
	dir string

	queue []*Item
}

func NewItemPool(dir string) (*ItemPool, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("cache directory is unreadable: %w", err)
	}

	if !stat.IsDir() {
		return nil, errors.New("cache directory is not a directory")
	}

	return &ItemPool{
		dir:   dir,
		queue: make([]*Item, 0),
	}, nil
}

func (t *ItemPool) GetItem(key string) *Item {
	return NewItem(t.dir, key)
}

func (t *ItemPool) GetItems(keys []string) []*Item {
	var items = make([]*Item, len(keys))
	for i, key := range keys {
		items[i] = t.GetItem(key)
	}

	return items
}

func (t *ItemPool) HasItem(key string) bool {
	filename := path.Join(t.dir, key)

	if _, err := os.Stat(filename); err != nil {
		if !os.IsNotExist(err) {
			logrus.WithError(err).WithField("filename", filename).Errorf("failed to get file stat")
		}

		return false
	}

	return true
}

func (t *ItemPool) Clear() bool {
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

func (t *ItemPool) DeleteItem(key string) bool {
	filename := path.Join(t.dir, key)

	if err := os.Remove(filename); err != nil {
		logrus.WithError(err).WithField("filename", filename).Error("cannot delete cache item")

		return false
	}

	return true
}

func (t *ItemPool) DeleteItems(keys []string) bool {
	success := true

	for _, key := range keys {
		if !t.DeleteItem(key) {
			success = false
		}
	}

	return success
}

func (t *ItemPool) Save(item *Item) bool {
	filename := path.Join(t.dir, item.GetKey())

	if err := os.WriteFile(filename, item.Get(), 0644); err != nil {
		logrus.WithError(err).WithField("filename", filename).Error("failed to write file cache")

		return false
	}

	return true
}

func (t *ItemPool) SaveDeferred(item *Item) bool {
	t.queue = append(t.queue, item)

	return true
}

func (t *ItemPool) Commit() bool {
	success := true

	for _, item := range t.queue {
		if !t.Save(item) {
			success = false
		}
	}

	return success
}
