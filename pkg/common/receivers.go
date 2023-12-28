package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func (f *File) GetFilePath() string {
	path := fmt.Sprintf("%s\n", filepath.Join(f.Path, f.FileInfo.Name))
	return path
}

// File Receivers
func (f *File) GetDefaultServerPath() (string, error) {
	defaultPath, err := GetDefaultPath()

	if err != nil {
		return "", err
	}
	return defaultPath, err
}

func (f *File) ServerWritePath() string {
	return filepath.Base(f.Path)
}

func (f *File) FillFileInfo() error {
	stat, err := os.Stat(f.Path)
	if err != nil {
		return err
	}
	f.FileInfo = FileInfo{
		Name:        stat.Name(),
		Size:        stat.Size(),
		Mode:        stat.Mode(),
		CurrModTime: stat.ModTime(),
		IsDir:       stat.IsDir(),
		Sys:         stat.Sys(),
	}
	return nil
}

// SingleSync Receivers
func ResetSingleSync(ss *SingleSync) {
	var wg sync.WaitGroup

	ss.EncodeComplete = false
	ss.SyncComplete = false
	ss.EncodeChan = make(chan string)
	ss.SyncChan = make(chan *File)
	ss.Wg = &wg
}

func (s *SingleSync) CompleteEncode() {
	if !s.EncodeComplete {
		close(s.EncodeChan)
		s.EncodeComplete = true
	}
}

func (s *SingleSync) CompleteSync() {
	if !s.SyncComplete {
		close(s.SyncChan)
		s.SyncComplete = true
	}
}

func (s *SingleSync) IsSyncComplete() bool {
	if s.EncodeComplete && s.SyncComplete {
		return true
	}
	return false
}

func (s *SingleSync) EncodeFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	encodedFile, err := FileToBase64(path)
	if err != nil {
		log.Printf("Error occurred when encoding file (%s)", path)
	} else {
		// Pass the encoded file through the SyncChan to be synced.
		s.SyncChan <- encodedFile
	}
}
