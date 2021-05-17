package application

import (
	"errors"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	logger "go-network/logrus"
	"go-network/pattern"
	"go-network/utils"
	"strings"
)

var defaultCommands map[string]bool = map[string]bool{
	"quit": true,
}

type CommandFunc func([]string)

type AppLogSettings struct {
	DisableLogging bool
	LoggingLevel logrus.Level
	LogDir string
	LogFilename string
	ReportCaller bool
	EnableStdout bool

	LogFormatter logrus.Formatter
	Options []rotatelogs.Option
}

type CustomLogSettings struct {
	LoggingLevel logrus.Level
	LogDir string
	LogFilename string
	ReportCaller bool
	EnableStdout bool

	LogFormatter logrus.Formatter
	Options []rotatelogs.Option
}

type AppSettings struct {
	AppLogSettings AppLogSettings
	CustomLogSettings CustomLogSettings
}

type Application struct {
	rootPath string
	stop chan struct{}
	commands map[string]CommandFunc
	DefaultLogger *logger.Logger
	Logger *logger.Logger
}

func (app *Application) Init(settings AppSettings) bool {
	app.rootPath = utils.GetExecutableRootPath()
	app.initAppLogger(settings.AppLogSettings)
	app.initCustomLogger(settings.CustomLogSettings)
	app.addDefaultCommands()
	app.Debug("Application::Init: Initialized success.")

	return true
}

func (app *Application) initAppLogger(settings AppLogSettings) {
	if settings.DisableLogging {
		return
	}

	app.DefaultLogger = logger.NewLogger(logger.LogSettings{
		Settings: settings.LogFormatter,
		LogDir: settings.LogDir,
		LogFilename: settings.LogFilename,
		LogLevel: settings.LoggingLevel,
		ReportCaller: settings.ReportCaller,
		EnableStdout: settings.EnableStdout,
	}, settings.Options...)
}

func (app *Application) initCustomLogger(settings CustomLogSettings) {
	app.Logger = logger.NewLogger(logger.LogSettings{
		Settings: settings.LogFormatter,
		LogDir: settings.LogDir,
		LogFilename: settings.LogFilename,
		LogLevel: settings.LoggingLevel,
		ReportCaller: settings.ReportCaller,
		EnableStdout: settings.EnableStdout,
	}, settings.Options...)
}

func (app *Application) Run() {
	app.DefaultLogger.Debug("Application::Run: App was started.")

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

func (app *Application) Debug(args ...interface{}) {
	if app.DefaultLogger == nil {
		return
	}

	app.DefaultLogger.Debug(args...)
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