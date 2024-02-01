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

package config

import (
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/models"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	var config *models.Config
	var err error
	config, err = NewConfig("./config.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, models.Timeout{Server: 30000000000, Write: 15000000000, Read: 15000000000, Idle: 5000000000}, config.Server.Timeout)
	assert.Equal(t, []string{"de.test.timref.messengerclient"}, config.Client)

	_, err = NewConfig("./anyPath/config.yaml")
	assert.Error(t, err, "config, err = NewConfig(\"./config.yaml\")")

}

func TestValidateConfigPath(t *testing.T) {
	err := ValidateConfigPath("./config.yaml")
	assert.NoError(t, err)
	err = ValidateConfigPath("./noConfig.yaml")
	assert.EqualErrorf(t, err, "stat ./noConfig.yaml: no such file or directory", "'../config' is a directory, not a normal file")
	err = ValidateConfigPath("../config")
	assert.EqualErrorf(t, err, "'../config' is a directory, not a normal file", "'../config' is a directory, not a normal file")
}
func TestParseFlags(t *testing.T) {
	var flags string
	var err error

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	flag.CommandLine = flag.NewFlagSet("flags set", flag.ExitOnError)
	os.Args = append([]string{"flags set"}, "-config", "./config.yaml")
	flags, err = ParseFlags()
	assert.NoError(t, err)
	assert.Equal(t, "./config.yaml", flags)

	flag.CommandLine = flag.NewFlagSet("flags set", flag.ExitOnError)
	os.Args = append([]string{"flags set"}, "-config", "./config/config.yaml")
	flags, err = ParseFlags()
	fmt.Println(flags)
	assert.Error(t, err)
}
