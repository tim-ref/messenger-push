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

package main

import (
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/config"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/firebase"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/logger"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/webserver"
	"go.uber.org/zap"
	"os"
)

func main() {
	// LogLevelEnv get the Loglevel from Environment.
	LogLevelEnv := os.Getenv("MATRIX_GO_PUSH_LOGLEVEL")
	// Initializes the Logger
	logger.InitializeLogger(LogLevelEnv)
	logr := logger.Logger
	logr.Info("Setting LogLevel to " + LogLevelEnv)

	// Gets the Config and Errors, if the Config can't get parsed
	cfgPath, err := config.ParseFlags()
	if err != nil {
		logr.Fatal("Couldn't Parse Config ", zap.String("Path: ", cfgPath), zap.String(" Error: ", err.Error()))
	}

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		logr.Fatal(err.Error())
	}

	err = firebase.CreateFBCClient()
	if err != nil {
		logr.Fatal(err.Error())
	}

	webserver.Start(*cfg)
}
