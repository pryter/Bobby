package events

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func generateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func testEvent(length time.Duration) {
	time.Sleep(length * time.Second)
	fmt.Printf("Task %s done laews\n", generateRandomString(5))
}

func TestConcurrentPool(t *testing.T) {
	pool := InitConcurrentPool(ConcurrentPoolOptions{MaxConcurrentTasks: 3})
	pool.Add(func() {
		testEvent(10)
	})
	pool.Add(func() {
		testEvent(10)
	})
	pool.Add(func() {
		testEvent(10)
	})
	pool.Add(func() {
		testEvent(10)
	})
	pool.Add(func() {
		testEvent(10)
	})

	for true {
	}
}
