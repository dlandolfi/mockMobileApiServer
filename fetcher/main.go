package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func fetchPriceObject() (string, error) {
	url := "https://api.metals.dev/v1/metal/spot"
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("api_key", "1XXURGAUJCZZAWTFJPHB808TFJPHB")
	q.Add("metal", "gold")
	q.Add("currency", "USD")

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	priceObject := string(body)
	return priceObject, nil
}

func main() {
	// ExampleClient()
	response, err := fetchPriceObject()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(response))
}

func ExampleClient() {

	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
