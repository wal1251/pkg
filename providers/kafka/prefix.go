package kafka

import (
	"fmt"
	"strings"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/cfg"
)

func Prefix(prefix string, env cfg.Environment) string {
	if env == "" {
		return prefix
	}

	if prefix == "" {
		return string(env)
	}

	return fmt.Sprintf("%s.%s", prefix, env)
}

func WithPrefix(prefix string) core.Map[string, string] {
	return func(topic string) string {
		if prefix == "" || topic == "" {
			return topic
		}

		return fmt.Sprintf("%s.%s", prefix, topic)
	}
}

func WithoutPrefix(prefix string) core.Map[string, string] {
	return func(topic string) string {
		return strings.TrimPrefix(topic, prefix+".")
	}
}
