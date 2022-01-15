package boot

import (
	"github.com/fooree/boot/app"
)

// Application 单例对象
var _app_ *app.Application

func _app() *app.Application {
	if _app_ == nil {
		_app_ = app.NewApplication()
	}
	return _app_
}

func Run(name string) {
	_app().Named(name).Run()
}

func RegisterApplicationListener(listener app.ApplicationListener) {
	_app().RegisterApplicationListener(listener)
}

func RegisterObject(object interface{}) {
	_app().RegisterObject(object)
}
