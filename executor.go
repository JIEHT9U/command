package command

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/jieht9u/logger"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Executor struct {
	app      *kingpin.Application
	commands map[string]Command

	envFile string
	sub     *Sub
}

func New(name, help string) *Executor {
	var app = kingpin.New(name, help)
	app.HelpFlag.Short('h')
	return &Executor{
		app:      app,
		commands: make(map[string]Command),
		envFile:  ".env",
	}
}

func (e *Executor) AddFlag(name, help string) *kingpin.FlagClause {
	return e.app.Flag(name, help)
}

func (e *Executor) ExecCommandName() string {
	if name, ok := e.sub.Clone().Next(); ok {
		return name
	}
	return "dummy"
}

func (e *Executor) SetEnvFilePath(path string) *Executor {
	e.envFile = path
	return e
}

func (e *Executor) Register(name string, command Command, helpText string) *Executor {
	if e.isCommandExisting(name) {
		logger.Fatalf("error register command, command already exists: %v", name)
	}
	command.FlagParse(e.app.Command(name, helpText))
	e.commands[name] = command
	return e
}

func (e *Executor) FlagParse() error {

	if err := loadEnvFromDotFiles(e.envFile); err != nil {
		return err
	}

	commands, err := e.app.Parse(os.Args[1:])
	if err != nil {
		e.app.Usage(os.Args[1:])
		return errors.Wrap(err, "error parsing input params")
	}

	if e.sub, err = NewSub(commands); err != nil {
		return err
	}

	return nil
}

func (e *Executor) Exec(ctx context.Context) error {
	var execCommand, ok = e.sub.Next()
	if !ok {
		return errors.New("error exec unspecified command")
	}

	cmd, err := e.getCommand(execCommand)
	if err != nil {
		return err
	}
	return cmd.Exec(ctx, e.sub)
}

func loadEnvFromDotFiles(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		if err := godotenv.Load(filename); err != nil {
			return errors.Wrapf(err, "error loading %v file: %v", filename, err)
		}
	}
	return nil
}

func (e *Executor) getCommand(name string) (Command, error) {
	if c, find := e.commands[name]; find {
		return c, nil
	}
	return nil, errors.Errorf("error get cmd %s", name)
}

func (e *Executor) isCommandExisting(command string) bool {
	if _, ok := e.commands[command]; ok {
		return true
	}
	return false
}
