package main

import (
	"log"
	client "slai.io/takehome/pkg/client"
	"sync"
	"time"
)

func main() {
	log.Println("Starting client...")

	c, err := client.NewClient("./testing/client")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Watching - Waiting for %d seconds between checks...\n", int(c.Watcher.Delay.Seconds()))
	for {
		var wg sync.WaitGroup

		go func() {
			wg.Add(1)
			defer wg.Done()
			defer c.Watcher.SingleSync.CompleteSync()
			c.SyncFromChannel()
		}()

		go func() {
			wg.Add(1)
			defer wg.Done()
			defer c.Watcher.SingleSync.CompleteEncode()
			c.Watcher.Run()
		}()

		wg.Wait()
		time.Sleep(c.Watcher.Delay)
	}
}
