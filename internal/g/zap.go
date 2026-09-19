// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-07, by liasica

package g

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 日志投递 Kafka
func zapKafkaCore() zapcore.Core {
	jsonEnc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	syncer := zapcore.AddSync(newKafkaWriter(cfg.Kafka.Addresses))
	return zapcore.NewCore(jsonEnc, syncer, zap.NewAtomicLevelAt(zap.DebugLevel))
}

// 日志写入控制台
func zapConsoleCore() zapcore.Core {
	consoleEnc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	return zapcore.NewCore(consoleEnc, zapcore.Lock(os.Stdout), zap.NewAtomicLevelAt(zap.DebugLevel))
}

func NewZap() *zap.Logger {
	// 集成多个 core
	var cores []zapcore.Core
	if cfg.Logger.Kafka {
		cores = append(cores, zapKafkaCore())
	}
	if cfg.Logger.Console {
		cores = append(cores, zapConsoleCore())
	}

	core := zapcore.NewTee(cores...)

	// logger 输出到 console 且标识调用代码行
	return zap.New(core).WithOptions(zap.AddCaller()).Named(cfg.Logger.LoggerName)
}
