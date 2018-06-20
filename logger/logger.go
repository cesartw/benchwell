package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger logs
type Logger interface {
	WithField(string, interface{}) Logger
	SetComponent(string) Logger
	Debug(args ...interface{})
	Debugf(f string, args ...interface{})
	Info(args ...interface{})
	Infof(f string, args ...interface{})
}

type logformatter struct{}

func (f *logformatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Time
	entry.Buffer.WriteString(fmt.Sprintf("%d ", entry.Time.Unix()))

	comp, ok := entry.Data["component"]
	if !ok {
		comp = "UNKNOWN"
	}

	// Component
	entry.Buffer.WriteString(" " + comp.(string) + " ")

	// Level
	entry.Buffer.WriteString(strings.ToUpper(entry.Level.String()) + ":")

	// Data
	for k, v := range entry.Data {
		if k == "component" {
			continue
		}
		entry.Buffer.WriteString(fmt.Sprintf(` %s="%v"`, k, v))
	}

	// Message
	entry.Buffer.WriteString(fmt.Sprintf(` "%s"`, entry.Message))
	entry.Buffer.WriteByte('\n')

	// buff.UnreadByte() // remove last space?

	return entry.Buffer.Bytes(), nil
}

// NewLogger ...
func NewLogger() Logger {
	f, err := os.OpenFile("debug.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetOutput(f)
	l.SetFormatter(&logformatter{})
	l.Info("Starting")

	return &logger{Entry: logrus.NewEntry(l)}
}

type logger struct {
	*logrus.Entry
}

func (l *logger) SetComponent(c string) Logger {
	return &logger{Entry: l.Entry.WithField("component", c)}
}

func (l *logger) WithField(f string, v interface{}) Logger {
	return &logger{Entry: l.Entry.WithField(f, v)}
}

func (l *logger) Debug(args ...interface{}) {
	l.Entry.Debug(args...)
}
