package main

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	// Create new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Test Redis connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		errorString := fmt.Sprintf("Connection to Redis failed: %s", err)
		panic(errorString)
	}
	fmt.Println("Connection to Redis successfully established")

	// Config
	url := "https://example.com/"
	numberOfRequests := 200_000

	// Hash url and save hash in Redis
	MD5Hash := CreateMD5HashShortened(url, 4)
	err = client.Set(ctx, MD5Hash, url, 0).Err()
	if err != nil {
		panic(err)
	}

	performSequentialTests(numberOfRequests, MD5Hash, client)
	performConcurrentTests(numberOfRequests, MD5Hash, client)
}

func CreateMD5HashShortened(text string, byteLength int) string {
	hash := md5.Sum([]byte(text))

	// Check if the length of the hash should be limited
	if byteLength > 0 {
		return hex.EncodeToString(hash[:byteLength])
	} else {

		return hex.EncodeToString(hash[:])
	}
}

func CreateSHA1HashShortened(text string, byteLength int) string {
	hash := sha1.New()
	hash.Write([]byte(text))

	// Check if the length of the hash should be limited
	if byteLength > 0 {
		return hex.EncodeToString(hash.Sum(nil)[:byteLength])
	} else {
		return hex.EncodeToString(hash.Sum(nil))
	}
}

func performSequentialTests(numberOfRequests int, MD5Hash string, client *redis.Client) {
	startTime := time.Now()
	for i := 0; i < numberOfRequests; i++ {
		client.Get(ctx, MD5Hash)
	}
	elapsedTime := time.Since(startTime)
	fmt.Printf("Looking up %d MD5 hashes sequentially took %.4ss\n", numberOfRequests, elapsedTime)
}

func performConcurrentTests(numberOfRequests int, MD5Hash string, client *redis.Client) {
	startTime := time.Now()

	// Create wait group to check for completion of all concurrently run function
	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())

	// Create x numbers of groups, i.e. as many CPUs you have
	for group := 0; group < runtime.NumCPU(); group++ {
		// Execute function calls concurrently
		go func(group int) {
			for i := group * (numberOfRequests / runtime.NumCPU()); i < (group+1)*(numberOfRequests/runtime.NumCPU()); i++ {
				client.Get(ctx, MD5Hash)
			}

			// Signal the wait group that this function call has been finished
			wg.Done()
		}(group)
	}

	// Wait until all functions in this wait group have been finished
	wg.Wait()
	elapsedTime := time.Since(startTime)
	fmt.Printf("Looking up %d MD5 hashes using go routines took %.4ss\n", numberOfRequests, elapsedTime)
}
