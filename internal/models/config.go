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

package models

import (
	"time"
)

type Config struct {
	Server Server   `yaml:"Server"`
	Client []string `yaml:"Clients"`
}

type Server struct {
	Host    string  `yaml:"Host"`
	Port    string  `yaml:"Port"`
	Timeout Timeout `yaml:"Timeout"`
}

type Timeout struct {
	Server time.Duration `yaml:"Server"`
	Write  time.Duration `yaml:"Read"`
	Read   time.Duration `yaml:"Write"`
	Idle   time.Duration `yaml:"Idle"`
}
