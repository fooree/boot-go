package context

import "context"

type Context interface {
	Context() context.Context
	Exec(f func(context Context))
	Get(key string) interface{}
	Set(key string, value interface{})
	CommandLineValue(key string) (val string, ok bool)
}
