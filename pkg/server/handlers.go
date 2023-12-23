package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"slai.io/takehome/pkg/common"
)

func HandleEcho(msg []byte, client *Client) error {
	log.Println("Received ECHO request.")

	var request common.EchoRequest
	err := json.Unmarshal(msg, &request)

	if err != nil {
		log.Fatal("Invalid echo request.")
	}

	response := &common.EchoResponse{
		BaseResponse: common.BaseResponse{
			RequestId:   request.RequestId,
			RequestType: request.RequestType,
		},
		Value: request.Value,
	}

	responsePayload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = client.ws.WriteMessage(websocket.TextMessage, responsePayload)
	if err != nil {
		return err
	}

	return nil
}

func HandleSync(msg []byte, client *Client) error {
	var request common.SyncRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		log.Fatal("Invalid sync request.")
	}

	go func() {
		err := common.WriteFile(request.Data, OutputDir)
		if err != nil {
			fmt.Printf("Error writing file %s\n", request.Data.Path)
			return
		}
		log.Println("[SYNC] ", request.Data.FileInfo.Name)
	}()
	if err != nil {
		return err
	}

	response := &common.SyncResponse{
		BaseResponse: common.BaseResponse{
			RequestId:   request.RequestId,
			RequestType: request.RequestType,
		},
		Message: fmt.Sprintf("Syncing file \"%s\"", request.Data.Path),
	}

	responsePayload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = client.ws.WriteMessage(websocket.TextMessage, responsePayload)
	if err != nil {
		return err
	}

	return nil
}
