package generator

import (
	"strings"
)

type TestOTPDouble struct {
	length int
}

func NewTestDouble(length int) *TestOTPDouble {
	return &TestOTPDouble{length: length}
}

func (o *TestOTPDouble) Generate() (string, error) {
	otp := strings.Repeat("0", o.length)

	return otp, nil
}
