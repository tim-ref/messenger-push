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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TotalPostNotifyRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "matrix_go_push_processed_post_notify_requests_total",
		Help: "The total number of processed Notification events on Post /_matrix/push/v1/notify",
	})

	TotalFirebaseRequestsMetric = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "matrix_go_push_processed_firebase_requests_total",
			Help: "All Requests intended to sent to Firebase, that can be Rejected messages aswell as blocked client apps"},
	)

	TotalFirbaseSentRequests = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "matrix_go_push_processed_firebase_sent_requests_total",
			Help: "All Notifications that was sent to Firebase"},
	)

	TotalFirebaseRejectedRequestsMetric = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "matrix_go_push_processed_firebase_rejected_requests_total",
			Help: "All Notification Requests Firebase rejected"},
	)

	TotalWrongDeviceIdsProcessedMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "matrix_go_push_wrong_device_ids_sent_total",
		Help: "The total number of devices tried to send Push without permission",
	})

	TotalMatrixNotFoundRequestsMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "matrix_go_push_matrix_not_found_total",
		Help: "The total number of Requests the Endpoint was not found",
	})

	TotalMatrixNotAllowedRequestsMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "matrix_go_push_matrix_not_allowed_total",
		Help: "The total number of Requests that are not Allowed",
	})
)
