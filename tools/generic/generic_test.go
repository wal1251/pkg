package generic

import (
	"errors"
	"testing"
)

func TestEqualValues(t *testing.T) {
	type data struct {
		Name string
		Age  int
	}

	testCases := []struct {
		name     string
		val1     *data
		val2     *data
		expected bool
	}{
		{
			name:     "both values are nil",
			val1:     nil,
			val2:     nil,
			expected: true,
		},
		{
			name:     "one value is nil",
			val1:     &data{Name: "John", Age: 25},
			val2:     nil,
			expected: false,
		},
		{
			name:     "both values are equal",
			val1:     &data{Name: "John", Age: 25},
			val2:     &data{Name: "John", Age: 25},
			expected: true,
		},
		{
			name:     "values are not equal",
			val1:     &data{Name: "John", Age: 25},
			val2:     &data{Name: "Jane", Age: 30},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := EqualValues(tc.val1, tc.val2)
			if result != tc.expected {
				t.Errorf("Expected %v but got %v", tc.expected, result)
			}
		})
	}
}

func TestRetryWithBackoff(t *testing.T) {
	var (
		testValue  = "test"
		maxRetries = 3
	)

	tests := []struct {
		name         string
		config       ParamsRetryWithBackoff[string]
		expected     *string
		expectedErr  error
		expectedRuns int
	}{
		{
			name: "successful attempt",
			config: ParamsRetryWithBackoff[string]{
				MaxRetries: maxRetries,
				ErrMapFunc: func(err error) error { return err },
				CheckRetryNecessity: func(err error) bool {
					return true
				},
				Operation: func() (*string, error) {
					return &testValue, nil
				},
			},
			expected:     func() *string { v := testValue; return &v }(),
			expectedErr:  nil,
			expectedRuns: 1,
		},
		{
			name: "succeeds after retries",
			config: ParamsRetryWithBackoff[string]{
				MaxRetries: maxRetries,
				ErrMapFunc: func(err error) error { return err },
				CheckRetryNecessity: func(err error) bool {
					return true
				},
				Operation: func() func() (*string, error) {
					attempts := 0
					return func() (*string, error) {
						attempts++
						if attempts < maxRetries {
							return nil, errors.New("RequestRetryFailed")
						}
						return &testValue, nil
					}
				}(),
			},
			expected:     func() *string { v := testValue; return &v }(),
			expectedErr:  nil,
			expectedRuns: maxRetries,
		},
		{
			name: "non retryable error",
			config: ParamsRetryWithBackoff[string]{
				MaxRetries: maxRetries,
				ErrMapFunc: func(err error) error { return err },
				CheckRetryNecessity: func(err error) bool {
					return false
				},
				Operation: func() (*string, error) {
					return nil, errors.New("BtsInternalServerError error")
				},
			},
			expected:     nil,
			expectedErr:  errors.New("BtsInternalServerError error"),
			expectedRuns: 1,
		},
		{
			name: "mapped error stops retries",
			config: ParamsRetryWithBackoff[string]{
				MaxRetries: maxRetries,
				ErrMapFunc: func(err error) error {
					return nil
				},
				CheckRetryNecessity: func(err error) bool {
					return true
				},
				Operation: func() (*string, error) {
					return nil, errors.New("BtsCoidRequestInternalServerError")
				},
			},
			expected:     nil,
			expectedErr:  nil,
			expectedRuns: 1,
		},
		{
			name: "max retries",
			config: ParamsRetryWithBackoff[string]{
				MaxRetries: maxRetries,
				ErrMapFunc: func(err error) error { return err },
				CheckRetryNecessity: func(err error) bool {
					return true
				},
				Operation: func() (*string, error) {
					return nil, errors.New("retryable error")
				},
			},
			expected:     nil,
			expectedErr:  errors.New("retryable error"),
			expectedRuns: maxRetries,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var attempts int
			originalOperation := tt.config.Operation
			tt.config.Operation = func() (*string, error) {
				attempts++
				return originalOperation()
			}

			result, err := RetryWithBackoff(tt.config)

			if (result == nil) != (tt.expected == nil) || (result != nil && *result != *tt.expected) {
				t.Errorf("expected result %v, got %v", tt.expected, result)
			}
			if (err == nil) != (tt.expectedErr == nil) || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if attempts != tt.expectedRuns {
				t.Errorf("expected %d attempts, got %d", tt.expectedRuns, attempts)
			}
		})
	}
}
