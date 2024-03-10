package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/deadlinefen/tinyPortMapper-manager-ipv6/pkg/config"
	filename "github.com/keepeye/logrus-filename"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/t-tomalak/logrus-prefixed-formatter"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLog(cfg *config.Log) {
	level, err := log.ParseLevel(cfg.Level)
	if err != nil {
		fmt.Printf("parse log level %v failed: %v\n", cfg.Level, err)
		os.Exit(1)
	}
	log.SetLevel(level)
	log.SetFormatter(&prefixed.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceFormatting: true,
	})

	if level > log.InfoLevel {
		// log.SetReportCaller(true), caller by filename
		filenameHook := filename.NewHook()
		filenameHook.Field = "line"
		log.AddHook(filenameHook)
	}

	// TODO: Add write permission check
	if cfg.ToStdoutOnly {
		log.SetOutput(os.Stdout)
	} else {
		output := &lumberjack.Logger{
			Filename:   cfg.Path,
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 30,
			LocalTime:  true,
			Compress:   false,
		}
		if cfg.AlsoToStderr {
			writer := io.MultiWriter(output, os.Stderr)
			log.SetOutput(writer)
		} else {
			log.SetOutput(output)
		}
	}
}
