package clock_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/tools/clock"
)

func TestRetryingWrapper(t *testing.T) {
	var retries int
	retriesInc := func(error) bool {
		retries++
		return true
	}

	tests := []struct {
		name            string
		funcToRetry     func() error
		retriesCount    int
		expectedErr     bool
		expectedRetries int
	}{
		{
			name: "Wrapped function returns nil",
			funcToRetry: func() error {
				return nil
			},
			retriesCount:    3,
			expectedErr:     false,
			expectedRetries: 0,
		},
		{
			name: "Wrapped function always returns error",
			funcToRetry: func() error {
				return errors.New("error")
			},
			retriesCount:    3,
			expectedErr:     true,
			expectedRetries: 4,
		},
		{
			name: "Wrapped function returns 3 times error and then nil",
			funcToRetry: func() error {
				if retries < 3 {
					return errors.New("error")
				}
				return nil
			},
			retriesCount:    5,
			expectedErr:     false,
			expectedRetries: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			retries = 0
			retryFunc := clock.RetryingWrapper(tc.retriesCount, 100*time.Millisecond)
			err := retryFunc(tc.funcToRetry, retriesInc)

			if tc.expectedErr {
				assert.NotNil(t, err, "Expected error")
			} else {
				assert.Nil(t, err, "Expected no error")
			}

			assert.Equal(t, tc.expectedRetries, retries)
		})
	}
}

func TestRetryingWrapperWE(t *testing.T) {
	// Test successful retry
	successfulFunc := func() (bool, error) {
		return true, nil
	}
	retryFunc := clock.RetryingWrapperWE(3, 100*time.Millisecond)
	err := retryFunc(successfulFunc)
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}

	// Test unsuccessful retry
	unsuccessfulFunc := func() (bool, error) {
		return true, errors.New("error")
	}
	retryFunc = clock.RetryingWrapperWE(3, 100*time.Millisecond)
	err = retryFunc(unsuccessfulFunc)
	if err == nil {
		t.Errorf("Expected error to be non-nil, but got nil")
	}
}
