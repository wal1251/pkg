package grpcx

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wal1251/pkg/core/errs"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

const (
	// StatusErrorDefault используется как код ошибки по умолчанию, если тип ошибки не определен
	// в маппинге. Соответствует codes.InvalidArgument, который сигнализирует о некорректных входных данных.
	StatusErrorDefault = codes.InvalidArgument

	// StatusErrorSystemFailure код ошибки, используемый для обще-системных сбоев (например, проблемы сервера).
	// Соответствует codes.Internal, который сигнализирует о внутренних ошибках сервера.
	StatusErrorSystemFailure = codes.Internal

	DetailedErrorPrefix = "DETAILED_ERROR_" // Префикс для дополнительных полей ошибки.
)

// NewGrpcError преобразует ошибку в gRPC статус, если ошибка относится к типу errs.Error.
// Если тип ошибки определен в маппинге defaultErrorToGrpcStatusMapper, возвращается соответствующий gRPC статус.
// Если тип ошибки не определен, возвращается код по умолчанию (StatusErrorDefault).
// Если ошибка не является типом errs.Error, возвращается системная ошибка (StatusErrorSystemFailure).
//
// Пример использования:
//
//	func (s *Server) ExampleMethod(ctx context.Context) (*pb.ExampleResponse, error) {
//		result, err := s.uc.SomeMethod(ctx, req)
//		if err != nil {
//			return nil, grpcx.NewGrpcError(err)
//		}
//		return &pb.ExampleResponse{}, nil
//	}
func NewGrpcError(err error) error {
	if err == nil {
		return nil
	}

	mapper := getGrpcErrorMapper()

	var typedErr errs.Error
	if errors.As(err, &typedErr) {
		if mappedErr, exists := mapper[typedErr.Type]; exists {
			return status.Error(mappedErr, typedErr.Error()) //nolint:wrapcheck // already should be wrapped
		}

		return status.Error(StatusErrorDefault, typedErr.Error()) //nolint:wrapcheck // already should be wrapped
	}

	return status.Error(StatusErrorSystemFailure, err.Error()) //nolint:wrapcheck // already should be wrapped
}

func NewDetailedGrpcError(serviceID int, err error) error {
	if err == nil {
		return nil
	}

	var typedErr errs.Error
	if !errors.As(err, &typedErr) {
		return createDetailedError(
			StatusErrorSystemFailure,
			err.Error(),
			"0", // no ErrNum for system errors
			nil,
		)
	}

	// Get GRPC status code from error type.
	mapper := getGrpcErrorMapper()
	code, exists := mapper[typedErr.Type]
	if !exists {
		code = StatusErrorDefault
	}
	var errNum string
	if typedErr.ErrNum != "" {
		errNum = fmt.Sprintf("%d.%s", serviceID, typedErr.ErrNum)
	} else {
		errNum = fmt.Sprintf("%d.0", serviceID)
	}
	details := typedErr.Details

	return createDetailedError(
		code,
		typedErr.Error(),
		errNum,
		details,
	)
}

// Создает ошибку GRPC с дополнительными полями.
func createDetailedError(code codes.Code, message string, errNum string, details map[string]string) error {
	status := status.New(code, message)

	// Create error details
	errorInfo := &errdetails.ErrorInfo{
		Reason: message,
		Metadata: map[string]string{
			"error_num": errNum,
		},
	}

	for key, value := range details {
		errorInfo.Metadata[DetailedErrorPrefix+key] = value
	}

	detailedStatus, err := status.WithDetails(errorInfo)
	if err != nil {
		err = status.Err()
	} else {
		err = detailedStatus.Err()
	}

	return err // nolint:wrapcheck // already should be wrapped
}

// getGrpcErrorMapper - возвращает маппер типов пользовательских ошибок на gRPC коды статусов.
// Маппер позволяет конвертировать различные типы ошибок (например, ошибка аутентификации или конфликта)
// в подходящие для gRPC статусы.
func getGrpcErrorMapper() map[errs.Type]codes.Code {
	return map[errs.Type]codes.Code{
		errs.TypeIllegalArgument:    codes.InvalidArgument,
		errs.TypeAuthFailure:        codes.Unauthenticated,
		errs.TypeForbidden:          codes.PermissionDenied,
		errs.TypeNotFound:           codes.NotFound,
		errs.TypeConflict:           codes.AlreadyExists,
		errs.TypeCancelled:          codes.Canceled,
		errs.TypeHasReferences:      codes.FailedPrecondition,
		errs.TypeSystemFailure:      codes.Internal,
		errs.TypeNotImplemented:     codes.Unimplemented,
		errs.TypeServiceUnavailable: codes.Unavailable,
		errs.TypeTooManyRequests:    codes.ResourceExhausted,
		errs.TypeUnauthenticated:    codes.Unauthenticated,
		errs.TypePermissionDenied:   codes.PermissionDenied,
		errs.TypeResourceExhausted:  codes.ResourceExhausted,
		errs.TypeFailedPrecondition: codes.FailedPrecondition,
		errs.TypeAborted:            codes.Aborted,
		errs.TypeOutOfRange:         codes.OutOfRange,
		errs.TypeUnimplemented:      codes.Unimplemented,
		errs.TypeInternal:           codes.Internal,
		errs.TypeUnavailable:        codes.Unavailable,
		errs.TypeDataLoss:           codes.DataLoss,
		errs.TypeCanceled:           codes.Canceled,
	}
}

// ExtractErrorNum извлекает из ошибки статус gRPC и возвращает его в виде строки.
func ExtractErrorNum(err error) string {
	if err == nil {
		return ""
	}

	status, ok := status.FromError(err)
	if !ok {
		return ""
	}

	// Extract error details
	for _, detail := range status.Details() {
		if errorInfo, ok := detail.(*errdetails.ErrorInfo); ok {
			if errNum, exists := errorInfo.GetMetadata()["error_num"]; exists {
				return errNum
			}
		}
	}

	// If no error details found, return the original error
	return ""
}

func ExtractDetails(status status.Status) map[string]string {
	// get error_num from status details using errdetails.FromStatus
	customDetails := make(map[string]string)
	count := 0
	for _, detail := range status.Details() {
		if errorInfo, ok := detail.(*errdetails.ErrorInfo); ok {
			for key, value := range errorInfo.GetMetadata() {
				if strings.HasPrefix(key, DetailedErrorPrefix) {
					customDetails[strings.TrimPrefix(key, DetailedErrorPrefix)] = value
					count++
				}
			}
		}
	}
	if count == 0 {
		customDetails = nil
	}

	return customDetails
}

// ErrorFromGRPC преобразует ошибку gRPC в ошибку системной логики.
func ErrorFromGRPC(err error) errs.Error {
	status, ok := status.FromError(err)
	if ok {
		// get error_num from status details using errdetails.FromStatus
		var grpcError errs.Error
		errNum := ExtractErrorNum(err)
		details := ExtractDetails(*status)

		code := status.Code()
		switch {
		case code == codes.InvalidArgument:
			grpcError = errs.Reasons(status.Message(), errs.TypeIllegalArgument, errNum)
		case code == codes.Unauthenticated:
			grpcError = errs.Reasons(status.Message(), errs.TypeAuthFailure, errNum)
		case code == codes.PermissionDenied:
			grpcError = errs.Reasons(status.Message(), errs.TypeForbidden, errNum)
		case code == codes.NotFound:
			grpcError = errs.Reasons(status.Message(), errs.TypeNotFound, errNum)
		case code == codes.AlreadyExists:
			grpcError = errs.Reasons(status.Message(), errs.TypeConflict, errNum)
		case code == codes.Canceled:
			grpcError = errs.Reasons(status.Message(), errs.TypeCancelled, errNum)
		case code == codes.FailedPrecondition:
			grpcError = errs.Reasons(status.Message(), errs.TypeHasReferences, errNum)
		case code == codes.Unimplemented:
			grpcError = errs.Reasons(status.Message(), errs.TypeNotImplemented, errNum)
		case code == codes.Unavailable:
			grpcError = errs.Reasons(status.Message(), errs.TypeServiceUnavailable, errNum)
		case code == codes.ResourceExhausted:
			grpcError = errs.Reasons(status.Message(), errs.TypeResourceExhausted, errNum)
		case code == codes.OutOfRange:
			grpcError = errs.Reasons(status.Message(), errs.TypeOutOfRange, errNum)
		case code == codes.Internal:
			grpcError = errs.Reasons(status.Message(), errs.TypeInternal, errNum)
		case code == codes.DataLoss:
			grpcError = errs.Reasons(status.Message(), errs.TypeDataLoss, errNum)
		case code == codes.DeadlineExceeded:
			grpcError = errs.Reasons(status.Message(), errs.TypeTooManyRequests, errNum)
		case code == codes.Aborted:
			grpcError = errs.Reasons(status.Message(), errs.TypeAborted, errNum)
		// Default case - unchanged
		default:
			if errNum == "" {
				grpcError = errs.Reasons(status.Message(), errs.TypeSystemFailure, "1.0")
			} else {
				grpcError = errs.Reasons(status.Message(), errs.TypeSystemFailure, errNum)
			}
		}
		if details != nil {
			return grpcError.WithDetails(details)
		}

		return grpcError
	}

	return errs.Reasons("Fail", errs.TypeSystemFailure, "1.0")
}
