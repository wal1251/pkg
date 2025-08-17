package tests

import (
	"context"
	"fmt"
	"strings"
)

type SampleInterface interface {
	WithArgsAndReturn(string, string) (string, error)
	NoArgsAndNoReturn()
	WithContextArgAndReturn(context.Context, string, string) (string, error)
	WithContextVarArgAndReturn(context.Context, string, ...string) (string, error)
}

var _ SampleInterface = (*SampleObject)(nil)

type SampleObject struct {
	HookWithArgsAndReturn func(arg1, arg2 string) (string, error)
	HookNoArgsAndNoReturn func()
}

func (o *SampleObject) WithContextVarArgAndReturn(ctx context.Context, arg string, args ...string) (string, error) {
	return fmt.Sprint(ctx.Value(arg), " ", strings.Join(args, " ")), nil
}

func (o *SampleObject) WithContextArgAndReturn(ctx context.Context, arg1, arg2 string) (string, error) {
	return fmt.Sprint(ctx.Value(arg1), " ", arg2), nil
}

func (o *SampleObject) WithArgsAndReturn(arg1, arg2 string) (string, error) {
	if o.HookWithArgsAndReturn != nil {
		return o.HookWithArgsAndReturn(arg1, arg2)
	}

	return "", nil
}

func (o *SampleObject) NoArgsAndNoReturn() {
	if o.HookNoArgsAndNoReturn != nil {
		o.HookNoArgsAndNoReturn()
	}
}
