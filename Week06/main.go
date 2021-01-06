package main

import (
	"fmt"
	"time"

	"geektime/Go-000/Week06/rolling"
)

func main() {
	opt := rolling.RollingCounterOpts{Size: 4, BucketDuration: 100 * time.Millisecond}
	rollingCounter := rolling.NewRollingCounter(opt)

	rollingCounter.Add(4)
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("rolling counter avg:%f,Value:%d\n", rollingCounter.Avg(), rollingCounter.Value())
	rollingCounter.Add(6)
	fmt.Printf("rolling counter avg:%f,Value:%d\n", rollingCounter.Avg(), rollingCounter.Value())

	time.Sleep(200 * time.Millisecond)
	rollingCounter.Add(1)
	rollingCounter.Add(3)
	fmt.Printf("rolling counter avg:%f,Value:%d\n", rollingCounter.Avg(), rollingCounter.Value())

	time.Sleep(100 * time.Millisecond)
	rollingCounter.Add(10)
	fmt.Printf("rolling counter avg:%f,Value:%d\n", rollingCounter.Avg(), rollingCounter.Value())

	time.Sleep(500 * time.Millisecond)
	rollingCounter.Add(100)
	fmt.Printf("rolling counter avg:%f,Value:%d\n", rollingCounter.Avg(), rollingCounter.Value())
}
