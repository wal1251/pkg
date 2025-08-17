package confluent

import (
	"context"
	"fmt"
	"strings"
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/providers/kafka"
	"github.com/wal1251/pkg/tools/collections"
)

type (
	// Admin клиент-админ KAFKA.
	Admin struct {
		client  *confluent.AdminClient
		prefix  string
		timeout time.Duration
	}
)

// TopicMetadata возвращает информацию о топиках.
func (a *Admin) TopicMetadata() (*confluent.Metadata, error) {
	metadata, err := a.client.GetMetadata(nil, true, int(a.timeout.Milliseconds()))
	if err != nil {
		return nil, fmt.Errorf("failed to get all topics metadata: %w", err)
	}

	topics := make(map[string]confluent.TopicMetadata)
	for name, topic := range metadata.Topics {
		if !strings.HasPrefix(name, a.prefix) {
			continue
		}
		topics[name] = topic
	}

	metadata.Topics = topics

	return metadata, nil
}

// TopicCreateIfNotExists создать топики, если они еще не существуют. Если при создании одного из топиков возникла
// ошибка, то вернут ошибку kafka.ErrTopicCreateFailed.
func (a *Admin) TopicCreateIfNotExists(ctx context.Context, topics ...kafka.TopicMetadata) (bool, error) {
	logger := logs.FromContext(ctx)

	metadata, err := a.client.GetMetadata(nil, true, int(a.timeout.Milliseconds()))
	if err != nil {
		return false, fmt.Errorf("failed to retrieve metadata: %w", err)
	}

	toCreate := collections.Filter(topics, func(topicMetadata kafka.TopicMetadata) bool {
		_, ok := metadata.Topics[topicMetadata.Name]

		return !ok
	})

	if len(toCreate) == 0 {
		logger.Debug().Msg("all topics exist")

		return false, nil
	}

	results, err := a.client.CreateTopics(ctx, collections.Map(toCreate, MakeTopicSpec))
	if err != nil {
		return false, fmt.Errorf("failed to create topics: %w", err)
	}

	hasErrors := false
	for _, result := range results {
		if result.Error.Code() != 0 {
			hasErrors = true

			logger.Err(result.Error).Msgf("failed to create topic: %s", result.Topic)

			continue
		}

		logger.Info().Msgf("topic crated: %s", result.Topic)
	}

	if hasErrors {
		return false, kafka.ErrTopicCreateFailed
	}

	return true, nil
}

// TopicDelete удаляет топики с указанными именами. Если при удалении одного из топиков произошла ошибка, то вернет
// ошибку kafka.ErrTopicDeleteFailed.
func (a *Admin) TopicDelete(ctx context.Context, names ...string) error {
	results, err := a.client.DeleteTopics(ctx, names)
	if err != nil {
		return fmt.Errorf("failed to delete topics %v: %w", names, err)
	}

	logger := logs.FromContext(ctx)

	hasErrors := false
	for _, result := range results {
		if result.Error.Code() != 0 {
			logger.Err(result.Error).Msgf("failed to delete topic: %s", result.Topic)

			hasErrors = true

			continue
		}

		logger.Info().Msgf("topic deleted: %s", result.Topic)
	}

	if hasErrors {
		return kafka.ErrTopicDeleteFailed
	}

	return nil
}

// Client возвращает более низкоуровневый клиент.
func (a *Admin) Client() *confluent.AdminClient {
	return a.client
}

// Close закрыть соединение с брокером.
func (a *Admin) Close() {
	a.client.Close()
}

// NewAdmin возвращает новый экземпляр клиента-админа Admin.
func NewAdmin(cfg *kafka.Config, opts ...kafka.ClientOption) (*Admin, error) {
	clientConfig, err := kafka.AdminConfig(cfg, opts...)
	if err != nil {
		return nil, fmt.Errorf("can't create kafka admin: %w", err)
	}

	adm, err := confluent.NewAdminClient(ConfigMap(clientConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka admin client: %w", err)
	}

	return &Admin{
		client:  adm,
		prefix:  cfg.Prefix,
		timeout: cfg.Timeout,
	}, nil
}
