package main

import (
	"log"
	"sync"
	"time"

	client "slai.io/takehome/pkg/client"
)

func main() {
	log.Println("Starting client...")

	c, err := client.NewClient("./testing/client")
	if err != nil {
		log.Fatal(err)
	}

	if c == nil {
		log.Fatal("Error creating client")
	}

	log.Printf("Watching - Waiting for %d seconds between checks...\n", int(c.Watcher.Delay.Seconds()))
	for {
		var wg sync.WaitGroup

		wg.Add(2)
		go func() {
			defer wg.Done()
			defer c.Watcher.SingleSync.CompleteSync()
			c.SyncFromChannel()
		}()

		go func() {
			defer wg.Done()
			defer c.Watcher.SingleSync.CompleteEncode()
			c.Watcher.Run()
		}()

		time.Sleep(c.Watcher.Delay)
	}
}
