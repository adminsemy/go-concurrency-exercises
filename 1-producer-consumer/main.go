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
)

func producer(stream Stream) <-chan *Tweet {
	tweetCH := make(chan *Tweet)

	go func() {
		defer close(tweetCH)
		for {
			tweet, err := stream.Next()
			if err != nil {
				return
			}
			tweetCH <- tweet
		}
	}()
	return tweetCH
}

func consumer(tweets <-chan *Tweet) {
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

	// Producer
	tweets := producer(stream)

	// Consumer
	consumer(tweets)
	fmt.Printf("Process took %s\n", time.Since(start))
}
