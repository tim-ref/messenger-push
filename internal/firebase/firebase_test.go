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

package firebase

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/logger"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/models"
	"testing"
)

func TestSetupFirebase(t *testing.T) {
	_, _, _, err := SetupFirebase()
	if err != nil {
		t.Errorf("SetupFirebase() error = %v", err)
	}
}

var requestBody = models.RequestBody{
	Notification: models.Notification{
		Devices: []models.Device{
			{
				AppId:     "de.akquinet.timref.messenger-client",
				PushKey:   "MySecretkey",
				PushKeyTS: 12345678,
			},
		},
	},
}

var testconfig = models.Config{
	Server: models.Server{
		Host: "0.0.0.0",
		Port: "8080",
	},
	Client: []string{
		"de.akquinet.timref.messengerclient.data_message",
		"de.akquinet.timref.messenger-client",
	},
}

func TestSendToToken(t *testing.T) {
	logger.InitializeLogger("DEBUG")
	ClientFirebase = &MockMessagingClient{}
	mResponse, err := SendToToken(requestBody, testconfig)
	assert.NoError(t, err)
	fmt.Println(mResponse)
}
