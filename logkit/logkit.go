package logkit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	color_red = uint8(iota + 91)
	color_green
	color_yellow
	color_blue
	color_magenta //洋红
)

var (
	info *log.Logger
	warn *log.Logger
	succ *log.Logger
	erro *log.Logger

	//log文件
	infofile *log.Logger
	warnfile *log.Logger
	succfile *log.Logger
	errofile *log.Logger

	lock = new(sync.Mutex)
)

func Info(data ...interface{}) {
	infofile.Println(data...)
	info.Println(data...)
}
func InfoF(format string, data ...interface{}) {
	infofile.Printf(format, data...)
	info.Printf(format, data...)
}

func Err(data ...interface{}) {
	_, file, line, ok := runtime.Caller(0)
	if !ok {
		file = "???"
		line = 0
	}
	errofile.Println([]interface{}{file, line, data})
	erro.Println([]interface{}{file, line, data})
}
func ErrF(format string, data ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	errofile.Printf(file+":"+strconv.Itoa(line)+format, data)
	erro.Printf(file+":"+strconv.Itoa(line)+format, data)
}

func Warn(data ...interface{}) {
	warnfile.Println(data...)
	warn.Println(data...)
}

func Warnf(format string, data ...interface{}) {
	warnfile.Printf(format, data...)
	warn.Printf(format, data...)
}
func Succ(data ...interface{}) {
	succfile.Println(data...)
	succ.Println(data...)
}

func Succf(format string, data ...interface{}) {
	succfile.Printf(format, data...)
	succ.Printf(format, data...)
}
func LogInit(logdir string) {
	lastTime := time.Now()
	logStart(lastTime, logdir)
	go func() {
		for {
			select {
			//24小时重启一次
			case now := <-time.After(500 * time.Millisecond):
				if now.Unix()%24*3600 == 0 && lastTime != now {
					logStart(now, logdir)
					lastTime = now
				}
			}
		}
	}()
}
func logStart(now time.Time, logdir string) {
	lock.Lock()
	defer lock.Unlock()
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if logdir[:2] == "./" {
		logdir = logdir[1:]
	}
	tail := len(logdir) - 1
	if logdir[tail] != '/' {
		logdir += "/"
	}
	dir += logdir + now.Format("20060102") + ".log"
	errFile, err := os.OpenFile(dir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
		return
	}
	switch runtime.GOOS {
	case "windows":
		succ = log.New(os.Stdout, "[SUCCESS]", log.Ldate|log.Ltime)
		info = log.New(os.Stdout, "[INFO]", log.Ldate|log.Ltime)
		warn = log.New(os.Stdout, "[WARN]", log.Ldate|log.Ltime)
		erro = log.New(os.Stdout, "[ERR]", log.Ldate|log.Ltime)
	default:
		succ = log.New(os.Stdout, color(color_green, "[SUCCESS]"), log.Ldate|log.Ltime)
		info = log.New(os.Stdout, color(color_blue, "[INFO]"), log.Ldate|log.Ltime)
		warn = log.New(os.Stdout, color(color_yellow, "[WARN]"), log.Ldate|log.Ltime)
		erro = log.New(os.Stdout, color(color_red, "[ERR]"), log.Ldate|log.Ltime)
	}
	succfile = log.New(errFile, "[SUCCESS]", log.Ldate|log.Ltime)
	infofile = log.New(errFile, "[INFO]", log.Ldate|log.Ltime)
	warnfile = log.New(errFile, "[WARN]", log.Ldate|log.Ltime)
	errofile = log.New(errFile, "[ERR]", log.Ldate|log.Ltime)
}

func color(color uint8, s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, s)
}
