package logging

import (
	"fmt"
	"io"
	"os"
)

type Level int

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case ERROR:
		return "ERROR"
	case CRITICAL:
		return "CRITICAL"
	}
	return ""
}

const (
	DEBUG Level = iota
	INFO
	ERROR
	CRITICAL
)

type Event interface {
	Msg(s string) (n int, err error)
}

type e struct {
	currentLevel,
	level Level
	stdout io.Writer
	errout io.Writer
}

func (e e) Msg(s string) (n int, err error) {
	new := fmt.Sprintf("%s: %s\n", e.level.String(), s)
	if e.level >= e.currentLevel {
		return e.stdout.Write([]byte(new))
	}

	return 0, nil
}

type I interface {
	Debug() Event
	Info() Event
}

type i struct {
	lvl Level
	w   io.Writer
}

type Options struct {
	Level  Level
	Writer io.Writer
}

func New(opts ...Options) I {
	opt := Options{
		Level:  INFO,
		Writer: os.Stdout,
	}
	if len(opts) == 1 {
		opt = opts[0]
	}
	return i{lvl: opt.Level, w: opt.Writer}
}

func (i i) Debug() Event {
	return e{
		currentLevel: i.lvl,
		level:        DEBUG,
		stdout:       os.Stdout,
	}
}

func (i i) Info() Event {
	return e{
		currentLevel: i.lvl,
		level:        INFO,
		stdout:       os.Stdout,
	}
}
