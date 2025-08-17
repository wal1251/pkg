package tests

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/proxy"
)

func TestProxy_noArgsNoReturn(t *testing.T) {
	var sequence *[]string

	hook := func() {
		s := append(*sequence, "NoArgsAndNoReturn")
		sequence = &s
	}

	tests := []struct {
		name         string
		object       *SampleObject
		middlewares  []proxy.MethodInvocationMiddleware
		wantSequence []string
	}{
		{
			name:         "Вызов метода без параметров и результатов без middleware",
			object:       &SampleObject{HookNoArgsAndNoReturn: hook},
			wantSequence: []string{"NoArgsAndNoReturn"},
		},
		{
			name:   "Вызов метода без параметров и результатов с middleware",
			object: &SampleObject{HookNoArgsAndNoReturn: hook},
			middlewares: []proxy.MethodInvocationMiddleware{
				func(object any, name string, args []reflect.Value, next proxy.GenericFunction) []reflect.Value {
					s := append(*sequence, "Before1")
					sequence = &s
					result := next(args)
					s = append(*sequence, "After1")
					sequence = &s
					return result
				},
				func(object any, name string, args []reflect.Value, next proxy.GenericFunction) []reflect.Value {
					s := append(*sequence, "Before2")
					sequence = &s
					result := next(args)
					s = append(*sequence, "After2")
					sequence = &s
					return result
				},
			},
			wantSequence: []string{"Before1", "Before2", "NoArgsAndNoReturn", "After2", "After1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := make([]string, 0)
			sequence = &s

			p := NewSampleInterfaceProxy(tt.object, tt.middlewares...)
			p.NoArgsAndNoReturn()
			assert.Equal(t, tt.wantSequence, *sequence, "Последовательность вызовов не равна ожидаемой")
		})
	}
}

func TestProxy_ArgsAndReturn(t *testing.T) {
	var sequence *[]string

	hook := func(arg1, arg2 string) (string, error) {
		result := fmt.Sprintf("WithArgsAndReturn,%s,%s", arg1, arg2)
		s := append(*sequence, result)
		sequence = &s
		return result, errors.New("fake err")
	}

	tests := []struct {
		name         string
		object       *SampleObject
		middlewares  []proxy.MethodInvocationMiddleware
		arg1         string
		arg2         string
		wantResult1  string
		wantResult2  error
		wantSequence []string
	}{
		{
			name:         "Вызов метода с параметрами и результатами без middleware",
			object:       &SampleObject{HookWithArgsAndReturn: hook},
			arg1:         "foo",
			arg2:         "bar",
			wantResult1:  "WithArgsAndReturn,foo,bar",
			wantResult2:  errors.New("fake err"),
			wantSequence: []string{"WithArgsAndReturn,foo,bar"},
		},
		{
			name:   "Вызов метода с параметрами и результатами с middleware",
			object: &SampleObject{HookWithArgsAndReturn: hook},
			middlewares: []proxy.MethodInvocationMiddleware{
				func(object any, name string, args []reflect.Value, next proxy.GenericFunction) []reflect.Value {
					s := append(*sequence, "Before1")
					sequence = &s
					result := next(args)
					s = append(*sequence, "After1")
					sequence = &s
					return result
				},
				func(object any, name string, args []reflect.Value, next proxy.GenericFunction) []reflect.Value {
					s := append(*sequence, "Before2")
					sequence = &s
					result := next(args)
					s = append(*sequence, "After2")
					sequence = &s
					return result
				},
			},
			arg1:         "foo",
			arg2:         "bar",
			wantResult1:  "WithArgsAndReturn,foo,bar",
			wantResult2:  errors.New("fake err"),
			wantSequence: []string{"Before1", "Before2", "WithArgsAndReturn,foo,bar", "After2", "After1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := make([]string, 0)
			sequence = &s

			p := NewSampleInterfaceProxy(tt.object, tt.middlewares...)
			result, err := p.WithArgsAndReturn(tt.arg1, tt.arg2)
			assert.Equal(t, tt.wantResult1, result, "Результат 1 не равен ожидаемому")
			assert.Equal(t, tt.wantResult2, err, "Результат 2 не равен ожидаемому")
			assert.Equal(t, tt.wantSequence, *sequence, "Последовательность вызовов не равна ожидаемой")
		})
	}
}

func TestProxy_panicInterceptedAndPropagated(t *testing.T) {
	var sequence *[]string

	tests := []struct {
		name         string
		proxy        SampleInterface
		wantSequence []string
	}{
		{
			name: "Паника в методе объекта",
			proxy: NewSampleInterfaceProxyWithPanicHook(&SampleObject{
				HookNoArgsAndNoReturn: func() {
					s := append(*sequence, "Invoke")
					sequence = &s
					panic(errors.New("fake panic"))
				},
			}, func(err any, stack []byte) any {
				s := append(*sequence, "PanicRecover")
				sequence = &s
				return err
			}),
			wantSequence: []string{"Invoke", "PanicRecover"},
		},
		{
			name: "Паника в middleware",
			proxy: NewSampleInterfaceProxyWithPanicHook(&SampleObject{
				HookNoArgsAndNoReturn: func() {
					s := append(*sequence, "Invoke")
					sequence = &s
				},
			},
				func(err any, stack []byte) any {
					s := append(*sequence, "PanicRecover")
					sequence = &s
					return err
				},
				func(object any, name string, args []reflect.Value, next proxy.GenericFunction) []reflect.Value {
					s := append(*sequence, "Before1")
					sequence = &s
					panic(errors.New("fake panic"))
				},
				func(object any, name string, args []reflect.Value, next proxy.GenericFunction) []reflect.Value {
					s := append(*sequence, "Before2")
					sequence = &s
					result := next(args)
					s = append(*sequence, "After2")
					sequence = &s
					return result
				},
			),
			wantSequence: []string{"Before1", "PanicRecover"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := make([]string, 0)
			sequence = &s
			if assert.Panics(t, tt.proxy.NoArgsAndNoReturn, "Паника ожидалась") {
				assert.Equal(t, tt.wantSequence, *sequence, "Последовательность вызовов не равна ожидаемой")
			}
		})
	}
}

func TestProxy_panicInterceptedAndRecovered(t *testing.T) {
	var sequence *[]string

	tests := []struct {
		name         string
		proxy        SampleInterface
		wantSequence []string
	}{
		{
			name: "Восстановление после паники в методе объекта",
			proxy: NewSampleInterfaceProxyWithPanicHook(&SampleObject{
				HookNoArgsAndNoReturn: func() {
					s := append(*sequence, "Invoke")
					sequence = &s
					panic(errors.New("fake panic"))
				},
			}, func(err any, stack []byte) any {
				s := append(*sequence, "PanicRecover")
				sequence = &s
				return nil
			}),
			wantSequence: []string{"Invoke", "PanicRecover"},
		},
		{
			name: "Восстановление после паники в middleware",
			proxy: NewSampleInterfaceProxyWithPanicHook(&SampleObject{
				HookNoArgsAndNoReturn: func() {
					s := append(*sequence, "Invoke")
					sequence = &s
				},
			},
				func(err any, stack []byte) any {
					s := append(*sequence, "PanicRecover")
					sequence = &s
					return nil
				},
				func(object any, name string, args []reflect.Value, next proxy.GenericFunction) []reflect.Value {
					s := append(*sequence, "Before1")
					sequence = &s
					panic(errors.New("fake panic"))
				},
				func(object any, name string, args []reflect.Value, next proxy.GenericFunction) []reflect.Value {
					s := append(*sequence, "Before2")
					sequence = &s
					result := next(args)
					s = append(*sequence, "After2")
					sequence = &s
					return result
				},
			),
			wantSequence: []string{"Before1", "PanicRecover"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := make([]string, 0)
			sequence = &s
			if assert.NotPanics(t, tt.proxy.NoArgsAndNoReturn, "Паника не ожидалась") {
				assert.Equal(t, tt.wantSequence, *sequence, "Последовательность вызовов не равна ожидаемой")
			}
		})
	}
}

func BenchmarkProxy(b *testing.B) {
	p := NewSampleInterfaceProxy(&SampleObject{
		HookNoArgsAndNoReturn: func() {},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.NoArgsAndNoReturn()
	}
}

func BenchmarkProxyArgs(b *testing.B) {
	p := NewSampleInterfaceProxy(&SampleObject{
		HookWithArgsAndReturn: func(arg1, arg2 string) (string, error) {
			return arg1 + arg2, nil
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.NoArgsAndNoReturn()
	}
}

func BenchmarkNoProxy(b *testing.B) {
	p := &SampleObject{
		HookNoArgsAndNoReturn: func() {},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.NoArgsAndNoReturn()
	}
}
