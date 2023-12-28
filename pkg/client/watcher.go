package client

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sync"
	"time"

	"slai.io/takehome/pkg/common"
)

type Watcher struct {
	Root       string
	FileMap    map[string]common.FileInfo
	Delay      time.Duration
	absRoot    string
	Wg         *sync.WaitGroup
	SingleSync *common.SingleSync
}

func (w *Watcher) ResetSingleSync() {
	common.ResetSingleSync(w.SingleSync)
}

func NewWatcher(path string, delay time.Duration) *Watcher {
	absRoot, err := filepath.Abs(path)
	if err != nil {
		return nil
	}

	w := &Watcher{
		Root:    path,
		FileMap: make(map[string]common.FileInfo),
		Delay:   delay,
		absRoot: absRoot,
		SingleSync: &common.SingleSync{
			EncodeComplete: false,
			SyncComplete:   false,
			EncodeChan:     make(chan string),
			SyncChan:       make(chan *common.File),
			Wg:             &sync.WaitGroup{},
		},
	}

	return w
}

func (w *Watcher) SeenPreviously(key string) bool {
	_, ok := w.FileMap[key]
	return ok
}

func (w *Watcher) AddKey(path string, info *fs.FileInfo) {
	w.FileMap[path] = common.NewFileInfo(*info)
}

func (w *Watcher) RemoveKey(key string) {
	delete(w.FileMap, key)
}

func (w *Watcher) IsModified(key string, CurrentModTime time.Time) bool {
	// This function makes the assumption that the key already exists.
	// If you don't know if the key exists, then use the SeenPreviously() function first.
	val, ok := w.FileMap[key]
	if !ok {
		log.Fatalf("File at path (%s) doens't exist. Please make sure the file exists before checking if it's been modified.", key)
	}

	if val.CurrModTime.After(CurrentModTime) {
		return true
	}

	return false
}

func (w *Watcher) UpdateLastSyncedTime(path string, CurrentModTime time.Time) {
	val, ok := w.FileMap[path]
	if !ok {
		log.Printf("Key \"%s\" doesn't exist.", path)
		return
	}
	val.LastSyncModTime = CurrentModTime
	w.FileMap[path] = val
}

func (w *Watcher) Run() {
	var wg sync.WaitGroup
	err := filepath.WalkDir(w.absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error occured when walking directory: %s\n", err)
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			if w.SeenPreviously(path) {
				if w.IsModified(path, info.ModTime()) {
					wg.Add(1)
					go w.SingleSync.EncodeFile(path, &wg)
				}
			} else {
				w.AddKey(path, &info)
				wg.Add(1)
				go w.SingleSync.EncodeFile(path, &wg)
			}
		}
		return nil
	})

	wg.Wait()
	if err != nil {
		log.Printf("ERROR %s", err)
	}
}
