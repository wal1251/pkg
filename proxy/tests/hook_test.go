package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/proxy"
)

func TestHook_noArgsNoReturn(t *testing.T) {
	var sequence *[]string

	hook := func() {
		s := append(*sequence, "NoArgsAndNoReturn")
		sequence = &s
	}

	tests := []struct {
		name         string
		object       *SampleObject
		onBeforeCall proxy.Hook
		onPostCall   proxy.Hook
		wantSequence []string
	}{
		{
			name:         "Вызов метода без параметров и результатов без хуков",
			object:       &SampleObject{HookNoArgsAndNoReturn: hook},
			wantSequence: []string{"NoArgsAndNoReturn"},
		},
		{
			name:   "Вызов метода без параметров и результатов с хуками",
			object: &SampleObject{HookNoArgsAndNoReturn: hook},
			onBeforeCall: proxy.Hook(func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "Before1")
				sequence = &s
				return ctx
			}).And(func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "Before2")
				sequence = &s
				return ctx
			}, func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "Before3")
				sequence = &s
				return ctx
			}),
			onPostCall: proxy.Hook(func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "After1")
				sequence = &s
				return ctx
			}).And(func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "After2")
				sequence = &s
				return ctx
			}, func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "After3")
				sequence = &s
				return ctx
			}),
			wantSequence: []string{"Before1", "Before2", "Before3", "NoArgsAndNoReturn", "After1", "After2", "After3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := make([]string, 0)
			sequence = &s

			p := NewSampleInterfaceHook(tt.object, tt.onBeforeCall, tt.onPostCall, nil)
			p.NoArgsAndNoReturn()
			assert.Equal(t, tt.wantSequence, *sequence, "Последовательность вызовов не равна ожидаемой")
		})
	}
}

func TestHook_WithContextArgAndReturn(t *testing.T) {
	tests := []struct {
		name         string
		object       *SampleObject
		onBeforeCall proxy.Hook
		arg1         string
		arg2         string
		wantResult1  string
		wantResult2  error
	}{
		{
			name:   "Вызов метода с контекстом",
			object: &SampleObject{},
			onBeforeCall: proxy.Hook(func(ctx context.Context, object any, method string, args []any) context.Context {
				return context.WithValue(ctx, "foo", fmt.Sprint("foo ", ctx.Value("foo")))
			}),
			arg1:        "foo",
			arg2:        "bar",
			wantResult1: "foo baz bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, tt.arg1, "baz")
			p := NewSampleInterfaceHook(tt.object, tt.onBeforeCall, nil, nil)
			result, err := p.WithContextArgAndReturn(ctx, tt.arg1, tt.arg2)
			assert.Equal(t, tt.wantResult1, result, "Результат 1 не равен ожидаемому")
			assert.Nil(t, tt.wantResult2, err, "Результат 2 не равен ожидаемому")
		})
	}
}

func TestHook_WithContextVarArgAndReturn(t *testing.T) {
	tests := []struct {
		name         string
		object       *SampleObject
		onBeforeCall proxy.Hook
		arg1         string
		args         []string
		wantResult1  string
		wantResult2  error
	}{
		{
			name:   "Вызов метода с контекстом",
			object: &SampleObject{},
			onBeforeCall: proxy.Hook(func(ctx context.Context, object any, method string, args []any) context.Context {
				return context.WithValue(ctx, "foo", fmt.Sprint("foo ", ctx.Value("foo")))
			}),
			arg1:        "foo",
			args:        []string{"foo", "bar", "baz"},
			wantResult1: "foo baz foo bar baz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, tt.arg1, "baz")
			p := NewSampleInterfaceHook(tt.object, tt.onBeforeCall, nil, nil)
			result, err := p.WithContextVarArgAndReturn(ctx, tt.arg1, tt.args...)
			assert.Equal(t, tt.wantResult1, result, "Результат 1 не равен ожидаемому")
			assert.Nil(t, tt.wantResult2, err, "Результат 2 не равен ожидаемому")
		})
	}
}

func TestHook_ArgsAndReturn(t *testing.T) {
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
		onBeforeCall proxy.Hook
		onPostCall   proxy.Hook
		arg1         string
		arg2         string
		wantResult1  string
		wantResult2  error
		wantSequence []string
	}{
		{
			name:         "Вызов метода с параметрами и результатами без хуков",
			object:       &SampleObject{HookWithArgsAndReturn: hook},
			arg1:         "foo",
			arg2:         "bar",
			wantResult1:  "WithArgsAndReturn,foo,bar",
			wantResult2:  errors.New("fake err"),
			wantSequence: []string{"WithArgsAndReturn,foo,bar"},
		},
		{
			name:   "Вызов метода с параметрами и результатами с хуками",
			object: &SampleObject{HookWithArgsAndReturn: hook},
			onBeforeCall: proxy.Hook(func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "Before1")
				sequence = &s
				return ctx
			}).And(func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "Before2")
				sequence = &s
				return ctx
			}, func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "Before3")
				sequence = &s
				return ctx
			}),
			onPostCall: proxy.Hook(func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "After1")
				sequence = &s
				return ctx
			}).And(func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "After2")
				sequence = &s
				return ctx
			}, func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "After3")
				sequence = &s
				return ctx
			}),
			arg1:         "foo",
			arg2:         "bar",
			wantResult1:  "WithArgsAndReturn,foo,bar",
			wantResult2:  errors.New("fake err"),
			wantSequence: []string{"Before1", "Before2", "Before3", "WithArgsAndReturn,foo,bar", "After1", "After2", "After3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := make([]string, 0)
			sequence = &s

			p := NewSampleInterfaceHook(tt.object, tt.onBeforeCall, tt.onPostCall, nil)
			result, err := p.WithArgsAndReturn(tt.arg1, tt.arg2)
			assert.Equal(t, tt.wantResult1, result, "Результат 1 не равен ожидаемому")
			assert.Equal(t, tt.wantResult2, err, "Результат 2 не равен ожидаемому")
			assert.Equal(t, tt.wantSequence, *sequence, "Последовательность вызовов не равна ожидаемой")
		})
	}
}

func TestHook_panicInterceptedAndPropagated(t *testing.T) {
	var sequence *[]string

	tests := []struct {
		name         string
		proxy        SampleInterface
		wantSequence []string
	}{
		{
			name: "Паника в методе объекта",
			proxy: NewSampleInterfaceHook(&SampleObject{
				HookNoArgsAndNoReturn: func() {
					s := append(*sequence, "Invoke")
					sequence = &s
					panic(errors.New("fake panic"))
				},
			}, nil, nil, func(msg any, stack []byte, object any, method string, args []any) any {
				s := append(*sequence, "PanicRecover")
				sequence = &s
				return msg
			}),
			wantSequence: []string{"Invoke", "PanicRecover"},
		},
		{
			name: "Паника в хуке",
			proxy: NewSampleInterfaceHook(&SampleObject{
				HookNoArgsAndNoReturn: func() {
					s := append(*sequence, "Invoke")
					sequence = &s
				},
			}, func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "Before")
				sequence = &s
				panic(errors.New("fake panic"))
			}, func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "After")
				sequence = &s
				return ctx
			}, func(msg any, stack []byte, object any, method string, args []any) any {
				s := append(*sequence, "PanicRecover")
				sequence = &s
				return msg
			}),
			wantSequence: []string{"Before", "PanicRecover"},
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

func TestHook_panicInterceptedAndRecovered(t *testing.T) {
	var sequence *[]string

	tests := []struct {
		name         string
		proxy        SampleInterface
		wantSequence []string
	}{
		{
			name: "Восстановление после паники в методе объекта",
			proxy: NewSampleInterfaceHook(&SampleObject{
				HookNoArgsAndNoReturn: func() {
					s := append(*sequence, "Invoke")
					sequence = &s
					panic(errors.New("fake panic"))
				},
			}, nil, nil, func(msg any, stack []byte, object any, method string, args []any) any {
				s := append(*sequence, "PanicRecover")
				sequence = &s
				return nil
			}),
			wantSequence: []string{"Invoke", "PanicRecover"},
		},
		{
			name: "Восстановление после паники в хуке",
			proxy: NewSampleInterfaceHook(&SampleObject{
				HookNoArgsAndNoReturn: func() {
					s := append(*sequence, "Invoke")
					sequence = &s
				},
			}, func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "Before")
				sequence = &s
				panic(errors.New("fake panic"))
			}, func(ctx context.Context, object any, method string, args []any) context.Context {
				s := append(*sequence, "After")
				sequence = &s
				return ctx
			}, func(msg any, stack []byte, object any, method string, args []any) any {
				s := append(*sequence, "PanicRecover")
				sequence = &s
				return nil
			}),
			wantSequence: []string{"Before", "PanicRecover"},
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

func BenchmarkHook(b *testing.B) {
	p := NewSampleInterfaceHook(&SampleObject{
		HookNoArgsAndNoReturn: func() {},
	}, nil, nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.NoArgsAndNoReturn()
	}
}

func BenchmarkHookArgs(b *testing.B) {
	p := NewSampleInterfaceHook(&SampleObject{
		HookWithArgsAndReturn: func(arg1, arg2 string) (string, error) {
			return arg1 + arg2, nil
		},
	}, nil, nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.NoArgsAndNoReturn()
	}
}

func BenchmarkNoHook(b *testing.B) {
	p := &SampleObject{
		HookNoArgsAndNoReturn: func() {},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.NoArgsAndNoReturn()
	}
}
