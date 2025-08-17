package redis_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/core/memorystore"
	"github.com/wal1251/pkg/providers/redis"
)

func TestClient(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T, *redis.Client)
	}{
		// Get and Set
		{
			name: "Set and Get string value",
			test: func(t *testing.T, client *redis.Client) {
				err := client.Set(context.Background(), "key1", "value1", 0)
				assert.NoError(t, err)

				val, err := client.Get(context.Background(), "key1")
				assert.NoError(t, err)

				valStr, err := val.String()
				assert.NoError(t, err)

				assert.Equal(t, "value1", valStr)
			},
		},
		{
			name: "Set and Get int value",
			test: func(t *testing.T, client *redis.Client) {
				err := client.Set(context.Background(), "key2", 12345, 0)
				assert.NoError(t, err)

				val, err := client.Get(context.Background(), "key2")
				assert.NoError(t, err)

				valInt, err := val.Int()
				assert.NoError(t, err)

				assert.Equal(t, 12345, valInt)
			},
		},
		{
			name: "Set and Get bool value",
			test: func(t *testing.T, client *redis.Client) {
				err := client.Set(context.Background(), "key3", true, 0)
				assert.NoError(t, err)

				val, err := client.Get(context.Background(), "key3")
				assert.NoError(t, err)

				valBool, err := val.Bool()
				assert.NoError(t, err)

				assert.Equal(t, true, valBool)
			},
		},
		{
			name: "Set and Get float64 value",
			test: func(t *testing.T, client *redis.Client) {
				err := client.Set(context.Background(), "key4", 123.45, 0)
				assert.NoError(t, err)

				val, err := client.Get(context.Background(), "key4")
				assert.NoError(t, err)

				valFloat, err := val.Float64()
				assert.NoError(t, err)

				assert.Equal(t, 123.45, valFloat)
			},
		},
		{
			name: "Set and Get time value",
			test: func(t *testing.T, client *redis.Client) {
				now := time.Now()
				err := client.Set(context.Background(), "key5", now, 0)
				assert.NoError(t, err)

				val, err := client.Get(context.Background(), "key5")
				assert.NoError(t, err)

				valTime, err := val.Time()
				assert.NoError(t, err)
				assert.ObjectsAreEqualValues(now, valTime)
			},
		},
		{
			name: "Set and Get struct",
			test: func(t *testing.T, client *redis.Client) {
				type TestStruct struct {
					Field1 string
					Field2 int
				}

				testStruct := TestStruct{
					Field1: "value1",
					Field2: 12345,
				}

				err := client.Set(context.Background(), "key5", testStruct, 0)
				assert.NoError(t, err)

				val, err := client.Get(context.Background(), "key5")
				assert.NoError(t, err)

				var valStruct TestStruct
				err = val.Struct(&valStruct)
				assert.NoError(t, err)

				assert.Equal(t, testStruct, valStruct)
			},
		},
		{
			name: "Get non-existing value",
			test: func(t *testing.T, client *redis.Client) {
				val, err := client.Get(context.Background(), "key1")
				assert.ErrorIs(t, err, memorystore.ErrKeyNotFound)
				assert.Nil(t, val)
			},
		},

		// GetList
		{
			name: "Get multiple values",
			test: func(t *testing.T, client *redis.Client) {
				keys := []string{"key1", "key2", "key3"}
				for _, key := range keys {
					err := client.Set(context.Background(), key, key, 0)
					assert.NoError(t, err)
				}

				vals, err := client.GetList(context.Background(), keys...)
				assert.NoError(t, err)

				for i, val := range vals {
					strVal, err := val.String()
					assert.NoError(t, err)
					assert.Equal(t, keys[i], strVal)
				}
			},
		},
		{
			name: "Get non-existing list",
			test: func(t *testing.T, client *redis.Client) {
				vals, err := client.GetList(context.Background(), "key1", "key2", "key3")
				assert.NoError(t, err)
				assert.Len(t, vals, 3)

				for _, val := range vals {
					assert.Nil(t, val)
				}
			},
		},
		{
			name: "One of the keys is non-existing in the list",
			test: func(t *testing.T, client *redis.Client) {
				for _, key := range []string{"key1", "key3"} {
					err := client.Set(context.Background(), key, key, 0)
					assert.NoError(t, err)
				}

				vals, err := client.GetList(context.Background(), "key1", "key2", "key3")
				assert.NoError(t, err)
				assert.Len(t, vals, 3)

				key1StrVal, err := vals[0].String()
				assert.NoError(t, err)
				key3StrVal, err := vals[2].String()
				assert.NoError(t, err)

				assert.Equal(t, "key1", key1StrVal)
				assert.Nil(t, vals[1])
				assert.Equal(t, "key3", key3StrVal)
			},
		},
		{
			name: "Get list keys count is greater than MaxBulkRequestSize",
			test: func(t *testing.T, client *redis.Client) {
				keys := make([]string, memorystore.DefaultMaxBulkRequestSize+1)
				for i := 0; i < memorystore.DefaultMaxBulkRequestSize+1; i++ {
					keys[i] = fmt.Sprintf("key%d", i)
				}

				vals, err := client.GetList(context.Background(), keys...)
				assert.ErrorIs(t, err, memorystore.ErrBulkRequestTooLarge)
				assert.Nil(t, vals)
			},
		},

		// Delete
		{
			name: "Delete non-existing value",
			test: func(t *testing.T, client *redis.Client) {
				deleted, err := client.Delete(context.Background(), "key1")
				assert.NoError(t, err)
				assert.Equal(t, 0, deleted)
			},
		},
		{
			name: "Delete value",
			test: func(t *testing.T, client *redis.Client) {
				err := client.Set(context.Background(), "key2", "value2", 0)
				assert.NoError(t, err)

				deleted, err := client.Delete(context.Background(), "key2")
				assert.NoError(t, err)
				assert.Equal(t, 1, deleted)
			},
		},
		{
			name: "Delete multiple values",
			test: func(t *testing.T, client *redis.Client) {
				keys := make([]string, 10)
				for i := 0; i < 10; i++ {
					keys[i] = fmt.Sprintf("key%d", i)
				}

				for _, key := range keys {
					err := client.Set(context.Background(), key, "value", 0)
					assert.NoError(t, err)
				}

				deleted, err := client.Delete(context.Background(), keys...)
				assert.NoError(t, err)
				assert.Equal(t, 10, deleted)
			},
		},
		{
			name: "Deleting keys count more than MaxBulkRequestSize",
			test: func(t *testing.T, client *redis.Client) {
				keys := make([]string, memorystore.DefaultMaxBulkRequestSize+1)
				for i := 0; i < memorystore.DefaultMaxBulkRequestSize+1; i++ {
					keys[i] = fmt.Sprintf("key%d", i)
				}

				deleted, err := client.Delete(context.Background(), keys...)
				assert.ErrorIs(t, err, memorystore.ErrBulkRequestTooLarge)
				assert.Equal(t, 0, deleted)
			},
		},
	}

	config := &redis.Config{
		Host:               "localhost",
		Port:               "6379",
		MaxBulkRequestSize: memorystore.DefaultMaxBulkRequestSize,
	}

	// Инициализация тестового сервера Redis.
	testServer := redis.NewTestRedisServer()
	err := testServer.Run(*config)
	defer testServer.Close()
	require.NoError(t, err)

	// Инициализация клиента Redis.
	client, err := redis.NewClient(context.Background(), config)
	defer client.Close(context.Background())
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Удаление всех ключей перед каждым тестом.
			testServer.FlushAll()
			tt.test(t, client)
		})
	}
}
