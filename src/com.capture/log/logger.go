package log

import (
	"io"
	"log"
	"os"
)

type LogConf struct {
	FileName *string
}

func (l *LogConf) InitLog() {
	if *l.FileName == "" {
		return
	}
	//初始化日志
	errFile, err := os.OpenFile(*l.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(io.MultiWriter(errFile))
	}
}
