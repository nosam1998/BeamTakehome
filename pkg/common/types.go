package common

import (
	"io/fs"
	"time"
)

type RequestType string

const (
	Echo RequestType = "ECHO"
	Sync RequestType = "SYNC"
)

type BaseRequest struct {
	RequestId   string `json:"request_id"`
	RequestType string `json:"request_type"`
}

type BaseResponse struct {
	RequestId   string `json:"request_id"`
	RequestType string `json:"request_type"`
}

type EchoRequest struct {
	BaseRequest
	Value string
}

type EchoResponse struct {
	BaseResponse
	Value string
}

type FileInfo struct {
	Name    string
	Size    int64
	Mode    fs.FileMode
	ModTime time.Time
	IsDir   bool
	Sys     any
}

type File struct {
	FileInfo
	Ext     string
	Path    string
	Content *string
}

type SyncRequest struct {
	BaseRequest
	Data *File
}

type SyncResponse struct {
	BaseResponse
	Message string
}

type StartServerOpts struct {
	OutputPath string
	Address    string
	Host       string
	Port       int
}
