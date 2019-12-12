//
//I0210 10:00:37.291011 info                                      tags=logs c=hello [example.go:main:62]
//{"@timestamp":"2017-02-10 10:02:43","@version":1,"c":"hello","tags":"logs","type":"hlog"}
//
//I0210 10:00:37.335020 info 123                                  a=b c=hello [example.go:main:63]
//W0210 10:00:37.335020 [1441090359-781398132]                    key=value [example.go:main:67]
//I0210 10:00:37.335020 1LP431GMDNa                               key=value [example.go:main:70]
//I0210 10:00:37.335020 nil id                                    [example.go:main:72]
//I0210 10:00:37.335020 nilid---                                  [example.go:main:75]
//I0210 10:00:37.335020 id is string                              [example.go:main:76]
//I0210 10:00:37.335020 test IdIface                              [123456] n=3 [example.go:main:80]
//I0210 10:00:37.335020 test IdIface asdf   asdf                  [123456] n=3 [example.go:main:81]
//I0210 10:00:37.335020 test IdIface                              [123456] N=3 [example.go:main:83]
//I0210 10:00:37.335020 test IdIface asdf asdf                    [123456] N=3 [example.go:main:84]
//I0210 10:00:37.335020 withiD                                    [123456] [example.go:main:86]
//W0210 10:00:37.335020 test IdIface asdf   asdf                  [123456] [example.go:main:87]
//I0210 10:00:37.335020 depth                                     [proc.go:main:183]
//I0210 10:00:37.335020 depth                                     [123456] [proc.go:main:183]
//I0210 10:00:37.335020 depth                                     [123456] [proc.go:main:183]
//I0210 10:00:37.335020 depth                                     key=value [proc.go:main:183]
//I0210 10:00:37.335020 --custome 1 head--test1test:2             [example.go:Log:23]
//I0210 10:00:37.335020 --custome 1 head--test1test:2             [example.go:main:94]
//I0210 10:00:37.335020 --custome 2 head--test1test:2             [proc.go:main:183]
//I0210 10:00:37.335020 ok ************************************** [example.go:main:96]
//I0210 10:00:37.335020 create merge task 1234 ---- process: 0% - [example.go:main:97]
//I0210 10:00:37.335020 create merge task 1234 ---- process: 0% - [example.go:main:98]
//I0210 10:00:38.335220 sleep                                     n=0 [example.go:main:102]
//I0210 10:00:39.335420 sleep                                     n=1 [example.go:main:102]
//reload
//I0210 10:00:40.335620 sleep                                     n=2 [example.go:main:102]
//I0210 10:00:40.335620  -------------------------- process: 0% - [example.go:main:108]
//W0210 10:00:40.335620  -------------------------- process: 0% - [entry.go:PadWarning:212]
//I0210 10:00:41.335820 ss ------------------------ process: 0% - [example.go:main:108]
//W0210 10:00:41.335820 ss ------------------------ process: 0% - [entry.go:PadWarning:212]
//I0210 10:00:42.336020 ssss ---------------------- process: 0% - [example.go:main:108]
//W0210 10:00:42.336020 ssss ---------------------- process: 0% - [entry.go:PadWarning:212]
//I0210 10:00:43.336220 ssssss -------------------- process: 0% - [example.go:main:108]
//W0210 10:00:43.336220 ssssss -------------------- process: 0% - [entry.go:PadWarning:212]
//I0210 10:00:44.336420 ssssssss ------------------ process: 0% - [example.go:main:108]
//W0210 10:00:44.336420 ssssssss ------------------ process: 0% - [entry.go:PadWarning:212]
//I0210 10:00:45.336620 --f1--                                    [example.go:f1:122]
//I0210 10:00:45.336620 --f1_1--                                  [example.go:f1_1:128]
//I0210 10:00:45.336620 long func                                 [example.go:longlonglonglonglonglong:134]
//F0210 10:00:45.336620 test depth                                [example.go:D:148]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/winxxp/hlog"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments] \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func Log(args ...interface{}) {
	hlog.PaddingColumns = 90
	hlog.Info("--custome 1 head--" + fmt.Sprint(args...))
	hlog.InfoDepth(1, "--custome 1 head--"+fmt.Sprint(args...))
	hlog.InfoDepth(2, "--custome 2 head--"+fmt.Sprint(args...))
}

type name int8

func (n name) ID() string {
	return "123456"
}

type NilId string

func (n NilId) String() string {
	return string(n) + "---"
}

func main() {
	defer func() {
		hlog.Info("-----", "defer 0")
	}()
	defer func() {
		if x := recover(); x != nil {
			hlog.WithFields(hlog.Fields{"stack": stacks(false)}).Error("panic:" + fmt.Sprint(x))
		}
	}()
	defer func() {
		hlog.Info("-----", "defer 2")
		defer hlog.Flush()
	}()

	flag.Parse()
	//hlog.Reload()
	hook, err := NewLogstashHook()
	if err != nil {
		hlog.Fatal("err", err)
	}
	hlog.AddHook(hook)

	hlog.WithField("tags", "logs").WithField("c", "hello").Info("info")
	hlog.WithField("a", "b").WithField("c", "hello").Infof("info %v", 123)

	hlog.WithFields(hlog.Fields{
		"key": "value",
	}).Warning("[1441090359-781398132]")

	hlog.WithID(nil).Info("nil id")
	hlog.WithID("string id").Info("string id test")

	nid := NilId("nilid")
	hlog.WithID(nid).Info(nid)
	hlog.WithID("1233466676").Info("id is string")

	n := name(3)

	hlog.WithField("n", n).WithID(n).Info("test IdIface")
	hlog.WithID(n).WithField("n", n).Info("test IdIface asdf   asdf")

	hlog.WithFields(hlog.Fields{"N": n}).WithID(n).Info("test IdIface")
	hlog.WithID(n).WithFields(hlog.Fields{"N": n}).Info("test IdIface asdf asdf")

	hlog.WithID(n).Info("withiD")
	hlog.WithID(n).Warning("test IdIface asdf   asdf")

	hlog.Depth(1).Info("depth")
	hlog.WithID(n).Depth(1).Info("depth")
	hlog.Depth(1).WithID(n).Info("depth")
	hlog.WithField("key", "value").Depth(1).Info("depth")

	Log("test", 1, "test:", 2.0)

	hlog.Padding('*').Info("ok")
	hlog.Padding('*').PadInfo("create merge task 1234", "process: 0%", '-')
	hlog.PadInfo("create merge task 1234", "process: 0%", '-')

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		hlog.WithField("n", i).Info("sleep")
	}

	hlog.WithResult(nil).Info("result test")
	hlog.WithResult(nil).Log("result test ok")
	hlog.WithResult(io.EOF).Log("result test error")

	for i := 0; i < 10; i += 2 {
		hlog.PadInfo(strings.Repeat("s", i), "process: "+strconv.Itoa(i*10)+"%", ' ')
		hlog.PadWarning(strings.Repeat("s", i), "process: 0%", '-')

		time.Sleep(time.Second)
	}

	f1()
	longlonglonglonglonglong()

	PanicTest()

}

func f1() {
	hlog.Info("--f1--")
	f1_1()

}

func f1_1() {
	hlog.Info("--f1_1--")
	hlog.V(2).Infoln("f1_1", 3, "f1_1")

}

func longlonglonglonglonglong() {
	hlog.Info("long func")
}

func PanicTest() {
	B()
}

func B() {
	C()
}
func C() {
	D()
}
func D() {
	hlog.Fatal("test depth")
}

func stacks(all bool) string {
	n := 10000
	if all {
		n = 100000
	}
	var trace []byte
	for i := 0; i < 5; i++ {
		trace = make([]byte, n)
		nbytes := runtime.Stack(trace, all)
		if nbytes < len(trace) {
			return string(trace[:nbytes])
		}
		n *= 2
	}
	return string(trace)
}

// LogstashHook to send logs to elastic.
type LogstashHook struct {
}

func NewLogstashHook() (*LogstashHook, error) {
	return &LogstashHook{}, nil
}

func (hook *LogstashHook) Fire(entry *hlog.Entry) error {
	if _, ok := entry.Data["tags"]; ok {
		data, err := (&LogstashFormatter{"hlog", "2006-01-02 15:04:05"}).Format(entry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
			return err
		}
		fmt.Fprintln(os.Stderr, string(data))
	}
	return nil
}

func (hook *LogstashHook) Severitys() []int {
	return []int{
		int(hlog.HookInfoLog),
		int(hlog.HookWarningLog),
		int(hlog.HookErrorLog),
		int(hlog.HookFatalLog),
	}

}

// Formatter generates json in logstash format.
type LogstashFormatter struct {
	Type string // if not empty use for logstash type field.

	// TimestampFormat sets the format used for timestamps.
	TimestampFormat string
}

func (f *LogstashFormatter) Format(entry *hlog.Entry) ([]byte, error) {
	fields := make(hlog.Fields)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/log/issues/137
			fields[k] = v.Error()
		default:
			fields[k] = v
		}
	}

	fields["@version"] = 1

	if f.TimestampFormat == "" {
		f.TimestampFormat = time.RFC3339
	}

	fields["@timestamp"] = time.Now().Format(f.TimestampFormat)

	// set message field
	v, ok := entry.Data["message"]
	if ok {
		fields["fields.message"] = v
	}

	// set level field
	v, ok = entry.Data["level"]
	if ok {
		fields["fields.level"] = v
	}

	// set type field
	if f.Type != "" {
		v, ok = entry.Data["type"]
		if ok {
			fields["fields.type"] = v
		}
		fields["type"] = f.Type
	}

	serialized, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
