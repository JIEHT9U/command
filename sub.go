package command

import (
	"strings"

	"github.com/pkg/errors"
)

type Sub struct {
	commands []string
}

func NewSub(commands string) (*Sub, error) {
	var subCommands = strings.Split(commands, " ")
	if len(subCommands) == 0 {
		return nil, errors.New("error create new Sub expect minimum one command")
	}
	return &Sub{commands: subCommands}, nil
}

func (s *Sub) Next() (string, bool) {
	if s == nil {
		return "", false
	}
	if len(s.commands) == 0 {
		return "", false
	}
	next := s.commands[0]
	s.commands = s.commands[1:]
	return next, true
}

func (s *Sub) Clone() *Sub {
	copyCommands := make([]string, len(s.commands))
	copy(copyCommands, s.commands)
	return &Sub{
		commands: copyCommands,
	}
}
