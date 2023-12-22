package client

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slai.io/takehome/pkg/common"
	"sync"
	"time"
)

type Watcher struct {
	Root      string
	Previous  map[string]common.FileInfo
	Delay     time.Duration
	SyncQueue []string
	absRoot   string
	wg        sync.WaitGroup
}

func NewWatcher(path string, delay time.Duration) *Watcher {
	absRoot, err := filepath.Abs(path)
	if err != nil {
		return nil
	}
	w := &Watcher{
		Root:      path,
		Previous:  make(map[string]common.FileInfo),
		Delay:     delay,
		SyncQueue: make([]string, 0),
		absRoot:   absRoot,
		wg:        sync.WaitGroup{},
	}

	return w
}

func (w *Watcher) walkFunc(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Printf("error occured when walking directory: %s\n", err)
		return err
	}

	if !d.IsDir() {
		info, err := d.Info()
		if err != nil {
			return err
		}
		if w.IsNew(path, &info) {
			w.SyncQueue = append(w.SyncQueue, path)
		}
	}
	return nil
}

func (w *Watcher) SeenPreviously(key string) bool {
	_, ok := w.Previous[key]
	return ok
}

func (w *Watcher) AddKey(path string, info *fs.FileInfo) {
	w.Previous[path] = common.NewFileInfo(*info)
}

func (w *Watcher) RemoveKey(key string) {
	delete(w.Previous, key)
}

func (w *Watcher) IsNew(key string, info *fs.FileInfo) bool {
	// Returns true if we've never seen the key before (New file)
	// OR if the file has been seen, but has a more recent modified time.
	val, ok := w.Previous[key]
	fi := common.NewFileInfo(*info)
	if ok && fi.ModTime.After(val.ModTime) {
		fmt.Println("New version of ", key)
		return true
	}

	if !ok {
		w.AddKey(key, info)
		fmt.Println("New file ", key)
		return true
	}

	return false
}

func (w *Watcher) Run() error {
	err := filepath.WalkDir(w.absRoot, w.walkFunc)
	if err != nil {
		return err
	}
	return nil
}
