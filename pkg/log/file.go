package log

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	_defaultLogFileNameTimeFormat = "20060102"
	_defaultFileMaxAge            = 7 * 24 * time.Hour
)

type fileHook struct {
	logPath string
	appName string

	levels   []logrus.Level
	f        *os.File
	lock     sync.Mutex
	lastTime time.Time
}

func NewFilekHook(level logrus.Level, path, appName string) *fileHook {
	hook := &fileHook{
		logPath: path,
		appName: appName,
	}

	for _, l := range logrus.AllLevels {
		if l <= level {
			hook.levels = append(hook.levels, l)
		}
	}

	return hook
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
	_, e := b.WriteString(fmt.Sprintf("%s[%s] %s\n", levelText, entry.Time.Format(_defaultTimeFormat), entry.Message))
	if e != nil {
		return nil, e
	}
	return b.Bytes(), nil
}

func (h *fileHook) getFile() {
	if h.f != nil {
		_ = h.f.Close()
	}

	if e := os.MkdirAll(h.logPath, os.ModePerm); e != nil {
		panic(e)
	}

	fileName := fmt.Sprintf("%s/%s-%s.log", h.logPath, h.appName, time.Now().Format(_defaultLogFileNameTimeFormat))
	if _, err := os.Stat(fileName); err == nil {
		f, err := os.OpenFile(filepath.Clean(fileName), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic(err)
		}
		h.f = f
	} else {
		f, err := os.Create(filepath.Clean(fileName))
		if err != nil {
			panic(err)
		}
		h.f = f
	}

	// delete expire file
	h.deleteExpireFile()
}

func (h *fileHook) deleteExpireFile() {
	files, err := os.ReadDir(h.logPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			re := regexp.MustCompile(`(\d{8})`)
			match := re.FindString(file.Name())
			if match == "" {
				continue
			}

			t, err := time.ParseInLocation(_defaultLogFileNameTimeFormat, match, time.Local)
			if err != nil {
				continue
			}

			if time.Since(t) > _defaultFileMaxAge {
				if err := os.Remove(filepath.Join(h.logPath, file.Name())); err != nil {
					panic(err)
				}
			}
		}
	}
}
