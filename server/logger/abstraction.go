package logger

import "context"

//ILogger ...
type ILogger interface {
	Instrumentation(
		method string,
		action string,
		f func(ctx context.Context, request interface{}) (interface{}, error),
	) func(ctx context.Context, request interface{}) (interface{}, error)
	Log(method string, action string, f func(ctx context.Context, request interface{}) (interface{}, error)) func(ctx context.Context, request interface{}) (interface{}, error)
}
