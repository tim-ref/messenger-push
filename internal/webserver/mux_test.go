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

package webserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/firebase"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/logger"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/models"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MyMockedObject struct {
	mock.Mock
}

var testconfig = WebserverConfig{
	Config: models.Config{
		Server: models.Server{
			Host: "0.0.0.0",
			Port: "8080",
		},
		Client: []string{
			"de.timref.messenger-client",
		},
	},
}

func SetupTestserver() (*httptest.Server, *zap.Logger) {
	logger.InitializeLogger("DEBUG")
	logr := logger.Logger

	testConfig := testconfig

	testwebserver := httptest.NewServer(NewRouter(&testConfig))
	return testwebserver, logr
}

func TestStart(t *testing.T) {
	logger.InitializeLogger("Info")
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		Start(testconfig.Config)
	}(ctx)
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
}

func TestNewRouter(t *testing.T) {
	routes := []string{"/_matrix/push/v1/notify", "/healthz", "/metrics"}
	router := NewRouter(&testconfig)
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		for _, v := range routes {
			if v == tpl {
				return nil
			}
		}
		err1 := fmt.Errorf("Error in TestNewRouter")
		return err1
	})
	require.NoError(t, err)
}
func TestMatrixPush(t *testing.T) {
	testwebserver, _ := SetupTestserver()
	defer testwebserver.Close()
	firebase.ClientFirebase = &firebase.MockMessagingClient{}

	requestBody := models.RequestBody{
		Notification: models.Notification{
			Devices: []models.Device{
				{
					AppId:     "de.timref.messenger-client",
					PushKey:   "someRandomKey",
					PushKeyTS: 12345678,
				},
			},
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	client := &http.Client{}
	req, err := http.NewRequest("POST", testwebserver.URL+"/_matrix/push/v1/notify", nil)
	require.NoError(t, err)
	req.Body = io.NopCloser(bytes.NewReader(body))
	response, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	newBody, err := io.ReadAll(response.Body)
	fmt.Println(newBody)
	assert.NoError(t, err)
	assert.Equal(t, "{}", string(newBody))
}

func TestMatrixNotFound(t *testing.T) {
	testwebserver, _ := SetupTestserver()
	defer testwebserver.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", testwebserver.URL+"/", nil)
	response, err := client.Do(req)
	require.NoError(t, err)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, "M_UNRECOGNIZED", string(body))
}

func TestMatrixNotAllowed(t *testing.T) {
	testwebserver, _ := SetupTestserver()
	defer testwebserver.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", testwebserver.URL+"/_matrix/push/v1/notify", nil)
	response, err := client.Do(req)
	require.NoError(t, err)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, response.StatusCode)
	assert.Equal(t, "M_UNRECOGNIZED", string(body))
}

func TestHealthHandler(t *testing.T) {
	testwebserver, _ := SetupTestserver()
	defer testwebserver.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", testwebserver.URL+"/healthz", nil)
	response, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "Server is healthy", string(resp))
}

func TestServerIsHealthy(t *testing.T) {
	resp := ServerIsHealthy(&testconfig)
	require.Equal(t, true, resp)
}

func TestMetrics(t *testing.T) {
	testwebserver, _ := SetupTestserver()
	defer testwebserver.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", testwebserver.URL+"/metrics", nil)
	response, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	require.NoError(t, err)
}
