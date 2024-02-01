/*
 * Copyright (C) 2023 akquinet GmbH
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

var Logger *zap.Logger

func InitializeLogger(logLevelEnv string) {
	var defaultLogLevel zapcore.Level

	loglevel, err := zapcore.ParseLevel(logLevelEnv)
	if err != nil {
		log.Fatal(err)
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	// This is not recommended, since we want json output
	// consoleEncoder := zapcore.NewConsoleEncoder(config)
	consoleEncoder := zapcore.NewJSONEncoder(config)

	if loglevel == 0 {
		defaultLogLevel = zapcore.InfoLevel
	} else {
		defaultLogLevel = loglevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
