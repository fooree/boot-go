package container

import (
	"context"
	"fmt"
	context2 "github.com/fooree/boot/context"
	"github.com/fooree/boot/reflects"
	"log"
	"os"
	"reflect"
	"strings"
)

type ContainerImpl struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cfg        *configuration
	beans      []*Bean
	namedBeans map[string]map[string]*Bean
	typedBeans map[reflect.Type][]*Bean
}

func (c *ContainerImpl) Exec(f func(context context2.Context)) {
	//TODO implement me
	return
}

func (c *ContainerImpl) Get(key string) interface{} {

	if v, ok := c.cfg.commandArgs[key]; ok { // 命令行参数优先
		return v
	}
	return nil
}

func (c *ContainerImpl) Set(key string, value interface{}) {
	//TODO implement me
}

func (c *ContainerImpl) Context() context.Context {
	return nil
}

func (c *ContainerImpl) RegisterObject(obj interface{}) *Bean {
	bean := c.newBean(obj)
	c.beans = append(c.beans, bean)
	return bean
}

func (c *ContainerImpl) Close() {
	c.cancel()
}

func (c *ContainerImpl) newBean(obj interface{}) *Bean {
	if obj == nil {
		panic("newBean bean from nil object")
	}

	val := reflects.ValueOf(obj)
	typ := val.Type()

	return &Bean{
		typ: typ,
		val: val,
	}
}

func (c *ContainerImpl) Start() (err error) {
	c.cfg.ReadCommandArgs(os.Args[1:])

	// bind beans
	err = c.bindBeans()

	// init beans
	err = c.initBeans()

	return nil
}

func (c *ContainerImpl) initBeans() (err error) {
	for _, bean := range c.beans {
		if obj, ok := bean.val.Interface().(Initializer); ok {
			err = obj.Initialize()
			if err == nil {
				continue
			} else {
				break
			}
		}
	}
	return
}

func (c *ContainerImpl) Stop() {
	c.cancel()

	for _, bean := range c.beans {
		if obj, ok := bean.val.Interface().(Finalizer); ok {
			obj.Finalize()
		}
	}
}

func (c *ContainerImpl) bindBeans() (err error) {
	err = c.findNamedBeans()
	if err == nil {
		for i := 0; i < len(c.beans); i++ {
			bean := c.beans[i]
			typ := bean.typ
			val := bean.val
			for typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
				val = val.Elem()
			}
			if typ.Kind() == reflect.Struct {
				num := typ.NumField()
				for j := 0; j < num; j++ {
					//f := typ.Field(i)
					//cast.ToString(8)
					//if name, ok := f.Tag.Lookup(Value); ok {
					//	v, o := c.CommandLineValue(name)
					//	if o {
					//
					//		kind := f.Type.Kind()
					//		switch kind {
					//
					//		}
					//	}
					//
					//}
				}
			}
		}
	}
	return
}

func (c *ContainerImpl) findNamedBeans() (err error) {
	for i := 0; i < len(c.beans); i++ {
		bean := c.beans[i]
		name := bean.name
		if name == "" {
			continue
		}
		if _, ok := c.namedBeans[name]; !ok {
			c.namedBeans[name] = make(map[string]*Bean)
		}
		typeBeans := c.namedBeans[name]

		typeName := strings.TrimLeft(bean.typ.String(), "*")
		log.Printf("bean: name=%s, type=%s", name, typeName)

		if b, ok := typeBeans[typeName]; ok {
			if bean == b {
				err = fmt.Errorf("register the same bean repeatedly: name=%s, type=%s", name, typeName)
			} else {
				err = fmt.Errorf("register multiple beans with same name: name=%s, type=%s", name, typeName)
			}
		}

		c.namedBeans[name][typeName] = bean
	}
	return
}

func (c *ContainerImpl) CommandLineValue(key string) (val string, ok bool) {
	val, ok = c.cfg.commandArgs[key]
	return val, ok
}

func NewContainer() *ContainerImpl {
	ctx, cancel := context.WithCancel(context.Background())
	return &ContainerImpl{
		ctx:        ctx,
		cancel:     cancel,
		cfg:        newConfiguration(),
		beans:      make([]*Bean, 0, 16),
		namedBeans: make(map[string]map[string]*Bean, 16),
		typedBeans: make(map[reflect.Type][]*Bean, 16),
	}
}
