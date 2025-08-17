package grpcx

import (
	"context"
	"errors"
	"testing"

	"github.com/wal1251/pkg/core/errs"
	"github.com/wal1251/pkg/tools/acceptlanguage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestUserInfoClientInterceptor(t *testing.T) {
	// Create a context with user information
	userID := "12345"
	phoneNumber := "+77778987788"
	name := "John"
	sessionID := "session_123123"
	// package rmr-pkg/httpx/mw.authorizer.go  Authorizer way how we put user info to context
	userInfo := map[string]interface{}{
		"user_id":      userID,
		"phone_number": phoneNumber,
		"name":         name,
		"session_id":   sessionID,
	}
	ctx := context.WithValue(context.Background(), UserInfoKey, userInfo)

	// Create a dummy invoker
	invoker := func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		opts ...grpc.CallOption,
	) error {
		// Extract metadata from the context
		md, ok := metadata.FromOutgoingContext(ctx)
		assert.True(t, ok, "metadata should be present in the context")
		assert.Equal(t, userID, md["user_id"][0], "user_id should be present in metadata")
		assert.Equal(t, phoneNumber, md["phone_number"][0], "phone_number should be present in metadata")
		assert.Equal(t, name, md["name"][0], "name should be present in metadata")
		assert.Equal(t, sessionID, md["session_id"][0], "session_id should be present in metadata")
		return nil
	}

	// Create the interceptor
	interceptor := UserInfoClientInterceptor()

	// Call the interceptor
	err := interceptor(ctx, "/test/method", nil, nil, nil, invoker)
	assert.NoError(t, err, "interceptor should not return an error")
}

func TestUserInfoServerInterceptor(t *testing.T) {
	// Setup test data
	userID := "12345"
	phoneNumber := "+77778987788"
	name := "John"
	sessionID := "session_123123"

	// Create metadata that simulates incoming context
	md := metadata.New(
		map[string]string{
			"user_id":       userID,
			"phone_number":  phoneNumber,
			"name":          name,
			"session_id":    sessionID,
			"grpc-internal": "should-be-ignored", // This should be filtered out
			"somehting":     "should-be-ignored", // This should be filtered out
		},
	)

	// Create incoming context with metadata
	ctx := metadata.NewIncomingContext(context.Background(), md)

	// Create a mock handler that will verify the context
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		// Extract user info from context
		userInfo, ok := ctx.Value(UserInfoKey).(map[string]interface{})
		assert.True(t, ok, "user info should be present in context")

		// Verify all expected values are present
		assert.Equal(t, userID, userInfo["user_id"], "user_id should match")
		assert.Equal(t, phoneNumber, userInfo["phone_number"], "phone_number should match")
		assert.Equal(t, name, userInfo["name"], "name should match")
		assert.Equal(t, sessionID, userInfo["session_id"], "session_id should match")

		// Verify grpc- prefixed keys are not present
		_, hasGrpcKey := userInfo["grpc-internal"]
		assert.False(t, hasGrpcKey, "grpc- prefixed keys should be filtered out")
		_, hasSomethingKey := userInfo["something"]
		assert.False(t, hasSomethingKey, "somehting should be filtered out")
		// also check that metadata is still present
		md, ok := metadata.FromIncomingContext(ctx)
		assert.True(t, ok, "metadata should be present in the context")
		assert.Equal(t, userID, md["user_id"][0], "user_id should be present in metadata")
		assert.Equal(t, phoneNumber, md["phone_number"][0], "phone_number should be present in metadata")
		assert.Equal(t, name, md["name"][0], "name should be present in metadata")
		assert.Equal(t, sessionID, md["session_id"][0], "session_id should be present in metadata")

		return nil, nil
	}

	// Create and call the interceptor
	interceptor := UserInfoServerInterceptor()
	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)

	assert.NoError(t, err, "interceptor should not return an error")
}

func TestSecureErrorInterceptor(t *testing.T) {
	serviceID := 2 // ID сервиса для отладки.

	t.Run("success case - no error", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "success", nil
		}

		interceptor := SecureErrorInterceptor(serviceID)
		resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

		assert.NoError(t, err)
		assert.Equal(t, "success", resp)
	})

	t.Run("error case - with status error", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, status.Error(codes.Internal, "internal error")
		}

		interceptor := SecureErrorInterceptor(serviceID)
		resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

		assert.Error(t, err)
		assert.Nil(t, resp)

		// Verify the error was converted to a system failure error
		st, ok := status.FromError(err)
		assert.True(t, ok)

		details := st.Details()
		assert.Len(t, details, 1)

		info, ok := details[0].(*errdetails.ErrorInfo)
		assert.True(t, ok)
		assert.Equal(t, "2.0", info.Metadata["error_num"])
	})

	t.Run("error case - default error", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, errors.New("def error")
		}

		interceptor := SecureErrorInterceptor(serviceID)
		resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

		assert.Error(t, err)
		assert.Nil(t, resp)

		// Verify the error was converted to a system failure error
		st, ok := status.FromError(err)
		assert.True(t, ok)

		details := st.Details()
		assert.Len(t, details, 1)

		info, ok := details[0].(*errdetails.ErrorInfo)
		assert.True(t, ok)
		assert.Equal(t, "0", info.Metadata["error_num"])
	})

	t.Run("error case - with custom error", func(t *testing.T) {
		detail := make(map[string]string)
		detail["fund"] = "not enough money"
		customErr := errs.Reasons("custom error", errs.TypeIllegalArgument, "500").WithDetails(detail)
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, customErr
		}

		interceptor := SecureErrorInterceptor(serviceID)
		resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

		assert.Error(t, err)
		assert.Nil(t, resp)

		err = ErrorFromGRPC(err)
		var typedErr errs.Error
		if !errors.As(err, &typedErr) {
			t.Errorf("ErrorFromGRPC failed to convert error to errs.Error")
		}
		assert.NotNil(t, typedErr.Details)
		assert.Equal(t, detail["fund"], typedErr.Details["fund"])
	})
}

func TestAcceptLanguageClientInterceptor(t *testing.T) {
	ctx := context.WithValue(context.Background(), acceptlanguage.AcceptLanguageKey, "ru")

	// Create a dummy invoker
	invoker := func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		opts ...grpc.CallOption,
	) error {
		// Extract metadata from the context
		md, ok := metadata.FromOutgoingContext(ctx)
		assert.True(t, ok, "metadata should be present in the context")
		assert.Equal(t, "ru", md["accept_language"][0], "accept_language should be present in metadata")
		return nil
	}

	// Create the interceptor
	interceptor := AcceptLanguageClientInterceptor()

	err := interceptor(ctx, "/test/method", nil, nil, nil, invoker)
	assert.NoError(t, err, "interceptor should not return an error")
}
