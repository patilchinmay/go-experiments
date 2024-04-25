package testcontainersdemo

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// func init() {
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     "0.0.0.0:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// }

func Set(ctx context.Context, client *redis.Client) error {
	err := client.Set(ctx, "foo", "bar", 0).Err()
	return err
}

func Get(ctx context.Context, client *redis.Client) string {
	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)

	return val
}
