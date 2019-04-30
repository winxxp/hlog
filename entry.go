package glog

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
	"unicode/utf8"
)

var (
	ErrTag        = "Err"
	ResultTag     = "result"
	ResultSucceed = "succeed"
	ResultFailed  = "failed"
)

type IdIface interface {
	ID() string
}

type ID string

func (id ID) ID() string {
	return string(id)
}

type Fields map[string]interface{}

type sortField struct {
	key string
	val interface{}
}

type sortFields []sortField

type Entry struct {
	Logger   *loggingT
	Data     Fields
	Depth_   int
	Padding_ byte
	Id       interface{}
}

func NewEntry(logger *loggingT) *Entry {
	return &Entry{
		Logger:   logger,
		Data:     make(Fields, 5),
		Depth_:   0,
		Padding_: ' ',
	}
}

func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Fields{key: value})
}

func (entry *Entry) WithID(id interface{}) *Entry {
	entry.Id = id
	return entry
}

func (entry *Entry) Depth(depth int) *Entry {
	entry.Depth_ = depth
	return entry
}

func (entry *Entry) Padding(padding byte) *Entry {
	entry.Padding_ = padding
	return entry
}

func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField(ErrTag, err)
}

func (entry *Entry) WithResult(err error) *Entry {
	result := ResultSucceed
	if err != nil {
		result = ResultFailed
		return entry.WithField(ResultTag, result).WithError(err)
	} else {
		return entry.WithField(ResultTag, result)
	}
}

// Add a map of fields to the Entry.
func (entry *Entry) WithFields(fields map[string]interface{}) *Entry {
	data := Fields{}
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return &Entry{Logger: entry.Logger,
		Data:     data,
		Depth_:   entry.Depth_,
		Padding_: ' ',
		Id:       entry.Id,
	}
}

func (entry *Entry) logf(s severity, format string, args ...interface{}) {
	buf, file, fn, line := entry.Logger.header(s, entry.Depth_)
	if format != "" {
		fmt.Fprintf(buf, format, args...)
	} else {
		fmt.Fprint(buf, args...)
	}

	buf.fillPading(entry.Padding_)

	switch id := entry.Id.(type) {
	case IdIface:
		buf.WriteByte('[')
		buf.WriteString(id.ID())
		buf.WriteByte(']')
		buf.Write(spacePad[:1])
	case string:
		buf.WriteByte('[')
		buf.WriteString(id)
		buf.WriteByte(']')
		buf.Write(spacePad[:1])
	}

	if result, found := entry.Data[ResultTag]; found {
		delete(entry.Data, ResultTag)
		buf.WriteByte('{')
		buf.WriteString(result.(string))
		buf.WriteByte('}')
		buf.Write(spacePad[:1])
	}
	if err, found := entry.Data[ErrTag]; found && err != nil {
		delete(entry.Data, ErrTag)
		buf.WriteByte('<')
		buf.WriteString(err.(error).Error())
		buf.WriteByte('>')
		buf.Write(spacePad[:1])
	}

	sf := make(sortFields, 0, len(entry.Data))
	for k, v := range entry.Data {
		sf = append(sf, sortField{k, v})
	}
	sort.Slice(sf, func(i, j int) bool {
		return strings.Compare(sf[i].key, sf[j].key) == -1
	})
	for _, v := range sf {
		fmt.Fprintf(buf, "%s=%v ", v.key, v.val)
	}

	fmt.Fprintf(buf, "[%s:%s:%d]", file, fn, line)

	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}

	entry.Logger.output(s, buf, file, line, false)
	entry.Logger.Hooks.Fire(int(s), entry)
}

func (entry *Entry) Log(args ...interface{}) {
	l := infoLog
	if err, found := entry.Data[ErrTag]; found && err != nil {
		l = errorLog
	} else {
		delete(entry.Data, ErrTag)
	}
	entry.logf(l, "", args...)
}

func (entry *Entry) Logf(format string, args ...interface{}) {
	l := infoLog
	if err, found := entry.Data[ErrTag]; found && err != nil {
		l = errorLog
	} else {
		delete(entry.Data, ErrTag)
	}
	entry.logf(l, format, args...)
}

func (entry *Entry) Info(args ...interface{}) {
	entry.logf(infoLog, "", args...)
}

func (entry *Entry) Warning(args ...interface{}) {
	entry.logf(warningLog, "", args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.logf(errorLog, "", args...)
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.logf(fatalLog, "", args...)
}

func (entry *Entry) Exit(args ...interface{}) {
	atomic.StoreUint32(&fatalNoStacks, 1)
	entry.logf(fatalLog, "", args...)
}

func (entry *Entry) padLog(s severity, ls, rs string, pad byte) {
	entry.Depth_ = entry.Depth_ + 1
	entry.logf(s, "", CreatPadInfo(ls, rs, pad, PaddingColumns))
}

func (entry *Entry) PadInfo(ls, rs string, pad byte) {
	entry.padLog(infoLog, ls, rs, pad)
}

func (entry *Entry) PadWarning(ls, rs string, pad byte) {
	entry.padLog(warningLog, ls, rs, pad)
}

func (entry *Entry) PadError(ls, rs string, pad byte) {
	entry.padLog(errorLog, ls, rs, pad)
}

func (entry *Entry) PadFatal(ls, rs string, pad byte) {
	entry.padLog(fatalLog, ls, rs, pad)
}

func (entry *Entry) PadExit(ls, rs string, pad byte) {
	atomic.StoreUint32(&fatalNoStacks, 1)
	entry.padLog(fatalLog, ls, rs, pad)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.logf(infoLog, format, args...)
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.logf(warningLog, format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.logf(errorLog, format, args...)
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.logf(fatalLog, format, args...)
}

func (entry *Entry) Exitf(format string, args ...interface{}) {
	atomic.StoreUint32(&fatalNoStacks, 1)
	entry.logf(fatalLog, format, args...)
}

func WithID(id interface{}) *Entry {
	return NewEntry(&logging).WithID(id)
}

func WithIDString(id string) *Entry {
	return NewEntry(&logging).WithID(ID(id))
}

func WithError(err error) *Entry {
	return NewEntry(&logging).WithField("Err", err)
}

func WithResult(err error) *Entry {
	return NewEntry(&logging).WithResult(err)
}

func WithField(key string, value interface{}) *Entry {
	return NewEntry(&logging).WithField(key, value)
}

func WithFields(fields Fields) *Entry {
	return NewEntry(&logging).WithFields(fields)
}

func Depth(depth int) *Entry {
	return NewEntry(&logging).Depth(depth)
}

func Padding(padding byte) *Entry {
	return NewEntry(&logging).Padding(padding)
}

func PadInfo(ls, rs string, pad byte) {
	NewEntry(&logging).Depth(1).PadInfo(ls, rs, pad)
}

func PadWarning(ls, rs string, pad byte) {
	NewEntry(&logging).PadWarning(ls, rs, pad)
}

func PadError(ls, rs string, pad byte) {
	NewEntry(&logging).PadError(ls, rs, pad)
}

func PadFatal(ls, rs string, pad byte) {
	NewEntry(&logging).PadFatal(ls, rs, pad)
}

func PadExit(ls, rs string, pad byte) {
	NewEntry(&logging).PadExit(ls, rs, pad)
}

var (
	ansi = regexp.MustCompile("\033\\[(?:[0-9]{1,3}(?:;[0-9]{1,3})*)?[m|K]")
)

func DisplayWidth(str string) int {
	return utf8.RuneCountInString(ansi.ReplaceAllLiteralString(str, ""))
}

func CreatPadInfo(ls, rs string, pad byte, width int) string {
	gap := width - 27 - DisplayWidth(ls) - DisplayWidth(rs)

	buf := bytes.NewBufferString(ls)
	buf.WriteByte(' ')

	if gap > 0 {
		buf.Write(bytes.Repeat([]byte{pad}, gap))
	} else {
		buf.WriteByte(pad)
	}

	buf.WriteByte(' ')
	buf.WriteString(rs)
	buf.WriteByte(' ')
	buf.WriteByte(pad)

	return buf.String()
}
