package common

import (
	"io/fs"
	"sync"
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
	Name            string
	Size            int64
	Mode            fs.FileMode
	CurrModTime     time.Time
	LastSyncModTime time.Time
	IsDir           bool
	Sys             any
}

type File struct {
	FileInfo
	Key     string
	Ext     string
	Path    string
	Content *string
}

type SingleSync struct {
	EncodeComplete bool
	SyncComplete   bool
	EncodeChan     chan string
	SyncChan       chan *File
	Wg             *sync.WaitGroup
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
