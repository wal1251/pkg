package grpcx

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/wal1251/pkg/core/logs"
)

func StartServer(ctx context.Context, cfg *Config, server *grpc.Server) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	log := logs.FromContext(ctx)
	log.Info().Msg("starting grpc server")

	// Включаем reflection, если это указано в конфигурации
	if cfg.Reflection {
		log.Info().Msg("grpc reflection enabled")
		reflection.Register(server)
	}

	serverErrors := make(chan error, 1)

	// Начинаем слушать gRPC запросы
	grpcListener, err := net.Listen("tcp", cfg.AddrServer())
	if err != nil {
		log.Error().Msg("failed to listen on grpc port")

		return fmt.Errorf("failed to listen on grpc port: %w", err)
	}

	go func() {
		serverErrors <- server.Serve(grpcListener)
	}()
	log.Info().Msgf("grpc service started at address: %s", cfg.AddrServer())

	select {
	case err := <-serverErrors:
		log.Info().Msg("grpc server has closed")

		return fmt.Errorf("grpc server has closed: %w", err)
	case sig := <-shutdown:
		log.Info().Str("signal", sig.String()).Msg("Start shutdown")
		server.GracefulStop()
	case <-ctx.Done():
		log.Info().Msg("closing grpc server due to context cancellation")
		server.GracefulStop()

		return nil
	}

	return nil
}

func SetOptions(cfg *Config) ([]grpc.ServerOption, error) {
	var serverOptions []grpc.ServerOption

	// Устанавливаем ограничения на размер сообщения
	serverOptions = append(
		serverOptions,
		grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.MaxSendMsgSize),
		grpc.ChainUnaryInterceptor(UserInfoServerInterceptor()),
	)

	// Настраиваем параметры keepalive
	serverOptions = append(
		serverOptions,
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionIdle:     cfg.MaxConnectionIdle,
				MaxConnectionAge:      cfg.MaxConnectionAge,
				MaxConnectionAgeGrace: cfg.MaxConnectionAgeGrace,
				Time:                  cfg.KeepAliveTime,
				Timeout:               cfg.KeepAliveTimeout,
			},
		),
	)

	// Настраиваем TLS, если предоставлены сертификаты
	if cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertFile, cfg.TLSKeyFile)
		if err != nil {
			log.Error().Msg("failed to load TLS credentials")

			return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
		}
		serverOptions = append(serverOptions, grpc.Creds(creds))
	}

	return serverOptions, nil
}
