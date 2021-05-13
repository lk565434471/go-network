package application

import (
	"go-network/pattern"
	"go-network/utils"
)

type Application struct {
	rootPath string
}

func (app *Application) Init() bool {
	app.rootPath = utils.GetExecutableRootPath()

	return true
}

func (app *Application) Run() {

}

func (app *Application) GetRootPath() string {
	return app.rootPath
}

var instance = pattern.NewSingleton(func() (interface{}, bool) {
	return &Application{
	}, true
})

func GetApp() func() *Application {
	return func() *Application {
		obj, ok := instance.Get()

		if !ok {
			return nil
		}

		app, ok := obj.(*Application)

		if !ok {
			return nil
		}

		return app
	}
}