package otp

import (
	"time"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/tools/checks"
)

const (
	CfgKeyOTPLength                cfg.Key = "OTP_LENGTH"
	CfgKeyOTPLifetime              cfg.Key = "OTP_LIFETIME"
	CfgKeyOTPNextCodeDelay         cfg.Key = "OTP_NEXT_CODE_DELAY"
	CfgKeyOTPMaxSendAttempts       cfg.Key = "OTP_MAX_SEND_ATTEMPTS"
	CfgKeyOTPSendBlockDuration     cfg.Key = "OTP_SEND_BLOCK_DURATION"
	CfgKeyOTPMaxValidateAttempts   cfg.Key = "OTP_MAX_VALIDATE_ATTEMPTS"
	CfgKeyOTPValidateBlockDuration cfg.Key = "OTP_VALIDATE_BLOCK_DURATION"

	CfgDefaultOTPLength                = 4
	CfgDefaultOTPLifetime              = time.Minute * 2
	CfgDefaultOTPNextCodeDelay         = time.Minute
	CfgDefaultOTPMaxSendAttempts       = 5
	CfgDefaultOTPSendBlockDuration     = time.Minute * 30
	CfgDefaultOTPMaxValidateAttempts   = 3
	CfgDefaultOTPValidateBlockDuration = time.Minute
)

type Config struct {
	Length                int
	Lifetime              time.Duration
	NextCodeDelay         time.Duration
	MaxSendAttempts       int
	SendBlockDuration     time.Duration
	MaxValidateAttempts   int
	ValidateBlockDuration time.Duration
}

func NewConfig(
	length int,
	lifetime time.Duration,
	nextCodeDelay time.Duration,
	maxSendAttempts int,
	sendBlockDuration time.Duration,
	maxValidateAttempts int,
	validateBlockDuration time.Duration,
) *Config {
	checks.MustBePositive(length, "Generator code length")
	checks.MustBePositive(int(lifetime.Seconds()), "OTP lifetime")
	checks.MustBePositive(int(nextCodeDelay.Seconds()), "OTP next code delay")
	checks.MustBePositive(maxSendAttempts, "OTP max send attempts")
	checks.MustBePositive(int(sendBlockDuration.Seconds()), "OTP send block duration")
	checks.MustBePositive(maxValidateAttempts, "OTP max validate attempts")
	checks.MustBePositive(int(validateBlockDuration.Seconds()), "OTP validate block duration")

	return &Config{
		Length:                length,
		Lifetime:              lifetime,
		NextCodeDelay:         nextCodeDelay,
		MaxSendAttempts:       maxSendAttempts,
		SendBlockDuration:     sendBlockDuration,
		MaxValidateAttempts:   maxValidateAttempts,
		ValidateBlockDuration: validateBlockDuration,
	}
}
