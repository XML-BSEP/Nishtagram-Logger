package Nishtagram_Logger

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"os"
	"path/filepath"
	"runtime"
	"strings"

)

type LoggerUseCase interface {
	InitializeLogger(serviceName string, ctx context.Context)
	CreateLogDir(filePath string) error
	CloseLogger()
	FormatFilePath(file string) string
}

type Logger struct {
	logger *logrus.Logger
	file *os.File
}

func (l *Logger) InitializeLogger(serviceName string, ctx context.Context) {
	if err := l.CreateLogDir(filepath.FromSlash("../log/logfiles/" + serviceName)); err != nil {
		logrus.Fatalf("Failed to create directory for log files | %v\n", err)
	}

	file := filepath.FromSlash("../log/logfiles/" + serviceName + "/" + serviceName + ".log")

	rotatingLogs, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename: file,
		MaxSize: 100,
		MaxBackups: 50,
		MaxAge: 14,
		Level: logrus.InfoLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: "01-01-2005 16:21:02",
			DataKey: "data",
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				return frame.Function, fmt.Sprintf("%s:%d", l.FormatFilePath(frame.File), frame.Line)
			}},
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize rotating hook | %v\n", err)
	}

	l.logger.AddHook(rotatingLogs)
	l.logger.SetReportCaller(true)
	l.logger.SetOutput(os.Stdout)
	l.logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "01-01-2005 16:21:02",
		FullTimestamp: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, fmt.Sprintf("%s:%d", l.FormatFilePath(frame.File), frame.Line)
		}})
	l.logger.SetLevel(logrus.InfoLevel)

}

func (l *Logger) FormatFilePath(file string) string {
	arr := strings.Split(filepath.ToSlash(file), "/")
	return arr[len(arr) - 1]
}

func (l *Logger) CreateLogDir(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return os.MkdirAll(filePath, os.ModeDir|0755)
	}
	return nil
}

func (l *Logger) CloseLogger() {
	l.file.Close()
}

func NewLogger() LoggerUseCase {
	return &Logger{}
}