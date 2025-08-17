package grpcx

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wal1251/pkg/core/errs"
)

func TestNewGrpcError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name: "MappedErrorIllegalArgument",
			args: args{
				err: errs.Error{
					Code: "invalid argument",
					Type: errs.TypeIllegalArgument,
				},
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "MappedErrorAuthFailure",
			args: args{
				err: errs.Error{
					Code: "authentication failed",
					Type: errs.TypeAuthFailure,
				},
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
		{
			name: "UnmappedCustomError",
			args: args{
				err: errs.Error{
					Code: "some custom error",
					Type: errs.Type("UnknownError"),
				},
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "StandardError",
			args: args{
				err: status.Error(codes.Unavailable, "service unavailable"),
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "NilError",
			args: args{
				err: nil,
			},
			wantErr:  false,
			wantCode: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewGrpcError(tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGrpcError() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				st, _ := status.FromError(err)
				if st.Code() != tt.wantCode {
					t.Errorf("NewGrpcError() code = %v, wantCode %v", st.Code(), tt.wantCode)
				}
			}
		})
	}
}
