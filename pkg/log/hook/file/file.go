// Package file package file
package file

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	_defaultTimeFormat            = time.RFC3339
	_defaultLogFileNameTimeFormat = time.DateOnly
	_defaultFileMaxAge            = 7 * 24 * time.Hour
)

var (
	singlton *fileHook
	once     sync.Once
)

type fileHook struct {
	logFolder string
	appName   string

	reportCaller     bool
	callerPrettyfier func(*runtime.Frame) (function string, file string)

	levels   []logrus.Level
	f        *os.File
	lock     sync.Mutex
	lastTime time.Time
}

func Get(level logrus.Level, path, appName string) *fileHook {
	if singlton != nil {
		return singlton
	}
	once.Do(func() {
		hook := &fileHook{
			logFolder: path,
			appName:   appName,
		}
		for _, l := range logrus.AllLevels {
			if l <= level {
				hook.levels = append(hook.levels, l)
			}
		}
		singlton = hook
	})
	return singlton
}

func (h *fileHook) SetReportCaller(reportCaller bool, f func(*runtime.Frame) (function string, file string)) {
	h.reportCaller = reportCaller
	h.callerPrettyfier = f
}

func (h *fileHook) Levels() []logrus.Level {
	return h.levels
}

func (h *fileHook) Fire(entry *logrus.Entry) error {
	defer h.lock.Unlock()
	h.lock.Lock()

	msg, err := h.Format(entry)
	if err != nil {
		return err
	}
	if h.lastTime.Day() != entry.Time.Day() || h.lastTime.IsZero() {
		h.getFile()
	}
	_, err = h.f.Write(msg)
	if err != nil {
		return err
	}
	h.lastTime = entry.Time
	return nil
}

func (h *fileHook) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	levelText := strings.ToUpper(entry.Level.String())[0:4]
	msg := fmt.Sprintf("%s[%s] %s\n", levelText, entry.Time.Format(_defaultTimeFormat), entry.Message)
	if h.reportCaller && h.callerPrettyfier != nil {
		caller, _ := h.callerPrettyfier(entry.Caller)
		msg = fmt.Sprintf("%s[%s]%s %s\n", levelText, entry.Time.Format(_defaultTimeFormat), caller, entry.Message)
	}
	_, e := b.WriteString(msg)
	if e != nil {
		return nil, e
	}
	return b.Bytes(), nil
}

func (h *fileHook) getFile() {
	if h.f != nil {
		_ = h.f.Close()
	}
	if e := os.MkdirAll(h.logFolder, 0o750); e != nil {
		panic(e)
	}
	path := fmt.Sprintf("%s/%s-%s.log", h.logFolder, h.appName, time.Now().Format(_defaultLogFileNameTimeFormat))
	if _, err := os.Stat(path); err == nil {
		f, err := os.OpenFile(filepath.Clean(path), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic(err)
		}
		h.f = f
	} else {
		f, err := os.Create(filepath.Clean(path))
		if err != nil {
			panic(err)
		}
		h.f = f
	}
	h.deleteExpireFile()
}

func (h *fileHook) deleteExpireFile() {
	files, err := os.ReadDir(h.logFolder)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
			match := re.FindString(file.Name())
			if match == "" {
				continue
			}
			t, err := time.ParseInLocation(_defaultLogFileNameTimeFormat, match, time.Local)
			if err != nil {
				continue
			}
			if time.Since(t) > _defaultFileMaxAge {
				if err := os.Remove(filepath.Join(h.logFolder, file.Name())); err != nil {
					panic(err)
				}
			}
		}
	}
}
