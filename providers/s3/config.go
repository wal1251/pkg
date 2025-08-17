package s3

import "github.com/wal1251/pkg/core/cfg"

const (
	SigningRegion = "ru-central1"
	PartitionID   = "yc"

	DefaultMaxKeys = 1000

	CfgDefaultAWSServiceURL      = ""
	CfgDefaultAWSBucketName      = ""
	CfgDefaultAWSAccessKeyID     = ""
	CfgDefaultAWSSecretAccessKey = ""
	CfgDefaultAWSRetryCount      = 3

	CfgKeyAWSServiceURL      cfg.Key = "AWS_SERVICE_URL"       // Точка входа для подключения к AWS сервису
	CfgKeyAWSBucketName      cfg.Key = "AWS_BUCKET_NAME"       // Имя контейнера(хранилища) объектов
	CfgKeyAWSAccessKeyID     cfg.Key = "AWS_ACCESS_KEY_ID"     // Ключ доступа
	CfgKeyAWSSecretAccessKey cfg.Key = "AWS_SECRET_ACCESS_KEY" // nolint Секретный ключ доступа
	CfgKeyAWSRetryCount      cfg.Key = "AWS_RETRY_COUNT"       // Число повторных попыток
)

type (
	Config struct {
		AWSAccessKeyID     string
		AWSSecretAccessKey string
		ServiceURL         string
		RetryCount         int
		ServiceTag         string
		BucketName         string
	}
)
