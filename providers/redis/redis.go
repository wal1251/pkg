// Package redis.
// Представляет собой адаптер redis-клиента.
package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"time"

	rv9 "github.com/redis/go-redis/v9"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/core/memorystore"
	"github.com/wal1251/pkg/tools/serial"
)

var _ memorystore.Manager = (*Client)(nil)

// Client является клиентом для работы с Redis.
// Включает в себя поддержку TLS для защищенных соединений и опциональную
// конфигурацию для работы с Redis через кластер.
type Client struct {
	client *rv9.Client // Низкоуровневый клиент Redis.
	cfg    *Config
}

// NewClient инициализирует и возвращает новый экземпляр клиента Redis.
// Проверяет доступность Redis и устанавливает соединение с учетом TLS и кластерной конфигурации.
// В случае, когда Redis недоступен, возвращает ошибку.
func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	rdb, err := createRedisClient(cfg)
	if err != nil {
		return nil, err
	}

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis is not responding: %w", err)
	}

	return &Client{client: rdb, cfg: cfg}, nil
}

// createRedisClient создает и возвращает клиент Redis с учетом конфигурации.
// Поддерживает как стандартное подключение, так и подключение через Sentinel для кластеров.
func createRedisClient(cfg *Config) (*rv9.Client, error) {
	var tlsConfig *tls.Config
	var err error
	if cfg.CertCA != "" {
		tlsConfig, err = getTLSConfig(cfg.CertCA)
		if err != nil {
			return nil, err
		}
	}

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	if cfg.ClusterEnabled {
		opt := &rv9.FailoverOptions{
			SentinelAddrs: []string{addr},
			Password:      cfg.Password,
			DB:            cfg.Database,
			MasterName:    cfg.MasterName,
			TLSConfig:     tlsConfig,
		}

		return rv9.NewFailoverClient(opt), nil
	}

	opt := &rv9.Options{
		Addr:      addr,
		Password:  cfg.Password,
		DB:        cfg.Database,
		TLSConfig: tlsConfig,
	}

	return rv9.NewClient(opt), nil
}

// getTLSConfig создает и возвращает настроенный TLS конфиг на основе пути к сертификату.
func getTLSConfig(certPath string) (*tls.Config, error) {
	caCert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("error while reading Redis certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{RootCAs: caCertPool, MinVersion: tls.VersionTLS12}, nil
}

// Set сохраняет значение по ключу в Redis с заданным временем истечения.
func (r *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	data, err := serial.ToBytes(value, serial.JSONEncode[any])
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return fmt.Errorf("can't set redis key %s: %w", key, err)
	}

	return nil
}

// Get извлекает и возвращает значение по ключу из Redis.
// Если ключ не существует, возвращает ошибку memorystore.ErrKeyNotFound.
func (r *Client) Get(ctx context.Context, key string) (*memorystore.Value, error) {
	cmd := r.client.Get(ctx, key)

	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), rv9.Nil) {
			return nil, memorystore.ErrKeyNotFound
		}

		return nil, fmt.Errorf("can't get redis key %s: %w", key, cmd.Err())
	}

	data, err := cmd.Bytes()
	if err != nil {
		return nil, fmt.Errorf("can't read redis key %s value: %w", key, err)
	}

	return memorystore.NewValue(data), nil
}

// GetList извлекает и возвращает значения для списка ключей из Redis.
// Ограничивает количество ключей, которые можно запросить за один раз, в соответствии с конфигурацией.
// Для несуществующих ключей в соответствующих позициях будет nil.
// Если превышено кол-во запрашиваемых элементов, возвращает ошибку memorystore.ErrBulkRequestTooLarge.
func (r *Client) GetList(ctx context.Context, keys ...string) ([]*memorystore.Value, error) {
	if len(keys) == 0 {
		return []*memorystore.Value{}, nil
	}

	if len(keys) > r.cfg.MaxBulkRequestSize {
		return nil, fmt.Errorf("can't get more than %d keys at once: %w", r.cfg.MaxBulkRequestSize, memorystore.ErrBulkRequestTooLarge)
	}

	results := make([]*memorystore.Value, len(keys))
	for idx, k := range keys {
		val, err := r.Get(ctx, k)
		if err != nil {
			if errors.Is(err, memorystore.ErrKeyNotFound) {
				results[idx] = nil

				continue
			}

			return nil, err
		}

		results[idx] = val
	}

	return results, nil
}

// Delete удаляет один или несколько ключей из Redis и возвращает количество успешно удаленных ключей.
// Если превышено кол-во удаляемых элементов, возвращает ошибку memorystore.ErrBulkRequestTooLarge.
func (r *Client) Delete(ctx context.Context, keys ...string) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	if len(keys) > r.cfg.MaxBulkRequestSize {
		return 0, fmt.Errorf("can't delete more than %d keys at once: %w", r.cfg.MaxBulkRequestSize, memorystore.ErrBulkRequestTooLarge)
	}

	result := r.client.Del(ctx, keys...)

	return int(result.Val()), result.Err()
}

// Close закрывает соединение клиента с Redis, логируя при этом информацию о закрытии.
func (r *Client) Close(ctx context.Context) {
	logs.FromContext(ctx).Debug().Msg("closing redis client")
	if err := r.client.Close(); err != nil {
		logs.FromContext(ctx).Warn().Err(err).Msgf("can't close redis client")
	}
}
