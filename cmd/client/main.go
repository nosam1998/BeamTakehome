package main

import (
	"fmt"
	"log"
	"time"

	client "slai.io/takehome/pkg/client"
	"slai.io/takehome/pkg/common"
)

func main() {
	log.Println("Starting client...")

	c, err := client.NewClient("./")
	if err != nil {
		log.Fatal(err)
	}

	someMessage := "hello there"

	// Test variables soon...
	alreadySent := false
	filePath := "/home/mason/code/BeamTakehome/cmd/client/testfile.txt"

	// TODO: Implement file watcher
	for {
		log.Printf("Sending: '%s'", filePath)

		if !alreadySent {
			data, err := common.FileToBase64(filePath)
			if err != nil {
				log.Fatal(err)
			}
			serverMessage, err := c.Sync(data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("[Server] %s", serverMessage)
			alreadySent = true
		} else {
			log.Printf("Sending: '%s'", someMessage)

			value, err := c.Echo(someMessage)
			if err != nil {
				log.Fatal("Unable to send request.")
			}

			log.Printf("Received: '%s'", value)
		}

		time.Sleep(time.Second)
	}

}
