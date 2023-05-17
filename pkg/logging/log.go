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
	case FATAL:
		return "FATAL"
	}
	return ""
}

const (
	DEBUG Level = iota
	INFO
	ERROR
	CRITICAL
	FATAL
)

type Event interface {
	Msg(s string) (n int, err error)
	Msgf(s string, a ...any) (n int, err error)
}

type e struct {
	currentLevel,
	level Level
	terminate bool
	stdout    io.Writer
	errout    io.Writer
}

func (e e) Msg(s string) (n int, err error) {
	new := fmt.Sprintf("%s: %s\n", e.level.String(), s)
	if e.level >= e.currentLevel {
		n, err = e.stdout.Write([]byte(new))
	}

	if e.terminate {
		os.Exit(1)
	}

	return n, err
}

func (e e) Msgf(s string, a ...any) (n int, err error) {
	new := fmt.Sprintf("%s: %s\n", e.level.String(), fmt.Sprintf(s, a...))
	if e.level >= e.currentLevel {
		n, err = e.stdout.Write([]byte(new))
	}

	if e.terminate {
		os.Exit(1)
	}

	return n, err
}

type I interface {
	Debug() Event
	Info() Event
	Error() Event
	Fatal() Event
	Writer() io.Writer
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
		Level: INFO,
	}

	if len(opts) == 1 {
		opt = opts[0]
	}

	if opt.Writer == nil {
		opt.Writer = os.Stdout
	}
	return i{lvl: opt.Level, w: opt.Writer}
}

func (i i) Writer() io.Writer {
	return i.w
}

func (i i) Debug() Event {
	return e{
		currentLevel: i.lvl,
		level:        DEBUG,
		stdout:       i.w,
	}
}

func (i i) Info() Event {
	return e{
		currentLevel: i.lvl,
		level:        INFO,
		stdout:       i.w,
	}
}

func (i i) Error() Event {
	return e{
		currentLevel: i.lvl,
		level:        ERROR,
		stdout:       i.w,
	}
}

func (i i) Fatal() Event {
	return e{
		currentLevel: i.lvl,
		level:        FATAL,
		stdout:       i.w,
		terminate:    true,
	}
}
