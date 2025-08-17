package kafka

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

// CfgFromViper загружает и возвращает конфигурацию подключения к KAFKA с помощью viper.
func CfgFromViper(loader *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	prefix := cfgFromViperPrefix(loader, keyMapping...)
	groupID := WithPrefix(prefix).Map(viperx.Get(loader, CfgKeyGroupID.Map(keyMapping...), CfgDefaultGroupID))

	return &Config{
		Hosts:            viperx.Get(loader, CfgKeyHosts.Map(keyMapping...), CfgDefaultHosts),
		SaslUser:         viperx.Get(loader, CfgKeySaslUser.Map(keyMapping...), ""),
		SaslPassword:     viperx.Get(loader, CfgKeySaslPassword.Map(keyMapping...), ""),
		SaslMechanisms:   viperx.Get(loader, CfgKeySaslMechanisms.Map(keyMapping...), ""),
		SecurityProtocol: viperx.Get(loader, CfgKeySecurityProtocol.Map(keyMapping...), ""),
		CertCA:           viperx.Get(loader, CfgKeyCertCA.Map(keyMapping...), ""),
		Timeout:          viperx.Get(loader, CfgKeyTimeout.Map(keyMapping...), CfgDefaultTimeout),
		PollTimeout:      viperx.Get(loader, CfgKeyPollTimeout.Map(keyMapping...), CfgDefaultPollTimeout),
		ClientID:         viperx.Get(loader, CfgKeyClientID.Map(keyMapping...), ""),
		ProducerConfig:   loader.GetStringMap(CfgKeyProducer.Map(keyMapping...).String()),
		ConsumerConfig:   loader.GetStringMap(CfgKeyConsumer.Map(keyMapping...).String()),
		AdminConfig:      loader.GetStringMap(CfgKeyAdmin.Map(keyMapping...).String()),
		Prefix:           prefix,
		GroupID:          groupID,
	}
}

func cfgFromViperPrefix(loader *viper.Viper, keyMapping ...cfg.KeyMap) string {
	return Prefix(viperx.Get(loader, CfgKeyPrefix.Map(keyMapping...), ""), viperx.Environment(loader))
}
