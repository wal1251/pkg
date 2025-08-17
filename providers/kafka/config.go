package kafka

import (
	"fmt"
	"os"
	"time"

	"github.com/wal1251/pkg/core/cfg"
)

const (
	ClientConfigBootstrapServers          = "bootstrap.servers" // Брокеры кластера.
	ClientConfigSecurityProtocol          = "security.protocol" // Протокол безопасности.
	ClientConfigSecurityProtocolPlainText = "PLAINTEXT"         // Без аутентификации и шифрования.

	ClientConfigSSLCALocation  = "ssl.ca.location" // Локальный путь до сертификата (SSL).
	ClientConfigSaslMechanisms = "sasl.mechanisms" // SASL: механизмы аутентификации.
	ClientConfigSaslUsername   = "sasl.username"   // SASL: имя пользователя (ACL).
	ClientConfigSaslPassword   = "sasl.password"   // SASL: пароль пользователя.

	ClientConfigID                      = "client.id"                // Идентификатор клиента.
	ClientConfigGroupID                 = "group.id"                 // Группа потребителей.
	ClientConfigAutoOffsetReset         = "auto.offset.reset"        // Сброс смещения, если нет валидного.
	ClientConfigAutoOffsetResetSmallest = "smallest"                 // Установить в начало.
	ClientConfigAutoOffsetResetLargest  = "largest"                  // Установить в конец.
	ClientConfigHeartbeatInterval       = "heartbeat.interval.ms"    // Допустимый интервал ожидания heartbeat сигнала потребителя.
	ClientConfigStatsInterval           = "statistics.interval.ms"   // Интервал сбора статистики клиента (librdkafka).
	ClientConfigEnableAutoOffsetStore   = "enable.auto.offset.store" // Автоматически сохранять смещение последнего сообщения (librdkafka).
	ClientConfigEnableAutoCommit        = "enable.auto.commit"       // Активировать авто-фиксацию смещений.
	ClientConfigAutoCommitInterval      = "auto.commit.interval.ms"  // Интервал авто-фиксации смещений.
	ClientConfigRequestRequiredACKs     = "request.required.acks"    // Сколько подтверждений ждет производитель: 0 (no), 1 (one), -1 (all).
	ClientConfigRequestRequiredACKsNo   = "0"                        // Производитель не ждет подтверждений.
	ClientConfigRequestRequiredACKsOne  = "1"                        // Производитель ждет одно подтверждение (от лидера).
	ClientConfigRequestRequiredACKsAll  = "-1"                       // Производитель ждет все подтверждения (от всех ISR).

	CfgKeyHosts            cfg.Key = "KAFKA_HOSTS"           // Брокеры кластера (string).
	CfgKeyPrefix           cfg.Key = "KAFKA_PREFIX"          // Префикс топиков и группы (string).
	CfgKeySaslUser         cfg.Key = "KAFKA_SASL_USER"       // SASL: имя пользователя (string).
	CfgKeySaslPassword     cfg.Key = "KAFKA_SASL_PASSWORD"   //nolint:gosec // SASL: пароль пользователя (string).
	CfgKeySaslMechanisms   cfg.Key = "KAFKA_SASL_MECHANISMS" // SASL: механизмы аутентификации (string).
	CfgKeySecurityProtocol cfg.Key = "KAFKA_SECURITY_PROTO"  // Протокол безопасности (string).
	CfgKeyCertCA           cfg.Key = "KAFKA_CERT_CA"         // Локальный путь до сертификата (string).
	CfgKeyTimeout          cfg.Key = "KAFKA_TIMEOUT"         // Таймаут для админа (duration).
	CfgKeyPollTimeout      cfg.Key = "KAFKA_POLL_TIMEOUT"    // Таймаут опроса брокера (duration).
	CfgKeyProducer         cfg.Key = "KAFKA_PRODUCER_CFG"    // Дополнительный параметры конфигурации клиента-производителя (map).
	CfgKeyConsumer         cfg.Key = "KAFKA_CONSUMER_CFG"    // Дополнительный параметры конфигурации клиента-потребителя (map).
	CfgKeyAdmin            cfg.Key = "KAFKA_ADMIN_CFG"       // Дополнительный параметры конфигурации клиента-админа (map).
	CfgKeyClientID         cfg.Key = "KAFKA_CLIENT_ID"       // Идентификатор клиента (string).
	CfgKeyGroupID          cfg.Key = "KAFKA_GROUP_ID"        // Группа потребителей (string).

	CfgDefaultHosts       = "localhost:9092" // Брокер по умолчанию.
	CfgDefaultGroupID     = "group.default"  // Группа потребителей по-умолчанию.
	CfgDefaultTimeout     = 30 * time.Second // Таймаут по умолчанию.
	CfgDefaultPollTimeout = 3 * time.Second  // Таймаут опроса по-умолчанию.
)

type (
	// ClientConfig каноничная конфигурация клиента KAFKA в формате ключ-значение.
	ClientConfig map[string]any

	// ClientOption опция конфигурации клиента.
	ClientOption func(ClientConfig) error

	// Config конфигурация подключения к KAFKA.
	Config struct {
		Hosts            string        // Брокеры кластера.
		SaslUser         string        // SASL: имя пользователя.
		SaslPassword     string        // SASL: пароль пользователя.
		SaslMechanisms   string        // SASL: механизмы аутентификации.
		SecurityProtocol string        // Протокол безопасности.
		CertCA           string        // Локальный путь до сертификата.
		Prefix           string        // Префикс топиков и группы.
		GroupID          string        // Группа потребителей.
		ClientID         string        // Идентификатор клиента.
		Timeout          time.Duration // Таймаут для админа.
		PollTimeout      time.Duration // Таймаут опроса брокера.
		ProducerConfig   ClientConfig  // Дополнительный параметры конфигурации клиента-производителя.
		ConsumerConfig   ClientConfig  // Дополнительный параметры конфигурации клиента-потребителя.
		AdminConfig      ClientConfig  // Дополнительный параметры конфигурации клиента-админа.
	}
)

// With возвращает конфиг клиента после применения к нему указанных опций. Вернет ошибку, если какая либо из
// опций вернет ошибку.
func (c ClientConfig) With(opts ...ClientOption) (ClientConfig, error) {
	return Options(opts...).ApplyTo(c)
}

// Append возвращает конфиг клиента после переноса в него указанной конфигурации.
func (c ClientConfig) Append(src ClientConfig) ClientConfig {
	for key, value := range src {
		c[key] = value
	}

	return c
}

// ApplyTo применяет опцию к указанному конфигу клиента.
func (o ClientOption) ApplyTo(clientCfg ClientConfig) (ClientConfig, error) {
	if o != nil {
		if err := o(clientCfg); err != nil {
			return nil, fmt.Errorf("can't apply kafka client config option: %w", err)
		}
	}

	return clientCfg, nil
}

// MakeClientConfig воздает новый экземпляр конфигурации клиента. Вернет ошибку, если не удалось создать конфигурацию
// клиента.
func MakeClientConfig(cfg *Config, opts ...ClientOption) (ClientConfig, error) {
	clientCfg := make(ClientConfig)

	return clientCfg.With(
		Options(opts...),
		Option(ClientConfigBootstrapServers, cfg.Hosts),
		Option(ClientConfigSecurityProtocol, cfg.SecurityProtocol),
		Option(ClientConfigSSLCALocation, cfg.CertCA),
		Option(ClientConfigSaslMechanisms, cfg.SaslMechanisms),
		Option(ClientConfigSaslUsername, cfg.SaslUser),
		Option(ClientConfigSaslPassword, cfg.SaslPassword),
	)
}

// ProducerConfig возвращает новый экземпляр конфигурации клиента-производителя.
func ProducerConfig(cfg *Config, opts ...ClientOption) (ClientConfig, error) {
	clientConfig, err := MakeClientConfig(cfg,
		func(config ClientConfig) error {
			config[ClientConfigRequestRequiredACKs] = ClientConfigRequestRequiredACKsNo

			return nil
		},
		Options(opts...),
	)
	if err != nil {
		return nil, err
	}

	return clientConfig.Append(cfg.ProducerConfig), nil
}

// ConsumerConfig возвращает новый экземпляр конфигурации клиента-потребителя.
func ConsumerConfig(cfg *Config, opts ...ClientOption) (ClientConfig, error) {
	clientConfig, err := MakeClientConfig(cfg,
		func(config ClientConfig) error {
			config[ClientConfigAutoOffsetReset] = ClientConfigAutoOffsetResetSmallest
			config[ClientConfigEnableAutoOffsetStore] = false
			config[ClientConfigEnableAutoCommit] = false

			return nil
		},
		Options(opts...),
		Option(ClientConfigGroupID, cfg.GroupID))
	if err != nil {
		return nil, err
	}

	return clientConfig.Append(cfg.ConsumerConfig), nil
}

// AdminConfig возвращает новый экземпляр конфигурации клиента-админа.
func AdminConfig(cfg *Config, opts ...ClientOption) (ClientConfig, error) {
	clientConfig, err := MakeClientConfig(cfg, opts...)
	if err != nil {
		return nil, err
	}

	return clientConfig.Append(cfg.AdminConfig), nil
}

// WithStats возвращает опцию конфигурации клиента, которая устанавливает периодичность генерации клиентом события
// статистики.
func WithStats(interval time.Duration) ClientOption {
	return Option(ClientConfigStatsInterval, Milliseconds(interval))
}

// WithClientID возвращает опцию конфигурации клиента, которая устанавливает идентификатор клиента. Если он не задан в
// конфигурации, тогда попытается в качестве идентификатора установить текущий хост.
func WithClientID(cfg *Config) ClientOption {
	return func(clientCfg ClientConfig) error {
		if cfg.ClientID != "" {
			clientCfg[ClientConfigID] = cfg.ClientID
		}

		hostname, err := os.Hostname()
		if err != nil {
			return fmt.Errorf("failed to get client id: %w", err)
		}

		clientCfg[ClientConfigID] = hostname

		return nil
	}
}

// Options объединяет указанные опции в единую опцию.
func Options(options ...ClientOption) ClientOption {
	return func(config ClientConfig) error {
		for _, opt := range options {
			if _, err := opt.ApplyTo(config); err != nil {
				return err
			}
		}

		return nil
	}
}

// Option возвращает новую опцию конфигурации, которая устанавливает значение для указанного ключа. Только если значение
// отличается от "нулевого".
func Option[T comparable](key string, value T) ClientOption {
	return func(clientCfg ClientConfig) error {
		ClientConfigSet(clientCfg, key, value)

		return nil
	}
}

// ClientConfigSet установит значение для указанного ключа, только если значение отличается от "нулевого".
func ClientConfigSet[T comparable](kafkaCfg ClientConfig, key string, value T) {
	var z T
	if value != z {
		kafkaCfg[key] = value
	}
}
