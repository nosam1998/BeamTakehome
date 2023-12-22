package common

import (
	"fmt"
	"os"
	"path/filepath"
)

// File Receivers
func (f *File) getClientPath() (string, error) {
	defaultPath, err := GetDefaultPath()
	if err != nil {
		return "", err
	}
	fmt.Printf("Default path is: %s", defaultPath)
	return defaultPath, err
}

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
		Name:    stat.Name(),
		Size:    stat.Size(),
		Mode:    stat.Mode(),
		ModTime: stat.ModTime(),
		IsDir:   stat.IsDir(),
		Sys:     stat.Sys(),
	}
	return nil
}
