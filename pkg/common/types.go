package common

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

type File struct {
	Name    string
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
