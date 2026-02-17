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
	"sync"
	"time"
)

func producer(stream Stream, ch chan *Tweet, wg *sync.WaitGroup) (tweets []*Tweet) {
	defer wg.Done()
	for {
		tweet, err := stream.Next()
		if tweet != nil {
			ch <- tweet
		}
		if err == ErrEOF {
			close(ch)
			return tweets
		}
	}
}

func consumer(ch chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		t, ok := <-ch
		if !ok {
			return
		}
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

	var wg sync.WaitGroup
	ch := make(chan *Tweet)

	wg.Go(func() { producer(stream, ch, &wg) })
	wg.Go(func() { consumer(ch, &wg) })

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
