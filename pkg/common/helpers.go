package common

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func EncodeToBase64(data *[]byte) *string {
	encodedData := base64.StdEncoding.EncodeToString(*data)
	return &encodedData
}

func DecodeBase64(data *string) (*[]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(*data)
	if err != nil {
		return nil, err
	}
	return &decodedData, nil
}

func WriteFile(file *File, path string) error {
	outputRoot := path
	root, err := filepath.Abs(outputRoot)
	if err != nil {
		return err
	}

	relFilePath := filepath.Join(root, file.Name)
	fmt.Println("Wrote file: ", relFilePath)

	filePtr, err := os.OpenFile(relFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, file.Mode)
	if err != nil {
		return err
	}

	content, err := DecodeBase64(file.Content)
	if err != nil {
		return err
	}

	_, err = filePtr.Write(*content)
	if err != nil {
		return err
	}

	err = filePtr.Close()
	if err != nil {
		return err
	}

	return nil
}

func NewFileInfo(info fs.FileInfo) FileInfo {
	return FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
		Sys:     info.Sys(),
	}
}

func FileToBase64(path string) (*File, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		e := fmt.Errorf("not a file")
		return nil, e
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	encodedData := EncodeToBase64(&data)

	return &File{
		FileInfo: NewFileInfo(fileInfo),
		Ext:      filepath.Ext(path),
		Path:     filepath.Dir(path),
		Content:  encodedData,
	}, nil
}

func MakeDir(path string) error {
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}

	return path, nil
}
