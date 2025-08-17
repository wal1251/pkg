package s3

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

func CfgFromViper(loader *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	return &Config{
		AWSAccessKeyID:     viperx.Get(loader, CfgKeyAWSAccessKeyID.Map(keyMapping...), CfgDefaultAWSAccessKeyID),
		AWSSecretAccessKey: viperx.Get(loader, CfgKeyAWSSecretAccessKey.Map(keyMapping...), CfgDefaultAWSSecretAccessKey),
		BucketName:         viperx.Get(loader, CfgKeyAWSBucketName.Map(keyMapping...), CfgDefaultAWSBucketName),
		ServiceURL:         viperx.Get(loader, CfgKeyAWSServiceURL.Map(keyMapping...), CfgDefaultAWSServiceURL),
		RetryCount:         viperx.Get(loader, CfgKeyAWSRetryCount.Map(keyMapping...), CfgDefaultAWSRetryCount),
		ServiceTag:         "S3",
	}
}
