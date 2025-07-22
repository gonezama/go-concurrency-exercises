//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
	"sync"
)

func producer(stream Stream, t chan <- *Tweet) {
	defer close(t)
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return
		}
		t<-tweet
	}
}

func consumer(tweets <-chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	
	t := make(chan *Tweet)
	var wg sync.WaitGroup
	wg.Add(1)

	// Producer
	go producer(stream, t)

	// Consumer
	go consumer(t, &wg)
	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))
}
