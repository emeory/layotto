package main

import (
	"context"
	"fmt"
	client "mosn.io/layotto/sdk/go-sdk/client"
)

const (
	store = "redis"
	key1  = "key1"
	key2  = "key2"
	key3  = "key3"
	key4  = "key4"
	key5  = "key5"
)

func main() {
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	ctx := context.Background()
	value := []byte("hello world")

	// Belows are CRUD examples.
	// save state
	testSave(ctx, cli, store, key1, value)

	// get state
	testGet(ctx, cli, store, key1)

	// SaveBulkState with options and metadata
	testSaveBulkState(ctx, cli, store, key1, value, key2)

	testGetBulkState(ctx, cli, store, key1, key2)

	// delete state
	testDelete(ctx, cli, store, key1)
	testDelete(ctx, cli, store, key2)
}

func testGetBulkState(ctx context.Context, cli client.Client, store string, key1 string, key2 string) {
	state, err := cli.GetBulkState(ctx, store, []string{key1, key2, key3, key4, key5}, nil, 3)
	if err != nil {
		panic(err)
	}
	for _, item := range state {
		fmt.Printf("GetBulkState succeeded.key:%v,value:%v\n", item.Key, string(item.Value))
	}
}

func testDelete(ctx context.Context, cli client.Client, store string, key string) {
	if err := cli.DeleteState(ctx, store, key); err != nil {
		panic(err)
	}
	fmt.Printf("DeleteState succeeded.key:%v\n", key)
}

func testSaveBulkState(ctx context.Context, cli client.Client, store string, key string, value []byte, key2 string) {
	item := &client.SetStateItem{
		Etag: &client.ETag{
			Value: "2",
		},
		Key: key,
		Metadata: map[string]string{
			"some-key-for-component": "some-value",
		},
		Value: value,
		Options: &client.StateOptions{
			Concurrency: client.StateConcurrencyLastWrite,
			Consistency: client.StateConsistencyStrong,
		},
	}
	item2 := *item
	item2.Key = key2

	if err := cli.SaveBulkState(ctx, store, item, &item2); err != nil {
		panic(err)
	}
	fmt.Printf("SaveBulkState succeeded.[key:%s etag:%s]: %s\n", item.Key, item.Etag.Value, string(item.Value))
	fmt.Printf("SaveBulkState succeeded.[key:%s etag:%s]: %s\n", item2.Key, item2.Etag.Value, string(item2.Value))
}

func testGet(ctx context.Context, cli client.Client, store string, key string) {
	item, err := cli.GetState(ctx, store, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetState succeeded.[key:%s etag:%s]: %s\n", item.Key, item.Etag, string(item.Value))
}

func testSave(ctx context.Context, cli client.Client, store string, key string, value []byte) {
	if err := cli.SaveState(ctx, store, key, value); err != nil {
		panic(err)
	}
	fmt.Printf("SaveState succeeded.key:%v , value: %v \n", key, string(value))
}
