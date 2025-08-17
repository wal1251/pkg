package otp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wal1251/pkg/core/memorystore"
	"github.com/wal1251/pkg/providers/otp/generator"
)

const (
	memoryStoreKeyPrefixOTP              = "otp:"
	memoryStorePrefixOTPValidateAttempts = "otp_validate_attempts:"
)

var (
	ErrTooManyAttempts     = errors.New("too many attempts")
	ErrTooFrequentAttempts = errors.New("too frequent attempts")
	ErrWrongCode           = errors.New("wrong code")
)

// otp представляет собой структуру для хранения информации о сгенерированном OTP.
type otp struct {
	Code         string    `json:"code"`
	ValidUntil   time.Time `json:"validUntil"`
	SendAttempts int       `json:"sendAttempts"`
	LastSendTime time.Time `json:"lastSendTime"`
}

// ManagerInterface определяет интерфейс для работы с OTP.
type ManagerInterface interface {
	Send(ctx context.Context, target string, msgTemplate string) (*time.Duration, error)
	Validate(ctx context.Context, target string, code string) error
}

// sender - интерфейс для отправителя OTP.
type sender interface {
	Send(ctx context.Context, to string, code string) error
}

// Manager представляет собой реализацию интерфейса ManagerInterface для работы с OTP.
type Manager struct {
	sender      sender
	generator   generator.Generator
	memoryStore memorystore.MemoryStore

	config *Config
}

// Send генерирует одноразовый пароль (OTP), отправляет его целевому адресату и возвращает
// задержку перед следующей попыткой отправки OTP. Если достигнуто максимальное количество
// попыток, функция вернет ошибку ErrTooManyAttempts.
//
// Входные параметры:
//   - ctx: Контекст выполнения операции.
//   - target: Целевой адресат OTP.
//   - msgTemplate: Шаблон сообщения, которое будет отправлено. В формате - "Ваш код: %s".
//     Если нужно отправить лишь код, то следует передать пустую строку.
//
// Возвращаемые значения:
//   - *time.Duration: Задержка перед следующей попыткой отправки OTP.
//   - error: Ошибка, если что-то пошло не так, например, при генерации OTP, сохранении
//     его в хранилище данных в памяти или отправке OTP.
//
// Пример использования:
//
//	delay, err := manager.Send(ctx, "+79123456789")
//	if err != nil {
//	    log.Error("Failed to send OTP:", err)
//	} else {
//	    log.Info("OTP sent successfully. Next attempt in:", delay)
//	}
func (m *Manager) Send(ctx context.Context, target string, msgTemplate string) (*time.Duration, error) {
	key := memoryStoreKeyPrefixOTP + target
	// Попытка получить существующий OTP из хранилища
	rawOTP, err := m.memoryStore.Get(ctx, key)
	if err != nil && !errors.Is(err, memorystore.ErrKeyNotFound) {
		// Обработка ошибок, не связанных с отсутствием ключа
		return nil, fmt.Errorf("failed to get OTP from memory store: %w", err)
	}

	var OTP otp
	if !errors.Is(err, memorystore.ErrKeyNotFound) {
		// Десериализация существующего OTP
		err = rawOTP.Struct(&OTP)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal OTP: %w", err)
		}

		// Проверка на слишком частые попытки отправки
		if OTP.LastSendTime.Add(m.config.NextCodeDelay).After(time.Now()) {
			return nil, ErrTooFrequentAttempts
		}
	}

	// Проверка на превышение максимального количества попыток отправки
	if OTP.SendAttempts >= m.config.MaxSendAttempts {
		return nil, ErrTooManyAttempts
	}

	// Генерация нового OTP
	newCode, err := m.generator.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Обновление данных OTP
	OTP.Code = newCode
	OTP.ValidUntil = time.Now().Add(m.config.Lifetime)
	OTP.SendAttempts++
	OTP.LastSendTime = time.Now()

	// Установка времени блокировки, если достигнуто максимальное количество попыток
	expiration := 0
	if OTP.SendAttempts == m.config.MaxSendAttempts {
		expiration = int(m.config.SendBlockDuration.Seconds())
	}

	if err := m.memoryStore.Set(ctx, key, OTP, time.Duration(expiration)*time.Second); err != nil {
		return nil, fmt.Errorf("failed to set OTP in memory store: %w", err)
	}

	// Обнуление количества попыток валидации OTP
	otpValidateKey := memoryStorePrefixOTPValidateAttempts + target
	if err := m.memoryStore.Set(ctx, otpValidateKey, 0, m.config.Lifetime); err != nil {
		return nil, fmt.Errorf("failed to set OTP validate count in memory store: %w", err)
	}

	// Формирование и отправка сообщения с OTP
	msg := OTP.Code
	if msgTemplate != "" {
		msg = fmt.Sprintf(msgTemplate, OTP.Code)
	}
	if err := m.sender.Send(ctx, target, msg); err != nil {
		return nil, fmt.Errorf("failed to send OTP: %w", err)
	}

	// Возврат задержки перед следующей попыткой отправки
	return &m.config.NextCodeDelay, nil
}

// Validate проверяет введенный одноразовый пароль (OTP) на соответствие сохраненному OTP
// для указанного целевого адресата. Если OTP совпадает и его срок действия не истек,
// функция удаляет OTP из хранилища и возвращает nil.
// В противном случае, возвращает ошибку ErrWrongCode.
// Если достигнуто максимальное количество попыток проверки, функция вернет ошибку ErrTooManyAttempts.
//
// Входные параметры:
//   - ctx: Контекст выполнения операции.
//   - target: Целевой адресат OTP.
//   - code: Введенный пользователем одноразовый пароль для проверки.
//
// Возвращаемые значения:
//   - error: Ошибка, если OTP неверный или его срок действия истек, или если произошла
//     ошибка при удалении OTP из хранилища.
//
// Пример использования:
//
//	err := manager.Validate(ctx, "+79123456789", "123456")
//	if err != nil {
//	    log.Error("OTP validation failed:", err)
//	} else {
//	    log.Info("OTP validated successfully.")
//	}
func (m *Manager) Validate(ctx context.Context, target string, code string) error {
	if err := m.checkOTPValidateAttemptsCount(ctx, target); err != nil {
		return err
	}

	otpKey := memoryStoreKeyPrefixOTP + target
	rawValidOTP, err := m.memoryStore.Get(ctx, otpKey)
	if err != nil {
		if errors.Is(err, memorystore.ErrKeyNotFound) {
			return ErrWrongCode
		}

		return fmt.Errorf("failed to get OTP from memory store: %w", err)
	}

	var validOTP otp
	if err := rawValidOTP.Struct(&validOTP); err != nil {
		return fmt.Errorf("failed to unmarshal OTP: %w", err)
	}

	if validOTP.Code != code || validOTP.ValidUntil.Before(time.Now()) {
		return ErrWrongCode
	}

	_, err = m.memoryStore.Delete(ctx, otpKey)
	if err != nil {
		return fmt.Errorf("failed to delete OTP from memory store: %w", err)
	}

	return nil
}

// NewNumericOTPManager создает новый экземпляр Manager с генератором числовых OTP.
func NewNumericOTPManager(sender sender, memoryStore memorystore.MemoryStore, config *Config) *Manager {
	otpGenerator := generator.NewNumericOTPGenerator(config.Length)

	return NewManager(sender, memoryStore, config, otpGenerator)
}

// NewManager создает новый экземпляр Manager.
func NewManager(sender sender, memoryStore memorystore.MemoryStore, config *Config, generator generator.Generator) *Manager {
	return &Manager{
		generator:   generator,
		sender:      sender,
		memoryStore: memoryStore,
		config:      config,
	}
}

func (m *Manager) checkOTPValidateAttemptsCount(ctx context.Context, target string) error {
	otpValidateAttemptsKey := memoryStorePrefixOTPValidateAttempts + target

	rawOTPValidateCount, err := m.memoryStore.Get(ctx, otpValidateAttemptsKey)
	if err != nil && !errors.Is(err, memorystore.ErrKeyNotFound) {
		return fmt.Errorf("failed to get OTP validate count from memory store: %w", err)
	}

	var otpValidateCount int
	if !errors.Is(err, memorystore.ErrKeyNotFound) {
		otpValidateCount, err = rawOTPValidateCount.Int()
		if err != nil {
			return fmt.Errorf("failed to parse OTP validate count: %w", err)
		}
	}

	if otpValidateCount >= m.config.MaxValidateAttempts {
		return ErrTooManyAttempts
	}
	otpValidateCount++

	expiration := m.config.Lifetime
	if otpValidateCount == m.config.MaxValidateAttempts {
		expiration = m.config.ValidateBlockDuration
	}

	err = m.memoryStore.Set(ctx, otpValidateAttemptsKey, otpValidateCount, expiration)
	if err != nil {
		return fmt.Errorf("failed to set OTP validate count in memory store: %w", err)
	}

	return nil
}
