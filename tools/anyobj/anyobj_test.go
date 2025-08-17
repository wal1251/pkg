package anyobj

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeCopy(t *testing.T) {
	type testStruct struct {
		Field1 string
		Field2 int
	}

	type args struct {
		src  any
		dest any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful copy - struct to struct",
			args: args{
				src:  testStruct{"Hello, world!", 42},
				dest: &testStruct{},
			},
			wantErr: false,
		},
		{
			name: "Successful copy - map to struct",
			args: args{
				src: map[string]any{
					"Field1": "Hello, world!",
					"Field2": 42,
				},
				dest: &testStruct{},
			},
			wantErr: false,
		},
		{
			name: "Failed copy - dest is not pointer",
			args: args{
				src:  testStruct{"Hello, world!", 42},
				dest: testStruct{},
			},
			wantErr: true,
		},
		{
			name: "Failed copy - dest is nil",
			args: args{
				src:  testStruct{"Hello, world!", 42},
				dest: nil,
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := SafeCopy(tc.args.src, tc.args.dest)
			if tc.wantErr {
				require.Error(t, err)
			}

			assert.ObjectsAreEqual(tc.args.src, tc.args.dest)
		})
	}
}

func TestSafeCopy_UnexportedFields(t *testing.T) {
	type testStructWithUnexported struct {
		Field1 string
		Field2 int
		field3 bool // неэкспортируемое поле
	}

	type args struct {
		src  any
		dest any
	}
	tests := []struct {
		name string
		args args
		want any // Ожидаемый результат после копирования
	}{
		{
			name: "struct to struct with unexported fields",
			args: args{
				src: testStructWithUnexported{
					Field1: "Test",
					Field2: 123,
					field3: true,
				},
				dest: &testStructWithUnexported{},
			},
			want: &testStructWithUnexported{
				Field1: "Test",
				Field2: 123,
				field3: false, // field3 остаётся в значении по умолчанию, потому что не копируется
			},
		},
		{
			name: "map to struct with unexported fields",
			args: args{
				src: map[string]any{
					"Field1": "Test",
					"Field2": 123,
					"field3": true,
				},
				dest: &testStructWithUnexported{},
			},
			want: &testStructWithUnexported{
				Field1: "Test",
				Field2: 123,
				field3: false, // field3 остаётся в значении по умолчанию, потому что не копируется
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := SafeCopy(tc.args.src, tc.args.dest)
			require.NoError(t, err)

			assert.ObjectsAreEqual(tc.args.dest, tc.want)
		})
	}
}
