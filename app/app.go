package app

import (
	"github.com/fooree/boot/container"
	"github.com/fooree/boot/context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ApplicationListener interface {
	OnApplicationStart(context context.Context) error // 启动事件监听
	OnApplicationStop(context context.Context)        // 停止事件监听
}

type Application struct {
	name      string
	container *container.ContainerImpl
	exit      chan struct{}
	listeners []ApplicationListener
}

func NewApplication() *Application {
	return &Application{
		container: container.NewContainer(),
		exit:      make(chan struct{}, 1),
		listeners: make([]ApplicationListener, 0, 16),
	}
}

func (app *Application) Run() {
	go app.monitorSignal()

	if err := app.start(); err == nil {
		log.Println("application is started")
		<-app.exit

		log.Println("application exiting")
		app.stop()
		log.Println("application exited")

		time.Sleep(time.Second * 5)
	} else {
		panic(err)
	}
}

func (app *Application) start() (err error) {
	err = app.container.Start()

	if err == nil {
		for i := 0; i < len(app.listeners); i++ {
			err = app.listeners[i].OnApplicationStart(app.container)
			if err == nil {
				continue
			}
			break
		}
	}

	return
}

func (app *Application) stop() {
	app.container.Stop()

	for i := 0; i < len(app.listeners); i++ {
		app.listeners[i].OnApplicationStop(app.container)
	}
}

func (app *Application) RegisterApplicationListener(listener ApplicationListener) {
	if listener == nil {
		panic("register nil application listener")
	}
	app.listeners = append(app.listeners, listener)
}

func (app *Application) Named(name string) *Application {
	if app.name == "" {
		app.name = name
		return app
	} else {
		panic("set application name repeatedly")
	}
}

func (app *Application) RegisterObject(object interface{}) {
	if nil == object {
		panic("register nil object")
	}
	app.container.RegisterObject(object)
}

func (app *Application) monitorSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer close(ch)
	defer close(app.exit)

	log.Printf("monitoring OS signal")
	sig := <-ch
	log.Printf("OS signal is '%s'", sig)
}
