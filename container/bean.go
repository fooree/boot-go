package container

import "reflect"

type Initializer interface {
	Initialize() error
}

type Finalizer interface {
	Finalize()
}

type Bean struct {
	name string
	typ  reflect.Type
	val  reflect.Value

	args []reflect.Value

	init     interface{}
	finalize interface{}
}

func (b *Bean) Named(name string) *Bean {
	if name == "" {
		panic("empty name")
	}
	if b.name == "" {
		b.name = name
		return b
	} else {
		panic("named bean repeatedly")
	}
}
