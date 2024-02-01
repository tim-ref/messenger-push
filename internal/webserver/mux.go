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
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/firebase"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/logger"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/metrics"
	"gitlab.spree.de/akquinet/health/timref/matrix-go-push.git/internal/models"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"time"
)

func Start(mc models.Config) {
	c := &WebserverConfig{mc}
	srv := &http.Server{
		Handler:      NewRouter(c),
		Addr:         mc.Server.Host + ":" + mc.Server.Port,
		WriteTimeout: mc.Server.Timeout.Write * time.Second,
		ReadTimeout:  mc.Server.Timeout.Read * time.Second,
		IdleTimeout:  mc.Server.Timeout.Idle * time.Second,
	}

	logger.Logger.Info("Starting Webserver on: " + srv.Addr)
	logger.Logger.Info("Setting timeouts : ", zap.Duration(" WriteTimeout in seconds: ", mc.Server.Timeout.Write), zap.Duration(" ReadTimeout in seconds: ", mc.Server.Timeout.Read), zap.Duration(" IdleTimeout in seconds: ", mc.Server.Timeout.Idle))
	log.Fatal(srv.ListenAndServe())
}

func NewRouter(c *WebserverConfig) *mux.Router {
	router := mux.NewRouter()

	// Matrix Spec Endpoint https://spec.matrix.org/v1.7/push-gateway-api/#post_matrixpushv1notify
	router.HandleFunc("/_matrix/push/v1/notify", c.MatrixPush).Methods("POST")
	// Health Check
	router.HandleFunc("/healthz", c.HealthHandler).Methods("GET")
	// Prometheus endpoint
	router.Path("/metrics").Handler(promhttp.Handler())

	// https://spec.matrix.org/unstable/push-gateway-api/#unsupported-endpoints
	notFound := http.HandlerFunc(MatrixNotFound)
	router.NotFoundHandler = notFound
	notAllowed := http.HandlerFunc(MatrixNotAllowed)
	router.MethodNotAllowedHandler = notAllowed

	return router
}

// MatrixPush processes the actual push message and Sends the Message to the Firbase SendToToken handler
// https://spec.matrix.org/v1.7/push-gateway-api/#post_matrixpushv1notify
func (c *WebserverConfig) MatrixPush(w http.ResponseWriter, r *http.Request) {
	var notify models.RequestBody
	var rejected models.MatrixResponse
	var err error
	defer metrics.TotalPostNotifyRequests.Inc()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Logger.Error("Error reading Body from Request")
	}
	err = json.Unmarshal(body, &notify)
	if err != nil {
		logger.Logger.Error("Error Unmarshal Request")
		fmt.Println("M_UNRECOGNIZED")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	logger.Logger.Info("Sending Notify to Firebase")
	rejected, err = firebase.SendToToken(notify, c.Config)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("M_UNRECOGNIZED")
		return
	}
	response, _ := json.Marshal(rejected)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(response)
}

// MatrixNotFound returns a 404 Code with M_UNRECOGNIZED
// https://spec.matrix.org/v1.7/push-gateway-api/#unsupported-endpoints
func MatrixNotFound(w http.ResponseWriter, _ *http.Request) {
	metrics.TotalMatrixNotFoundRequestsMetric.Inc()
	logger.Logger.Debug("M_UNRECOGNIZED")
	w.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprintf(w, "M_UNRECOGNIZED")
	if err != nil {
		logger.Logger.Error("MatrixNotFound Error")
		return
	}
}

// MatrixNotAllowed returns a 405 Code with M_UNRECOGNIZED
// https://spec.matrix.org/v1.7/push-gateway-api/#unsupported-endpoints
func MatrixNotAllowed(w http.ResponseWriter, _ *http.Request) {
	metrics.TotalMatrixNotAllowedRequestsMetric.Inc()
	logger.Logger.Debug("M_UNRECOGNIZED")
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, err := fmt.Fprintf(w, "M_UNRECOGNIZED")
	if err != nil {
		logger.Logger.Error("MatrixNotAllowed Error")
		return
	}
}

// HealthHandler Checks the Healthstatus of the Pushgateway
func (c *WebserverConfig) HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Check the health of the server and return a status code accordingly
	if ServerIsHealthy(c) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Server is healthy")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Server is not healthy")
	}
}
func ServerIsHealthy(c *WebserverConfig) bool {
	modelsConfig := models.Config{}
	config, err := json.Marshal(c)
	if err != nil {
		logger.Logger.Error("Couldn't Marshal json of Config")
		return false
	}
	err = json.Unmarshal(config, &modelsConfig)
	if err != nil {
		logger.Logger.Error("Couldn't UnMarshal json of original Config into modelsConfig")
		return false
	}

	return true
}
