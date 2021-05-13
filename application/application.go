package application

import (
	"errors"
	"fmt"
	"go-network/pattern"
	"go-network/utils"
	"strings"
)

var defaultCommands map[string]bool = map[string]bool{
	"quit": true,
}

type CommandFunc func([]string)

type Application struct {
	rootPath string
	stop chan struct{}
	commands map[string]CommandFunc
}

func (app *Application) Init() bool {
	app.rootPath = utils.GetExecutableRootPath()

	app.addDefaultCommands()

	return true
}

func (app *Application) Run() {
	for  {
		select {
		case <-app.stop:
			return
		default:
			var str string
			fmt.Scanln(&str)
			app.executeCommand(str)
		}
	}
}

func (app *Application) executeCommand(str string) {
	paramParts := strings.Split(str, " ")

	if len(paramParts) == 0 {
		return
	}

	cmd := paramParts[0]
	cmdFunc, ok := app.commands[cmd]

	if !ok {
		return
	}

	cmdFunc(paramParts[1:])
}

func (app *Application) AddCommand(name string, cmd CommandFunc) error {
	_, ok := defaultCommands[name]

	if ok {
		strErr := fmt.Sprintf("Command: %s is a keyword.", name)
		return errors.New(strErr)
	}

	app.commands[name] = cmd
	return nil
}

func (app *Application) RemoveCommand(name string) error {
	_, ok := defaultCommands[name]

	if ok {
		strErr := fmt.Sprintf("Command: %s is a keyword.", name)
		return errors.New(strErr)
	}

	delete(app.commands, name)
	return nil
}

func (app *Application) addDefaultCommands() {
	app.commands["quit"] = app.onQuit
}

func (app *Application) onQuit([]string) {
	close(app.stop)
}

func (app *Application) GetRootPath() string {
	return app.rootPath
}

var instance = pattern.NewSingleton(pattern.SingletonSettings{
	OnInit: func() (interface{}, bool) {
		return &Application{
			stop: make(chan struct{}),
			commands: map[string]CommandFunc{},
		}, true
	},
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