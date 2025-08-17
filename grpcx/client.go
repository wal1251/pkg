package grpcx

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const maxCallRecvMsgSize = 1024 * 1024 * 20

func ConnectServer(cfg *Config) (*grpc.ClientConn, error) {
	// Инициализация списка опций для подключения gRPC
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize),
			grpc.MaxCallSendMsgSize(maxCallRecvMsgSize),
		),
		grpc.WithBlock(), // Блокировать до установления соединения
		grpc.WithUnaryInterceptor(UserInfoClientInterceptor()),
	}

	// Настройка TLS, если предоставлены сертификаты
	if cfg.TLSCertFile != "" {
		creds, err := credentials.NewClientTLSFromFile(cfg.TLSCertFile, "")
		if err != nil {
			log.Error().Msgf("не удалось загрузить TLS сертификаты: %v", err)

			return nil, fmt.Errorf("не удалось загрузить TLS сертификаты: %w", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		// Если TLS не настроен, подключаемся небезопасно
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Настройка параметров keepalive
	keepaliveParams := keepalive.ClientParameters{
		Time:                cfg.KeepAliveTime,
		Timeout:             cfg.KeepAliveTimeout,
		PermitWithoutStream: true, // Отправлять пинги даже без активных потоков
	}
	opts = append(opts, grpc.WithKeepaliveParams(keepaliveParams))

	// Установка таймаута для контекста
	ctx, cancel := context.WithTimeout(context.Background(), cfg.KeepAliveTimeout)
	defer cancel()

	// Подключение к gRPC серверу
	conn, err := grpc.DialContext(ctx, cfg.Addr(), opts...)
	if err != nil {
		log.Error().Msgf("failed to connect to GRPC server at address: %s, error: %v", cfg.Addr(), err)

		return nil, fmt.Errorf("failed to connect to GRPC server at address %s: %w", cfg.Addr(), err)
	}

	log.Info().Msgf("successfully connected to GRPC server at address: %s", cfg.Addr())

	return conn, nil
}
