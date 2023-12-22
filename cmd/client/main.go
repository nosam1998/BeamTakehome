package main

import (
	"fmt"
	"log"
	"slai.io/takehome/pkg/common"
	"time"

	client "slai.io/takehome/pkg/client"
)

func main() {
	log.Println("Starting client...")

	c, err := client.NewClient("./testing/clientdir")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Watching - Waiting for %d seconds between checks...\n", int(c.Watcher.Delay.Seconds()))

	for {
		err := c.Watcher.Run()
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(c.Watcher.SyncQueue) > 0 {
			for _, path := range c.Watcher.SyncQueue {
				data, err := common.FileToBase64(path)
				if err != nil {
					fmt.Println(err)
				} else {
					serverMessage, err := c.Sync(data)
					if err != nil {
						return
					}
					fmt.Printf("[Server] %s\n", serverMessage)
				}
			}
			c.Watcher.SyncQueue = []string{}
		}
		time.Sleep(c.Watcher.Delay)
	}
}
