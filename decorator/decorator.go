package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

type piFunc func(int) float64

func wraplogger(f piFunc, logger *log.Logger) piFunc {
	return func(n int) float64 {
		fn := func(n int) (result float64) {
			defer func(t time.Time) {
				logger.Printf("Took=%v, n=%v, result=%v", time.Since(t), n, result)
			}(time.Now())

			return f(n)
		}
		return fn(n)
	}
}

func wrapcache(f piFunc, cache *sync.Map) piFunc {
	return func(n int) float64 {
		fn := func(n int) float64 {
			key := fmt.Sprintf("n=%d", n)
			val, ok := cache.Load(key)
			if ok {
				return val.(float64)
			}
			result := f(n)
			cache.Store(key, result)
			return result
		}

		return fn(n)
	}
}

func Pi(n int) float64 {
	ch := make(chan float64)

	for k := 0; k <= n; k++ {
		go func(ch chan float64, k float64) {
			ch <- 4 * math.Pow(-1, k) / (2*k + 1)
		}(ch, float64(k))
	}

	result := 0.0
	for k := 0; k <= n; k++ {
		result += <-ch
	}

	return result
}

func main() {
	fmt.Println(Pi(1000))
	fmt.Println(Pi(50000))

	f := wrapcache(Pi, &sync.Map{})
	g := wraplogger(f, log.New(os.Stdout, "test ", 1))
	g(100000)
	g(20000)
	g(100000)
}
