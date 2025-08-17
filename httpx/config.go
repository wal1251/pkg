package httpx

import (
	"net"
	"time"

	"github.com/wal1251/pkg/core/cfg"
)

const (
	CfgKeyPort            cfg.Key = "HTTP_SERVER_PORT"
	CfgKeyTimeout         cfg.Key = "HTTP_SERVER_TIMEOUT"
	CfgKeyCacheTTL        cfg.Key = "HTTP_SERVER_CACHE_TTL"
	CfgKeyBaseURL         cfg.Key = "HTTP_SERVER_BASE_URL"
	CfgKeyMaxRate         cfg.Key = "HTTP_SERVER_MAX_RATE"
	CfgKeyRatePeriod      cfg.Key = "HTTP_SERVER_RATE_PERIOD"
	CfgKeyCORSOrigins     cfg.Key = "HTTP_SERVER_CORS_ORIGINS"
	CfgKeyCORSMethods     cfg.Key = "HTTP_SERVER_CORS_METHODS"
	CfgKeyCORSHeaders     cfg.Key = "HTTP_SERVER_CORS_HEADERS"
	CfgKeyCORSCredentials cfg.Key = "HTTP_SERVER_CORS_CREDENTIALS" //nolint:gosec
	CfgKeyFrontMinVersion cfg.Key = "HTTP_SERVER_FRONT_MIN_VERSION"

	CfgDefaultPort              = "8080"
	CfgDefaultReadHeaderTimeout = 30 * time.Second
	CfgDefaultTimeout           = 3000 * time.Millisecond
	CfgDefaultCacheTTL          = 300 * time.Second
	CfgDefaultBaseURL           = "/api/"
	CfgDefaultMaxRate           = 10
	CfgDefaultRatePeriod        = 1 * time.Minute
	CfgDefaultCORSOrigins       = "*"
	CfgDefaultCORSMethod        = "GET,POST,PUT,DELETE,OPTIONS"
	CfgDefaultCORSHeaders       = "*"
)

type Config struct {
	Port              string
	ReadHeaderTimeout time.Duration
	Timeout           time.Duration
	CacheTTL          time.Duration
	BaseURL           string
	Rate              RateLimit
	CORS              CORSConfig
	FrontMinVersion   string
}

type RateLimit struct {
	MaxRate int
	Every   time.Duration
}

type CORSConfig struct {
	Origins     string
	Methods     string
	Headers     string
	Credentials bool
}

// Addr returns server address in format ":<port>".
func (c Config) Addr() string {
	return net.JoinHostPort("", c.Port)
}
