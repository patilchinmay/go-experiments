package testcontainersdemo_test

import (
	"context"
	"log"
	"testing"

	"github.com/patilchinmay/go-experiments/testcontainersdemo"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWithRedis(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "valkey/valkey:7.2.5",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}

	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop redis: %s", err)
		}
	}()

	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		t.Error(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	t.Run("Set foo bar", func(t *testing.T) {
		err := testcontainersdemo.Set(ctx, client)

		assert.NoError(t, err)
	})

	t.Run("Get foo", func(t *testing.T) {
		val := testcontainersdemo.Get(ctx, client)

		assert.Equal(t, "bar", val)
	})
}
