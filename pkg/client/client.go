package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"path/filepath"
	"slai.io/takehome/pkg/common"
	"time"
)

const maxConnectionAttempts = 100
const hostURL = "ws://localhost:5555/"

type Client struct {
	Directory string
	SessionId string
	ws        *websocket.Conn
	connected bool
	hostURL   string
	channels  map[string]chan []byte
	*Watcher
}

func NewClient(directory string) (*Client, error) {
	absDir, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Watching AbsDir: %s\n", absDir)

	w := NewWatcher(absDir, time.Second*5)

	var client *Client = &Client{
		Directory: directory,
		hostURL:   hostURL,
		Watcher:   w,
	}

	err = client.connect()
	if err != nil {
		return nil, err
	}

	client.connected = true
	client.channels = make(map[string]chan []byte)

	return client, nil
}

func (c *Client) connect() error {
	connected := false
	attempts := 0

	for {
		log.Println("Connection attempt: ", attempts)

		if attempts > maxConnectionAttempts {
			break
		}

		ws, _, err := websocket.DefaultDialer.Dial(c.hostURL, nil)
		c.ws = ws

		if err != nil {
			attempts++
			continue
		}

		connected = true
		break
	}

	// We weren't able to connect to the host, bail
	if !connected {
		return nil
	}

	// Start receiving messages
	go c.rx()

	return nil
}

func (c *Client) rx() {
	for {
		_, message, err := c.ws.ReadMessage()
		var ce *websocket.CloseError
		if errors.As(err, &ce) {

			switch ce.Code {
			case websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure:
				return
			}
		}

		var msg common.BaseResponse

		err = json.Unmarshal(message, &msg)
		if err != nil {
			continue
		} else {
			if _, ok := c.channels[msg.RequestId]; ok {
				c.channels[msg.RequestId] <- message
			} else {
				log.Println("channel not found")
			}
		}
	}
}

func (c *Client) tx(msg []byte) error {
	err := c.ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SyncWatcherQueue() error {
	var retryQueue []string
	for _, path := range c.Watcher.SyncQueue {
		data, err := common.FileToBase64(path)
		if err != nil {
			retryQueue = append(retryQueue, path)
			fmt.Println(err)
			continue
		}

		serverMessage, err := c.Sync(data)
		if err != nil {
			retryQueue = append(retryQueue, path)
			fmt.Println(err)
			continue
		}
		fmt.Printf("[Server] %s\n", serverMessage)
	}
	c.Watcher.SyncQueue = retryQueue
	return nil
}

// Request implementations
func (r *Client) Echo(value string) (string, error) {
	requestId := uuid.NewString()

	var request *common.EchoRequest = &common.EchoRequest{
		BaseRequest: common.BaseRequest{
			RequestId:   requestId,
			RequestType: string(common.Echo),
		},
		Value: value,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	r.channels[requestId] = make(chan []byte)

	err = r.tx(payload)
	if err != nil {
		return "", err
	}

	var response common.EchoResponse = common.EchoResponse{}

	msg := <-r.channels[requestId]
	err = json.Unmarshal(msg, &response)
	if err != nil {
		log.Println("Unable to handle echo response: ", err)
		return "", err
	}

	return response.Value, err
}

func (r *Client) Sync(file *common.File) (string, error) {
	requestId := uuid.NewString()

	var request *common.SyncRequest = &common.SyncRequest{
		BaseRequest: common.BaseRequest{
			RequestId:   requestId,
			RequestType: string(common.Sync),
		},
		Data: file,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	r.channels[requestId] = make(chan []byte)

	err = r.tx(payload)
	if err != nil {
		return "", err
	}

	var response common.SyncResponse = common.SyncResponse{}

	msg := <-r.channels[requestId]
	err = json.Unmarshal(msg, &response)
	if err != nil {
		log.Println("Unable to handle sync response: ", err)
		return "", err
	}

	return response.Message, err
}
