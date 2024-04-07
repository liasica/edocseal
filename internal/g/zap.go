// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-07, by liasica

package g

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapRedisWriter struct {
	cli *redis.Client
	key string
}

func (w *zapRedisWriter) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(context.Background(), w.key, p).Result()
	return int(n), err
}

// 日志写入Redis
func zapRedisCore() zapcore.Core {
	jsonEnc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	syncer := zapcore.AddSync(&zapRedisWriter{
		cli: NewRedis(),
		key: cfg.Logger.RedisKey,
	})
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
	if cfg.Logger.Redis {
		cores = append(cores, zapRedisCore())
	}
	if cfg.Logger.Console {
		cores = append(cores, zapConsoleCore())
	}

	core := zapcore.NewTee(cores...)

	// logger 输出到 console 且标识调用代码行
	return zap.New(core).WithOptions(zap.AddCaller())
}
