package application

import (
	"errors"
	"fmt"
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
	LogFilenameSuffix string
	ReportCaller bool

	LogFormatter logrus.Formatter
}

type CustomLogSettings struct {
	LoggingLevel logrus.Level
	LogDir string
	LogFilename string
	LogFilenameSuffix string
	ReportCaller bool

	LogFormatter logrus.Formatter
}

type AppSettings struct {
	AppLogSettings AppLogSettings
	CustomLogSettings CustomLogSettings
}

type Application struct {
	rootPath string
	stop chan struct{}
	commands map[string]CommandFunc
	logger *logger.Logger
	customLogger *logger.Logger
	Logger *logger.Logger
}

func (app *Application) Init(settings AppSettings) bool {
	app.rootPath = utils.GetExecutableRootPath()
	app.initAppLogger(settings.AppLogSettings)
	app.initCustomLogger(settings.CustomLogSettings)
	app.addDefaultCommands()
	app.debug("Application::Init: Initialized success.")

	return true
}

func (app *Application) initAppLogger(settings AppLogSettings) {
	if settings.DisableLogging {
		return
	}

	app.logger = logger.NewLogger(logger.LogSettings{
		Settings: settings.LogFormatter,
		LogDir: settings.LogDir,
		LogFilename: settings.LogFilename,
		LogFilenameSuffix: settings.LogFilenameSuffix,
		LogLevel: settings.LoggingLevel,
		ReportCaller: settings.ReportCaller,
	})
}

func (app *Application) initCustomLogger(settings CustomLogSettings) {
	app.customLogger = logger.NewLogger(logger.LogSettings{
		Settings: settings.LogFormatter,
		LogDir: settings.LogDir,
		LogFilename: settings.LogFilename,
		LogFilenameSuffix: settings.LogFilenameSuffix,
		LogLevel: settings.LoggingLevel,
		ReportCaller: settings.ReportCaller,
	})
	app.Logger = app.customLogger
}

func (app *Application) Run() {
	app.debug("Application::Run: App was started.")

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

func (app *Application) trace(args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Trace(args...)
}

func (app *Application) debug(args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Debug(args...)
}

func (app *Application) info(args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Info(args...)
}

func (app *Application) warn(args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Warn(args...)
}

func (app *Application) error(args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Error(args...)
}

func (app *Application) fatal(args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Fatal(args...)
}

func (app *Application) panic(args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Panic(args...)
}

func (app *Application) traceF(format string, args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Tracef(format, args...)
}

func (app *Application) debugF(format string, args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Debugf(format, args...)
}

func (app *Application) infoF(format string, args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Infof(format, args...)
}

func (app *Application) warnF(format string, args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Warnf(format, args...)
}

func (app *Application) errorF(format string, args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Errorf(format, args...)
}

func (app *Application) fatalF(format string, args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Fatalf(format, args...)
}

func (app *Application) panicF(format string, args ...interface{}) {
	if app.logger == nil {
		return
	}

	app.logger.Panicf(format, args...)
}

func (app *Application) Trace(args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Trace(args...)
}

func (app *Application) Debug(args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Debug(args...)
}

func (app *Application) Info(args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Info(args...)
}

func (app *Application) Warn(args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Warn(args...)
}

func (app *Application) Error(args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Error(args...)
}

func (app *Application) Fatal(args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Fatal(args...)
}

func (app *Application) Panic(args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Panic(args...)
}

func (app *Application) TraceF(format string, args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Tracef(format, args...)
}

func (app *Application) DebugF(format string, args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Debugf(format, args...)
}

func (app *Application) InfoF(format string, args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Infof(format, args...)
}

func (app *Application) WarnF(format string, args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Warnf(format, args...)
}

func (app *Application) ErrorF(format string, args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Errorf(format, args...)
}

func (app *Application) FatalF(format string, args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Fatalf(format, args...)
}

func (app *Application) PanicF(format string, args ...interface{}) {
	if app.customLogger == nil {
		return
	}

	app.customLogger.Panicf(format, args...)
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