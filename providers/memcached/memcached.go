// Package memcached.
// Представляет собой адаптер memcached-клиента.
package memcached

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/core/memorystore"
	"github.com/wal1251/pkg/tools/serial"
)

var _ memorystore.Manager = (*Client)(nil)

// Client является клиентом для работы с Memcached.
type Client struct {
	client *memcache.Client // Низкоуровневый клиент Memcached.
	cfg    *Config
}

// NewClient инициализирует и возвращает новый экземпляр клиента Memcached.
// Проверяет доступность Memcached и устанавливает соединение.
// В случае, когда Memcached недоступен, возвращает ошибку.
func NewClient(_ context.Context, cfg *Config) (*Client, error) {
	client := memcache.New(cfg.Hosts...)

	err := client.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to memcache: %w", err)
	}

	return &Client{client: client, cfg: cfg}, nil
}

// Set сохраняет значение по ключу в Memcached с заданным временем истечения.
func (c *Client) Set(_ context.Context, key string, value any, expiration time.Duration) error {
	data, err := serial.ToBytes(value, serial.JSONEncode[any])
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: int32(expiration.Seconds()),
	}

	err = c.client.Set(item)
	if err != nil {
		return fmt.Errorf("can't set memcache key %s: %w", key, err)
	}

	return nil
}

// Get извлекает и возвращает значение по ключу из Memcached.
// Если ключ не существует, возвращает ошибку memorystore.ErrKeyNotFound.
func (c *Client) Get(_ context.Context, key string) (*memorystore.Value, error) {
	item, err := c.client.Get(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return nil, memorystore.ErrKeyNotFound
		}

		return nil, fmt.Errorf("can't get memcache key %s: %w", key, err)
	}

	return memorystore.NewValue(item.Value), nil
}

// GetList извлекает и возвращает значения для списка ключей из Memcached.
// Ограничивает количество ключей, которые можно запросить за один раз, в соответствии с конфигурацией.
// Для несуществующих ключей в соответствующих позициях будет nil.
// Если превышено кол-во запрашиваемых элементов, возвращает ошибку memorystore.ErrBulkRequestTooLarge.
func (c *Client) GetList(_ context.Context, keys ...string) ([]*memorystore.Value, error) {
	if len(keys) == 0 {
		return []*memorystore.Value{}, nil
	}

	if len(keys) > c.cfg.MaxBulkRequestSize {
		return nil, fmt.Errorf("can't get more than %d keys at once: %w", c.cfg.MaxBulkRequestSize, memorystore.ErrBulkRequestTooLarge)
	}

	items, err := c.client.GetMulti(keys)
	if err != nil {
		return nil, fmt.Errorf("can't get memcache keys %v: %w", keys, err)
	}

	results := make([]*memorystore.Value, len(keys))
	for idx, key := range keys {
		item, ok := items[key]
		if !ok {
			results[idx] = nil

			continue
		}

		results[idx] = memorystore.NewValue(item.Value)
	}

	return results, nil
}

// Delete удаляет один или несколько ключей из Redis и возвращает количество успешно удаленных ключей.
// Если превышено кол-во удаляемых элементов, возвращает ошибку memorystore.ErrBulkRequestTooLarge.
func (c *Client) Delete(_ context.Context, keys ...string) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	if len(keys) > c.cfg.MaxBulkRequestSize {
		return 0, fmt.Errorf("can't delete more than %d keys at once: %w", c.cfg.MaxBulkRequestSize, memorystore.ErrBulkRequestTooLarge)
	}

	var deleted int
	for _, key := range keys {
		err := c.client.Delete(key)
		if err != nil {
			if errors.Is(err, memcache.ErrCacheMiss) {
				continue
			}

			return deleted, fmt.Errorf("can't delete memcache key %s: %w", key, err)
		}
		deleted++
	}

	return deleted, nil
}

// Close закрывает соединение клиента с Memcached, логируя при этом информацию о закрытии.
func (c *Client) Close(ctx context.Context) {
	logs.FromContext(ctx).Debug().Msg("closing memcache client")
	if err := c.client.Close(); err != nil {
		logs.FromContext(ctx).Warn().Err(err).Msgf("can't close memcache client")
	}
}
