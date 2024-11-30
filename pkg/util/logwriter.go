package util

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const (
	defaultFileMode  = 0666
	logRetentionDays = 90 // 保留日志的天数
	dateFormat       = "2006-01-02"
)

var (
	onceLogWriter sync.Once
	logWriter     *LogWriter
)

type LogWriter struct {
	filename string
	file     *os.File
	logName  string
	mu       sync.Mutex // 更清晰的互斥锁命名
}

func RedirectLog(logfile, crashfile string, mode os.FileMode) error {
	if err := os.MkdirAll(path.Dir(logfile), mode); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	if err := os.MkdirAll(path.Dir(crashfile), mode); err != nil {
		return fmt.Errorf("failed to create crash directory: %w", err)
	}

	absLogfile, err := filepath.Abs(logfile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if file, err := os.OpenFile(crashfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm); err == nil {
		syscall.Dup2(int(file.Fd()), 2)
		defer file.Close()
	}

	log.SetOutput(NewLogWriter(absLogfile))
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate | log.Lmicroseconds)
	return nil
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	dir := path.Dir(lw.filename)
	tm := time.Now()
	newLogFile := path.Join(dir, tm.Format(dateFormat)+".bak.log")

	lw.mu.Lock()
	defer lw.mu.Unlock()

	if newLogFile != lw.logName {
		if err := lw.rotateLog(newLogFile); err != nil {
			fmt.Printf("Failed to rotate log: %v\n", err)
		}
	}

	// 清理旧日志
	go lw.cleanOldLogs(dir, tm)

	if lw.file != nil {
		return lw.file.Write(p)
	}
	return 0, nil
}

func (lw *LogWriter) rotateLog(newLogFile string) error {
	f, err := os.OpenFile(newLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, defaultFileMode)
	if err != nil {
		return err
	}

	if _, err := f.WriteString("\n\n\n"); err != nil {
		f.Close()
		return err
	}

	if lw.file != nil {
		lw.file.Close()
	}

	lw.logName = newLogFile
	lw.file = f

	// 更新符号链接
	if err := os.Remove(lw.filename); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.Symlink(newLogFile, lw.filename)
}

func (lw *LogWriter) cleanOldLogs(dir string, now time.Time) {
	for i := logRetentionDays; i < logRetentionDays*2; i++ {
		oldTime := now.Add(-time.Hour * 24 * time.Duration(i))
		oldLogName := path.Join(dir, oldTime.Format(dateFormat)+".bak.log")
		if err := os.Remove(oldLogName); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Failed to remove old log %s: %v\n", oldLogName, err)
		}
	}
}

func NewLogWriter(logfile string) *LogWriter {
	onceLogWriter.Do(func() {
		logWriter = &LogWriter{
			filename: logfile,
		}
	})
	return logWriter
}
