package main

import (
	"fmt"
	"log"
	"time"

	client "slai.io/takehome/pkg/client"
)

func main() {
	log.Println("Starting client...")

	c, err := client.NewClient("./testing/client")
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
			err := c.SyncWatcherQueue()
			if err != nil {
				log.Println(err)
			}
		}
		time.Sleep(c.Watcher.Delay)
	}
}
