package otelx

import "github.com/wal1251/pkg/core/cfg"

const (
	CfgKeyExporterEndpoint  cfg.Key = "OTEL_EXPORTER_ENDPOINT"   // Эндпоинт для отправки данных.
	CfgKeyServiceName       cfg.Key = "OTEL_SERVICE_NAME"        // Название сервиса.
	CfgKeySamplerTraceRatio cfg.Key = "OTEL_SAMPLER_TRACE_RATIO" // Процент трейсов, которые нужно отправлять.

	CfgDefaultExporterEndpoint  = "http://localhost:4318"
	CfgDefaultServiceName       = "service"
	CfgDefaultSamplerTraceRatio = 1.0
)

type Config struct {
	ExporterEndpoint  string
	ServiceName       string
	SamplerTraceRatio float64
	Env               string
}
