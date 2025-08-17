package grpcx

import (
	"net"
	"time"

	"github.com/wal1251/pkg/core/cfg"
)

const (
	CfgKeyHost                  cfg.Key = "GRPC_SERVER_HOST"
	CfgKeyPort                  cfg.Key = "GRPC_SERVER_PORT"
	CfgKeyReflection            cfg.Key = "GRPC_SERVER_REFLECTION"
	CfgKeyMaxConnectionIdle     cfg.Key = "GRPC_SERVER_MAX_CONN_IDLE"
	CfgKeyMaxConnectionAge      cfg.Key = "GRPC_SERVER_MAX_CONN_AGE"
	CfgKeyMaxConnectionAgeGrace cfg.Key = "GRPC_SERVER_MAX_CONN_AGE_GRACE"
	CfgKeyKeepAliveTime         cfg.Key = "GRPC_SERVER_KEEP_ALIVE_TIME"
	CfgKeyKeepAliveTimeout      cfg.Key = "GRPC_SERVER_KEEP_ALIVE_TIMEOUT"
	CfgKeyMaxRecvMsgSize        cfg.Key = "GRPC_SERVER_MAX_RECV_MSG_SIZE"
	CfgKeyMaxSendMsgSize        cfg.Key = "GRPC_SERVER_MAX_SEND_MSG_SIZE"
	CfgKeyTLSCertFile           cfg.Key = "GRPC_SERVER_TLS_CERT_FILE"
	CfgKeyTLSKeyFile            cfg.Key = "GRPC_SERVER_TLS_KEY_FILE"

	CfgDefaultHost                  = "localhost"
	CfgDefaultHostKube              = "0.0.0.0"
	CfgDefaultPort                  = "50051"
	CfgDefaultReflection            = false
	CfgDefaultMaxConnectionIdle     = 5 * time.Minute
	CfgDefaultMaxConnectionAge      = 2 * time.Hour
	CfgDefaultMaxConnectionAgeGrace = 5 * time.Minute
	CfgDefaultKeepAliveTime         = 2 * time.Hour
	CfgDefaultKeepAliveTimeout      = 20 * time.Second
	CfgDefaultMaxRecvMsgSize        = 1024 * 1024 * 20
	CfgDefaultMaxSendMsgSize        = 1024 * 1024 * 20
	CfgDefaultTLSCertFile           = ""
	CfgDefaultTLSKeyFile            = ""
)

type Config struct {
	Host                  string
	Port                  string
	Reflection            bool
	MaxConnectionIdle     time.Duration
	MaxConnectionAge      time.Duration
	MaxConnectionAgeGrace time.Duration
	KeepAliveTime         time.Duration
	KeepAliveTimeout      time.Duration
	MaxRecvMsgSize        int
	MaxSendMsgSize        int
	TLSCertFile           string
	TLSKeyFile            string
	ServiceID             int
}

// Addr returns server address in format ":<port>".
func (c Config) Addr() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c Config) AddrServer() string {
	return net.JoinHostPort(CfgDefaultHostKube, c.Port)
}
