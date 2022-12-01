package logger

import (
	"log"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(message string)
	Infoln(message string)
	Infof(template string , args ...interface{})
	Debug(message string)
	Error(message string)
	Fatal(message string)
	Fatalf(template string, args ...interface{})
	Fatalln(message string)
}

type ZapperLogger struct {
	logger *zap.Logger
}

func (zl *ZapperLogger) Info(message string) {
	zl.logger.Info(message)
}

func(zl *ZapperLogger) Infoln(message string){
	zl.logger.Sugar().Infoln(message)
}

func (zl *ZapperLogger) Infof(template string, args ...interface{}){
	zl.logger.Sugar().Infof(template , args...)
}

func (zl *ZapperLogger) Debug(message string) {
	zl.logger.Debug(message)
}

func (zl *ZapperLogger) Error(message string) {
	zl.logger.Error(message)
}

func (zl *ZapperLogger) Fatal(message string){
	zl.logger.Fatal(message)
}

func (z1 *ZapperLogger) Fatalf(template string, args ...interface{}){
	z1.logger.Sugar().Fatalf(template, args...)
}

func (zl *ZapperLogger) Fatalln(message string){
	zl.logger.Sugar().Fatalln(message)
}

func NewZapperLogger() *ZapperLogger {
	// zapper configurations
	zapConfig := zap.NewProductionConfig()

	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalln("Error while initializing Zap Logger: " + err.Error())
	}

	return &ZapperLogger{logger: logger}
}

var Module = fx.Module("logger", fx.Provide(fx.Annotate(NewZapperLogger, fx.As(new(Logger)))))
