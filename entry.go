package glog

import "fmt"

type IdIface interface {
	ID() string
}

type Fields map[string]interface{}

type Entry struct {
	Logger *loggingT
	Data   Fields
	Id     IdIface
}

func NewEntry(logger *loggingT) *Entry {
	return &Entry{
		Logger: logger,
		Data:   make(Fields, 5),
	}
}

// Add a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Fields{key: value})
}

// Add a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	data := Fields{}
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: data, Id: entry.Id}
}

func (entry *Entry) WithId(id IdIface) *Entry {
	entry.Id = id
	return entry
}

func (entry *Entry) log(s severity, args ...interface{}) {
	buf, file, fn, line := entry.Logger.header(s, 0)

	if entry.Id != nil {
		fmt.Fprint(buf, entry.Id.ID())
	}
	fmt.Fprint(buf, args...)
	buf.fillPading()

	for k, v := range entry.Data {
		buf.WriteString(fmt.Sprintf(" %s=%v", k, v))
	}

	fmt.Fprintf(buf, " %s:%s:%d", file, fn, line)

	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	entry.Logger.output(s, buf, file, line, false)
}

func (entry *Entry) Info(args ...interface{}) {
	entry.log(infoLog, args...)
}

func (entry *Entry) Warning(args ...interface{}) {
	entry.log(warningLog, args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.log(errorLog, args...)
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.log(fatalLog, args...)
}

func WithId(id IdIface) *Entry {
	return NewEntry(&logging).WithId(id)
}

func WithField(key string, value interface{}) *Entry {
	return NewEntry(&logging).WithField(key, value)
}

func WithFields(fields Fields) *Entry {
	return NewEntry(&logging).WithFields(fields)
}
