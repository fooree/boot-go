package container

import (
	"context"
)

type Container interface {
	Context() context.Context
	RegisterObject(value interface{}) *Bean
	Close()
}
